package mongo

import (
	// Standard Library Imports
	"context"
	"time"

	// External Imports
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/ory/fosite"
	"github.com/sirupsen/logrus"

	// Internal Imports
	"github.com/matthewhartstonge/storage"
)

// CacheManager provides a cache implementation in MongoDB for auth
// sessions.
type CacheManager struct {
	DB *mgo.Database
}

// Configure sets up the Mongo collection for cache resources.
func (c *CacheManager) Configure(ctx context.Context) error {
	if err := c.configureAccessTokensCollection(ctx); err != nil {
		return err
	}

	if err := c.configureRefreshTokensCollection(ctx); err != nil {
		return err
	}

	return nil
}

// configureAccessTokensCollection sets indices for the Access Token Collection.
func (c *CacheManager) configureAccessTokensCollection(ctx context.Context) error {
	log := logger.WithFields(logrus.Fields{
		"package":    "datastore",
		"driver":     "mongo",
		"collection": storage.EntityCacheAccessTokens,
		"method":     "Configure",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = c.DB.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	collection := c.DB.C(storage.EntityCacheAccessTokens).With(mgoSession)

	// Ensure index on request id
	index := mgo.Index{
		Name:       IdxCacheRequestID,
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := collection.EnsureIndex(index)
	if err != nil {
		log.WithError(err).Error(logError)
		return err
	}

	// Ensure index on request signature
	index = mgo.Index{
		Name:       IdxCacheRequestSignature,
		Key:        []string{"signature"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = collection.EnsureIndex(index)
	if err != nil {
		log.WithError(err).Error(logError)
		return err
	}

	return nil
}

// configureRefreshTokensCollection sets indices for the Refresh Token
// Collection.
func (c *CacheManager) configureRefreshTokensCollection(ctx context.Context) error {
	log := logger.WithFields(logrus.Fields{
		"package":    "datastore",
		"driver":     "mongo",
		"collection": storage.EntityCacheRefreshTokens,
		"method":     "Configure",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = c.DB.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	collection := c.DB.C(storage.EntityCacheRefreshTokens).With(mgoSession)

	// Ensure index on request id
	index := mgo.Index{
		Name:       IdxCacheRequestID,
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := collection.EnsureIndex(index)
	if err != nil {
		log.WithError(err).Error(logError)
		return err
	}

	// Ensure index on request signature
	index = mgo.Index{
		Name:       IdxCacheRequestSignature,
		Key:        []string{"signature"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = collection.EnsureIndex(index)
	if err != nil {
		log.WithError(err).Error(logError)
		return err
	}

	return nil
}

// getConcrete returns a map of request id to token signature.
func (c *CacheManager) getConcrete(ctx context.Context, entityName, key string) (result storage.SessionCache, err error) {
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": entityName,
		"method":     "getConcrete",
		"key":        key,
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = c.DB.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Build Query
	query := bson.M{
		"id": key,
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "CacheManager",
		Method:  "getConcrete",
		Query:   query,
	})
	defer span.Finish()

	var sessionCache = storage.SessionCache{}
	collection := c.DB.C(entityName).With(mgoSession)
	if err := collection.Find(query).One(&sessionCache); err != nil {
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
	return sessionCache, nil
}

// Create creates a new Cache resource and returns the newly created Cache
// resource.
func (c *CacheManager) Create(ctx context.Context, entityName string, cacheObject storage.SessionCache) (result storage.SessionCache, err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": entityName,
		"method":     "Create",
		"key":        cacheObject.Key(),
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = c.DB.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	if cacheObject.CreateTime == 0 {
		cacheObject.CreateTime = time.Now().Unix()
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "CacheManager",
		Method:  "Create",
	})
	defer span.Finish()

	// Create resource
	collection := c.DB.C(entityName).With(mgoSession)
	err = collection.Insert(cacheObject)
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
		otLogQuery(span, cacheObject)
		otLogErr(span, err)
		return result, err
	}
	return cacheObject, nil
}

// Get returns the specified Cache resource.
func (c *CacheManager) Get(ctx context.Context, entityName string, key string) (result storage.SessionCache, err error) {
	return c.getConcrete(ctx, entityName, key)
}

// Update updates the Cache resource and attributes and returns the updated
// Cache resource.
func (c *CacheManager) Update(ctx context.Context, entityName string, updatedCacheObject storage.SessionCache) (result storage.SessionCache, err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": entityName,
		"method":     "Update",
		"key":        updatedCacheObject.Key(),
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = c.DB.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Update modified time
	updatedCacheObject.UpdateTime = time.Now().Unix()

	// Build Query
	selector := bson.M{
		"id": updatedCacheObject.ID,
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager:  "CacheManager",
		Method:   "Update",
		Selector: selector,
	})
	defer span.Finish()

	collection := c.DB.C(entityName).With(mgoSession)
	if err := collection.Update(selector, updatedCacheObject); err != nil {
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
		otLogQuery(span, updatedCacheObject)
		otLogErr(span, err)
		return result, err
	}
	return updatedCacheObject, nil
}

// Delete deletes the specified Cache resource.
func (c *CacheManager) Delete(ctx context.Context, entityName string, key string) error {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": entityName,
		"method":     "Delete",
		"key":        key,
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = c.DB.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Build Query
	query := bson.M{
		"id": key,
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "ClientManager",
		Method:  "Delete",
		Query:   query,
	})
	defer span.Finish()

	collection := c.DB.C(entityName).With(mgoSession)
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

// DeleteByValue deletes a Cache resource by matching on value.
func (c *CacheManager) DeleteByValue(ctx context.Context, entityName string, value string) error {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": entityName,
		"method":     "DeleteByValue",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = c.DB.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Build Query
	query := bson.M{
		"signature": value,
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "ClientManager",
		Method:  "DeleteByValue",
		Query:   query,
	})
	defer span.Finish()

	collection := c.DB.C(entityName).With(mgoSession)
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
