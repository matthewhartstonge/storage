package mongo

import (
	// Standard Library Imports
	"context"
	"encoding/json"
	"time"

	// External Imports
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/ory/fosite"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"

	// Internal Imports
	"github.com/matthewhartstonge/storage"
)

// RequestManager manages the main Mongo Session for a Request.
type RequestManager struct {
	// DB contains the Mongo connection that holds the base session that can be
	// copied and closed.
	DB *mgo.Database

	// Cache provides access to Cache entities in order to create, read,
	// update and delete resources from the caching collection.
	Cache storage.CacheStorer

	// Clients provides access to Client entities in order to create, read,
	// update and delete resources from the clients collection.
	// A client is required when cross referencing scope access rights.
	Clients storage.ClientStorer

	// Users provides access to User entities in order to create, read, update
	// and delete resources from the user collection.
	// Users are required when the Password Credentials Grant, is implemented
	// in order to find and authenticate users.
	Users storage.UserStorer
}

// Configure implements storage.Configurer.
func (r *RequestManager) Configure(ctx context.Context) error {
	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = r.DB.Session.Copy()
		defer mgoSession.Close()
	}

	// In terms of the underlying entity for session data, the model is the
	// same across the following entities. I have decided to logically break
	// them into separate collections rather than have a 'SessionType'.
	collections := []string{
		storage.EntityAccessTokens,
		storage.EntityAuthorizationCodes,
		storage.EntityOpenIDSessions,
		storage.EntityPKCESessions,
		storage.EntityRefreshTokens,
	}

	// Build Indices
	indices := []mgo.Index{
		{
			Name:       IdxSessionID,
			Key:        []string{"id"},
			Unique:     true,
			DropDups:   true,
			Background: true,
			Sparse:     true,
		},
		{
			Name:       IdxSignatureID,
			Key:        []string{"signature"},
			Unique:     true,
			DropDups:   true,
			Background: true,
			Sparse:     true,
		},
		{
			Name:       IdxCompoundRequester,
			Key:        []string{"clientId", "userId"},
			Unique:     false,
			DropDups:   false,
			Background: true,
			Sparse:     true,
		},
	}

	for _, collection := range collections {
		log := logger.WithFields(logrus.Fields{
			"package":    "mongo",
			"collection": collection,
			"method":     "Configure",
		})

		collection := r.DB.C(collection).With(mgoSession)
		for _, index := range indices {
			err := collection.EnsureIndex(index)
			if err != nil {
				log.WithError(err).Error(logError)
				return err
			}
		}
	}
	return nil
}

// getConcrete returns a Request resource.
func (r *RequestManager) getConcrete(ctx context.Context, entityName string, requestID string) (result storage.Request, err error) {
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": entityName,
		"method":     "getConcrete",
		"id":         requestID,
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = r.DB.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Build Query
	query := bson.M{
		"id": requestID,
	}

	// Trace how long the Mongo operation takes to complete.
	span, _ := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "getConcrete",
		Query:   query,
	})
	defer span.Finish()

	request := storage.Request{}
	collection := r.DB.C(entityName).With(mgoSession)
	if err := collection.Find(query).One(&request); err != nil {
		if err == mgo.ErrNotFound {
			log.WithError(err).Debug(logNotFound)
			return result, fosite.ErrNotFound
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return result, err
	}
	return request, nil
}

// List returns a list of Request resources that match the provided inputs.
func (r *RequestManager) List(ctx context.Context, entityName string, filter storage.ListRequestsRequest) (results []storage.Request, err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": entityName,
		"method":     "List",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = r.DB.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Build Query
	query := bson.M{}
	if filter.ClientID != "" {
		query["clientId"] = filter.ClientID
	}
	if filter.UserID != "" {
		query["userId"] = filter.UserID
	}
	if len(filter.ScopesIntersection) > 0 {
		query["scopes"] = bson.M{"$all": filter.ScopesIntersection}
	}
	if len(filter.ScopesUnion) > 0 {
		query["scopes"] = bson.M{"$in": filter.ScopesUnion}
	}
	if len(filter.GrantedScopesIntersection) > 0 {
		query["scopes"] = bson.M{"$all": filter.GrantedScopesIntersection}
	}
	if len(filter.GrantedScopesUnion) > 0 {
		query["scopes"] = bson.M{"$in": filter.GrantedScopesUnion}
	}

	// Trace how long the Mongo operation takes to complete.
	span, _ := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "List",
		Query:   query,
	})
	defer span.Finish()

	var requests []storage.Request
	collection := r.DB.C(entityName).With(mgoSession)
	err = collection.Find(query).All(&requests)
	if err != nil {
		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return results, err
	}
	return requests, nil
}

// Create creates the new Request resource and returns the newly created Request
// resource.
func (r *RequestManager) Create(ctx context.Context, entityName string, request storage.Request) (result storage.Request, err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": entityName,
		"method":     "Create",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = r.DB.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Enable developers to provide their own IDs
	if request.ID == "" {
		request.ID = uuid.New()
	}
	if request.CreateTime == 0 {
		request.CreateTime = time.Now().Unix()
	}
	if request.RequestedAt.IsZero() {
		request.RequestedAt = time.Now()
	}

	// Trace how long the Mongo operation takes to complete.
	span, _ := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "Create",
	})
	defer span.Finish()

	// Create resource
	collection := r.DB.C(entityName).With(mgoSession)
	err = collection.Insert(request)
	if err != nil {
		if mgo.IsDup(err) {
			// Log to StdOut
			log.WithError(err).Debug(logConflict)
			// Log to OpenTracing
			otLogErr(span, err)
			return result, storage.ErrResourceExists
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogQuery(span, request)
		otLogErr(span, err)
		return result, err
	}
	return request, nil
}

// Get returns the specified Request resource.
func (r *RequestManager) Get(ctx context.Context, entityName string, requestID string) (result storage.Request, err error) {
	return r.getConcrete(ctx, entityName, requestID)
}

// GetBySignature returns a Request resource, if the presented signature returns
// a match.
func (r *RequestManager) GetBySignature(ctx context.Context, entityName string, signature string) (result storage.Request, err error) {
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": entityName,
		"method":     "GetBySignature",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = r.DB.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Build Query
	query := bson.M{
		"signature": signature,
	}

	// Trace how long the Mongo operation takes to complete.
	span, _ := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "GetBySignature",
		Query:   query,
	})
	defer span.Finish()

	request := storage.Request{}
	collection := r.DB.C(entityName).With(mgoSession)
	if err := collection.Find(query).One(&request); err != nil {
		if err == mgo.ErrNotFound {
			log.WithError(err).Debug(logNotFound)
			return result, fosite.ErrNotFound
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return result, err
	}
	return request, nil
}

// Update updates the Request resource and attributes and returns the updated
// Request resource.
func (r *RequestManager) Update(ctx context.Context, entityName string, requestID string, updatedRequest storage.Request) (result storage.Request, err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": entityName,
		"method":     "Update",
		"id":         requestID,
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = r.DB.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Deny updating the entity Id
	updatedRequest.ID = requestID
	// Update modified time
	updatedRequest.UpdateTime = time.Now().Unix()

	// Build Query
	selector := bson.M{
		"id": requestID,
	}

	// Trace how long the Mongo operation takes to complete.
	span, _ := traceMongoCall(ctx, dbTrace{
		Manager:  "RequestManager",
		Method:   "Update",
		Selector: selector,
	})
	defer span.Finish()

	collection := r.DB.C(entityName).With(mgoSession)
	if err := collection.Update(selector, updatedRequest); err != nil {
		if err == mgo.ErrNotFound {
			// Log to StdOut
			log.WithError(err).Debug(logNotFound)
			// Log to OpenTracing
			otLogErr(span, err)
			return result, fosite.ErrNotFound
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogQuery(span, updatedRequest)
		otLogErr(span, err)
		return result, err
	}
	return updatedRequest, nil
}

// Delete deletes the specified Request resource.
func (r *RequestManager) Delete(ctx context.Context, entityName string, requestID string) error {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": entityName,
		"method":     "Delete",
		"id":         requestID,
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = r.DB.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Build Query
	query := bson.M{
		"id": requestID,
	}

	// Trace how long the Mongo operation takes to complete.
	span, _ := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "Delete",
		Query:   query,
	})
	defer span.Finish()

	collection := r.DB.C(entityName).With(mgoSession)
	if err := collection.Remove(query); err != nil {
		if err == mgo.ErrNotFound {
			// Log to StdOut
			log.WithError(err).Debug(logNotFound)
			// Log to OpenTracing
			otLogErr(span, err)
			return fosite.ErrNotFound
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return err
	}
	return nil
}

// DeleteBySignature deletes the specified Cache resource, if the presented
// signature returns a match.
func (r *RequestManager) DeleteBySignature(ctx context.Context, entityName string, signature string) error {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": entityName,
		"method":     "DeleteBySignature",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = r.DB.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Build Query
	query := bson.M{
		"signature": signature,
	}

	// Trace how long the Mongo operation takes to complete.
	span, _ := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "DeleteBySignature",
		Query:   query,
	})
	defer span.Finish()

	collection := r.DB.C(entityName).With(mgoSession)
	if err := collection.Remove(query); err != nil {
		if err == mgo.ErrNotFound {
			// Log to StdOut
			log.WithError(err).Debug(logNotFound)
			// Log to OpenTracing
			otLogErr(span, err)
			return fosite.ErrNotFound
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return err
	}
	return nil
}

// RevokeRefreshToken finds a token stored in cache based on request ID and
// deletes the session by signature.
func (r *RequestManager) RevokeRefreshToken(ctx context.Context, requestID string) error {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityCacheRefreshTokens,
		"method":     "RevokeRefreshToken",
		"id":         requestID,
	})

	// Copy a new DB session if none specified
	_, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession := r.DB.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "RevokeRefreshToken",
		Query:   requestID,
	})
	defer span.Finish()

	cacheObject, err := r.Cache.Get(ctx, storage.EntityCacheRefreshTokens, requestID)
	if err != nil {
		if err == fosite.ErrNotFound {
			// Log to OpenTracing
			otLogErr(span, err)
			log.WithError(err).Debug(logNotFound)
			return err
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return err
	}

	err = r.DeleteRefreshTokenSession(ctx, cacheObject.Value())
	if err != nil {
		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return err
	}

	err = r.Cache.Delete(ctx, storage.EntityCacheRefreshTokens, cacheObject.Key())
	if err != nil {
		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return err
	}

	return nil
}

// RevokeAccessToken finds a token stored in cache based on request ID and
// deletes the session by signature.
func (r *RequestManager) RevokeAccessToken(ctx context.Context, requestID string) error {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityCacheAccessTokens,
		"method":     "RevokeAccessToken",
		"id":         requestID,
	})

	// Copy a new DB session if none specified
	_, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession := r.DB.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "RevokeAccessToken",
		Query:   requestID,
	})
	defer span.Finish()

	cacheObject, err := r.Cache.Get(ctx, storage.EntityCacheAccessTokens, requestID)
	if err != nil {
		if err == fosite.ErrNotFound {
			// Log to OpenTracing
			otLogErr(span, err)
			log.WithError(err).Debug(logNotFound)
			return err
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return err
	}

	err = r.DeleteAccessTokenSession(ctx, cacheObject.Value())
	if err != nil {
		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return err
	}

	return nil
}

// toMongo transforms a fosite.Request to a storage.Request
// Signature is a hash that relates to the underlying request method and may not
// be a strict 'signature', for example, authorization code grant passes in an
// authorization code.
func toMongo(signature string, r fosite.Requester) storage.Request {
	session, _ := json.Marshal(r.GetSession())
	return storage.Request{
		ID:                r.GetID(),
		RequestedAt:       r.GetRequestedAt(),
		Signature:         signature,
		ClientID:          r.GetClient().GetID(),
		UserID:            r.GetSession().GetSubject(),
		RequestedScope:    r.GetRequestedScopes(),
		GrantedScope:      r.GetGrantedScopes(),
		RequestedAudience: r.GetRequestedAudience(),
		GrantedAudience:   r.GetGrantedAudience(),
		Form:              r.GetRequestForm(),
		Active:            true,
		Session:           session,
	}
}
