package mongo

import (
	// Standard Library Imports
	"context"
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

// requestMongoManager manages the main Mongo Session for a Request.
type requestMongoManager struct {
	// db contains the Mongo connection that holds the base session that can be
	// copied and closed.
	db *mgo.Database

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

func (r *requestMongoManager) Configure(ctx context.Context) error {
	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = r.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// In terms of the underlying entity for session data, the model is the
	// same across the following entities. I have decided to logically break
	// them into separate collections rather than have a 'SessionType'.
	collections := []string{
		CollectionOpenIDSessions,
		CollectionAccessTokens,
		CollectionRefreshTokens,
		CollectionAuthorizationCodes,
	}

	// Build Indices
	indices := []mgo.Index{
		{
			Name:       IdxSessionId,
			Key:        []string{"id"},
			Unique:     true,
			DropDups:   true,
			Background: true,
			Sparse:     true,
		},
		{
			Name:       IdxSignatureId,
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

		collection := r.db.C(collection).With(mgoSession)
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
func (r *requestMongoManager) getConcrete(ctx context.Context, entityName string, requestID string) (storage.Request, error) {
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": entityName,
		"method":     "getConcrete",
		"id":         requestID,
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = r.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Build Query
	query := bson.M{
		"id": requestID,
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "requestMongoManager",
		Method:  "getConcrete",
		Query:   query,
	})
	defer span.Finish()

	result := storage.Request{}
	collection := r.db.C(entityName).With(mgoSession)
	if err := collection.Find(query).One(&result); err != nil {
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
	return result, nil
}

func (r *requestMongoManager) List(ctx context.Context, entityName string, filter storage.ListRequestsRequest) ([]storage.Request, error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": entityName,
		"method":     "List",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = r.db.Session.Copy()
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
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "requestMongoManager",
		Method:  "List",
		Query:   query,
	})
	defer span.Finish()

	var results []storage.Request
	collection := r.db.C(CollectionClients).With(mgoSession)
	err := collection.Find(query).All(&results)
	if err != nil {
		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return nil, err
	}
	return results, nil
}

func (r *requestMongoManager) Create(ctx context.Context, entityName string, request storage.Request) (storage.Request, error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": entityName,
		"method":     "Create",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = r.db.Session.Copy()
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
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "requestMongoManager",
		Method:  "Create",
	})
	defer span.Finish()

	// Create resource
	collection := r.db.C(entityName).With(mgoSession)
	err := collection.Insert(request)
	if err != nil {
		if mgo.IsDup(err) {
			// Log to StdOut
			log.WithError(err).Debug(logConflict)
			// Log to OpenTracing
			otLogErr(span, err)
			return request, storage.ErrResourceExists
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogQuery(span, request)
		otLogErr(span, err)
		return request, err
	}
	return request, nil
}

func (r *requestMongoManager) Get(ctx context.Context, entityName string, requestID string) (storage.Request, error) {
	return r.getConcrete(ctx, entityName, requestID)
}

func (r *requestMongoManager) GetBySignature(ctx context.Context, entityName string, signature string) (storage.Request, error) {
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": entityName,
		"method":     "GetBySignature",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = r.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Build Query
	query := bson.M{
		"signature": signature,
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "requestMongoManager",
		Method:  "GetBySignature",
		Query:   query,
	})
	defer span.Finish()

	result := storage.Request{}
	collection := r.db.C(entityName).With(mgoSession)
	if err := collection.Find(query).One(&result); err != nil {
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
	return result, nil
}

func (r *requestMongoManager) Update(ctx context.Context, entityName string, requestID string, updatedRequest storage.Request) (storage.Request, error) {
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
		mgoSession = r.db.Session.Copy()
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
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager:  "requestMongoManager",
		Method:   "Update",
		Selector: selector,
	})
	defer span.Finish()

	collection := r.db.C(CollectionClients).With(mgoSession)
	if err := collection.Update(selector, updatedRequest); err != nil {
		if err == mgo.ErrNotFound {
			// Log to StdOut
			log.WithError(err).Debug(logNotFound)
			// Log to OpenTracing
			otLogErr(span, err)
			return updatedRequest, fosite.ErrNotFound
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogQuery(span, updatedRequest)
		otLogErr(span, err)
		return updatedRequest, err
	}
	return updatedRequest, nil
}

func (r *requestMongoManager) Delete(ctx context.Context, entityName string, requestID string) error {
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
		mgoSession = r.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Build Query
	query := bson.M{
		"id": requestID,
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "requestMongoManager",
		Method:  "Delete",
		Query:   query,
	})
	defer span.Finish()

	collection := r.db.C(entityName).With(mgoSession)
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

// RevokeRefreshToken finds a token stored in cache based on request ID and deletes the session by signature.
func (r *requestMongoManager) RevokeRefreshToken(ctx context.Context, requestID string) error {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": CollectionCacheRefreshTokens,
		"method":     "RevokeRefreshToken",
		"id":         requestID,
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = r.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "requestMongoManager",
		Method:  "RevokeRefreshToken",
		Query:   requestID,
	})
	defer span.Finish()

	cacheObject, err := r.Cache.Get(ctx, requestID, CollectionCacheRefreshTokens)
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

	err = r.Cache.Delete(ctx, cacheObject.Key(), CollectionCacheRefreshTokens)
	if err != nil {
		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return err
	}

	return nil
}

// RevokeAccessToken finds a token stored in cache based on request ID and deletes the session by signature.
func (r *requestMongoManager) RevokeAccessToken(ctx context.Context, requestID string) error {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": CollectionCacheAccessTokens,
		"method":     "RevokeAccessToken",
		"id":         requestID,
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = r.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "requestMongoManager",
		Method:  "RevokeAccessToken",
		Query:   requestID,
	})
	defer span.Finish()

	cacheObject, err := r.Cache.Get(ctx, requestID, CollectionCacheAccessTokens)
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

	err = r.Cache.Delete(ctx, cacheObject.Key(), CollectionCacheAccessTokens)
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
func toMongo(signature string, r fosite.Requester) storage.Request {
	return storage.Request{
		ID:            r.GetID(),
		RequestedAt:   r.GetRequestedAt(),
		Signature:     signature,
		ClientID:      r.GetClient().GetID(),
		UserID:        r.GetSession().GetSubject(),
		Scopes:        r.GetRequestedScopes(),
		GrantedScopes: r.GetGrantedScopes(),
		Form:          r.GetRequestForm(),
		Session:       r.GetSession(),
	}
}
