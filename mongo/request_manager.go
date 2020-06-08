package mongo

import (
	// Standard Library Imports
	"context"
	"encoding/json"
	"time"

	// External Imports
	"github.com/ory/fosite"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// Internal Imports
	"github.com/matthewhartstonge/storage"
)

// RequestManager manages the main Mongo Session for a Request.
type RequestManager struct {
	// DB contains the Mongo connection that holds the base session that can be
	// copied and closed.
	DB *mongo.Database

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
func (r *RequestManager) Configure(ctx context.Context) (err error) {
	log := logger.WithFields(logrus.Fields{
		"package": "mongo",
		"struct":  "requestManager",
		"method":  "Configure",
	})

	// Copy a new DB session if none specified
	_, ok := ContextToSession(ctx)
	if !ok {
		var closer func()
		ctx, _, closer, err = newSession(ctx, r.DB)
		if err != nil {
			log.WithError(err).Debug("error starting session")
			return err
		}
		defer closer()
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
	indices := []mongo.IndexModel{
		{
			Keys: bson.M{
				"id": 1,
			},
			Options: options.Index().
				SetBackground(true).
				SetName(IdxSessionID).
				SetSparse(true).
				SetUnique(true),
		},
		{
			Keys: bson.M{
				"signature": 1,
			},
			Options: options.Index().
				SetBackground(true).
				SetName(IdxSignatureID).
				SetSparse(true).
				SetUnique(true),
		},
		{
			Keys: bson.M{
				"clientId": 1,
				"userId":   1,
			},
			Options: options.Index().
				SetBackground(true).
				SetName(IdxCompoundRequester).
				SetSparse(true),
		},
	}

	for _, collection := range collections {
		log := logger.WithFields(logrus.Fields{
			"package":    "mongo",
			"collection": collection,
			"method":     "Configure",
		})

		collection := r.DB.Collection(collection)
		_, err = collection.Indexes().CreateMany(ctx, indices)
		if err != nil {
			log.WithError(err).Error(logError)
			return err
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
	_, ok := ContextToSession(ctx)
	if !ok {
		var closer func()
		ctx, _, closer, err = newSession(ctx, r.DB)
		if err != nil {
			log.WithError(err).Debug("error starting session")
			return result, err
		}
		defer closer()
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

	var request storage.Request
	collection := r.DB.Collection(entityName)
	err = collection.FindOne(ctx, query).Decode(&request)
	if err != nil {
		if err == mongo.ErrNoDocuments {
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
	_, ok := ContextToSession(ctx)
	if !ok {
		var closer func()
		ctx, _, closer, err = newSession(ctx, r.DB)
		if err != nil {
			log.WithError(err).Debug("error starting session")
			return results, err
		}
		defer closer()
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

	collection := r.DB.Collection(entityName)
	cursor, err := collection.Find(ctx, query)
	if err != nil {
		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return results, err
	}

	var requests []storage.Request
	err = cursor.All(ctx, &requests)
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
	_, ok := ContextToSession(ctx)
	if !ok {
		var closer func()
		ctx, _, closer, err = newSession(ctx, r.DB)
		if err != nil {
			log.WithError(err).Debug("error starting session")
			return result, err
		}
		defer closer()
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
	collection := r.DB.Collection(entityName)
	_, err = collection.InsertOne(ctx, request)
	if err != nil {
		if isDup(err) {
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
	_, ok := ContextToSession(ctx)
	if !ok {
		var closer func()
		ctx, _, closer, err = newSession(ctx, r.DB)
		if err != nil {
			log.WithError(err).Debug("error starting session")
			return result, err
		}
		defer closer()
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

	var request storage.Request
	collection := r.DB.Collection(entityName)
	err = collection.FindOne(ctx, query).Decode(&request)
	if err != nil {
		if err == mongo.ErrNoDocuments {
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
	_, ok := ContextToSession(ctx)
	if !ok {
		var closer func()
		ctx, _, closer, err = newSession(ctx, r.DB)
		if err != nil {
			log.WithError(err).Debug("error starting session")
			return result, err
		}
		defer closer()
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

	collection := r.DB.Collection(entityName)
	res, err := collection.ReplaceOne(ctx, selector, updatedRequest)
	if err != nil {
		if isDup(err) {
			// Log to StdOut
			log.WithError(err).Debug(logConflict)
			// Log to OpenTracing
			otLogErr(span, err)
			return result, storage.ErrResourceExists
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogQuery(span, updatedRequest)
		otLogErr(span, err)
		return result, err
	}

	if res.MatchedCount == 0 {
		// Log to StdOut
		log.WithError(err).Debug(logNotFound)
		// Log to OpenTracing
		otLogErr(span, err)
		return result, fosite.ErrNotFound
	}

	return updatedRequest, nil
}

// Delete deletes the specified Request resource.
func (r *RequestManager) Delete(ctx context.Context, entityName string, requestID string) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": entityName,
		"method":     "Delete",
		"id":         requestID,
	})

	// Copy a new DB session if none specified
	_, ok := ContextToSession(ctx)
	if !ok {
		var closer func()
		ctx, _, closer, err = newSession(ctx, r.DB)
		if err != nil {
			log.WithError(err).Debug("error starting session")
			return err
		}
		defer closer()
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

	collection := r.DB.Collection(entityName)
	res, err := collection.DeleteOne(ctx, query)
	if err != nil {
		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return err
	}

	if res.DeletedCount == 0 {
		// Log to StdOut
		log.WithError(err).Debug(logNotFound)
		// Log to OpenTracing
		otLogErr(span, err)
		return fosite.ErrNotFound
	}

	return nil
}

// DeleteBySignature deletes the specified Cache resource, if the presented
// signature returns a match.
func (r *RequestManager) DeleteBySignature(ctx context.Context, entityName string, signature string) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": entityName,
		"method":     "DeleteBySignature",
	})

	// Copy a new DB session if none specified
	_, ok := ContextToSession(ctx)
	if !ok {
		var closer func()
		ctx, _, closer, err = newSession(ctx, r.DB)
		if err != nil {
			log.WithError(err).Debug("error starting session")
			return err
		}
		defer closer()
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

	collection := r.DB.Collection(entityName)
	res, err := collection.DeleteOne(ctx, query)
	if err != nil {
		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return err
	}

	if res.DeletedCount == 0 {
		// Log to StdOut
		log.WithError(err).Debug(logNotFound)
		// Log to OpenTracing
		otLogErr(span, err)
		return fosite.ErrNotFound
	}

	return nil
}

// RevokeRefreshToken finds a token stored in cache based on request ID and
// deletes the session by signature.
func (r *RequestManager) RevokeRefreshToken(ctx context.Context, requestID string) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityCacheRefreshTokens,
		"method":     "RevokeRefreshToken",
		"id":         requestID,
	})

	// Copy a new DB session if none specified
	_, ok := ContextToSession(ctx)
	if !ok {
		var closer func()
		ctx, _, closer, err = newSession(ctx, r.DB)
		if err != nil {
			log.WithError(err).Debug("error starting session")
			return err
		}
		defer closer()
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
func (r *RequestManager) RevokeAccessToken(ctx context.Context, requestID string) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityCacheAccessTokens,
		"method":     "RevokeAccessToken",
		"id":         requestID,
	})

	// Copy a new DB session if none specified
	_, ok := ContextToSession(ctx)
	if !ok {
		var closer func()
		ctx, _, closer, err = newSession(ctx, r.DB)
		if err != nil {
			log.WithError(err).Debug("error starting session")
			return err
		}
		defer closer()
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

	err = r.Cache.Delete(ctx, storage.EntityCacheAccessTokens, cacheObject.Key())
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
		Active:            false,
		Session:           session,
	}
}
