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

// cacheMongoManager provides a cache implementation in MongoDB for auth
// sessions.
type cacheMongoManager struct {
	db *mgo.Database
}

// Configure sets up the Mongo collection for cache resources.
func (c *cacheMongoManager) Configure() error {
	if err := c.configureAccessTokensCollection(); err != nil {
		return err
	}

	if err := c.configureRefreshTokensCollection(); err != nil {
		return err
	}

	return nil
}

// configureAccessTokensCollection sets indices for the Access Token Collection.
func (c *cacheMongoManager) configureAccessTokensCollection() error {
	log := logger.WithFields(logrus.Fields{
		"package":    "datastore",
		"driver":     "mongo",
		"collection": CollectionCacheAccessTokens,
		"method":     "Configure",
	})

	collection := c.db.C(CollectionCacheAccessTokens).With(c.db.Session.Clone())
	defer collection.Database.Session.Close()

	// Ensure index on request id
	index := mgo.Index{
		Name:       IdxCacheRequestId,
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
func (c *cacheMongoManager) configureRefreshTokensCollection() error {
	log := logger.WithFields(logrus.Fields{
		"package":    "datastore",
		"driver":     "mongo",
		"collection": CollectionCacheRefreshTokens,
		"method":     "Configure",
	})

	collection := c.db.C(CollectionCacheRefreshTokens).With(c.db.Session.Clone())
	defer collection.Database.Session.Close()

	// Ensure index on request id
	index := mgo.Index{
		Name:       IdxCacheRequestId,
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
func (c *cacheMongoManager) getConcrete(ctx context.Context, entityName, key string) (value storage.SessionCache, err error) {
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": entityName,
		"method":     "getConcrete",
		"key":        key,
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = c.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Build Query
	query := bson.M{
		"id": key,
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "cacheMongoManager",
		Method:  "getConcrete",
		Query:   query,
	})
	defer span.Finish()

	result := storage.SessionCache{}
	collection := c.db.C(entityName).With(mgoSession)
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

func (c *cacheMongoManager) Create(ctx context.Context, entityName string, cacheObject storage.SessionCache) (storage.SessionCache, error) {
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
		mgoSession = c.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	if cacheObject.CreateTime == 0 {
		cacheObject.CreateTime = time.Now().Unix()
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "cacheMongoManager",
		Method:  "Create",
	})
	defer span.Finish()

	// Create resource
	collection := c.db.C(entityName).With(mgoSession)
	err := collection.Insert(c)
	if err != nil {
		if mgo.IsDup(err) {
			// Log to StdOut
			log.WithError(err).Debug(logConflict)
			// Log to OpenTracing
			otLogErr(span, err)
			return cacheObject, storage.ErrCacheSessionExists
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogQuery(span, cacheObject)
		otLogErr(span, err)
		return cacheObject, err
	}
	return cacheObject, nil
}

func (c *cacheMongoManager) Get(ctx context.Context, entityName string, key string) (storage.SessionCache, error) {
	return c.getConcrete(ctx, entityName, key)
}

func (c *cacheMongoManager) Update(ctx context.Context, entityName string, updatedCacheObject storage.SessionCache) (storage.SessionCache, error) {
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
		mgoSession = c.db.Session.Copy()
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
		Manager:  "cacheMongoManager",
		Method:   "Update",
		Selector: selector,
	})
	defer span.Finish()

	collection := c.db.C(entityName).With(mgoSession)
	if err := collection.Update(selector, updatedCacheObject); err != nil {
		if err == mgo.ErrNotFound {
			// Log to StdOut
			log.WithError(err).Debug(logNotFound)
			// Log to OpenTracing
			otLogErr(span, err)
			return updatedCacheObject, fosite.ErrNotFound
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogQuery(span, updatedCacheObject)
		otLogErr(span, err)
		return updatedCacheObject, err
	}
	return updatedCacheObject, nil
}

func (c *cacheMongoManager) Delete(ctx context.Context, entityName string, key string) error {
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
		mgoSession = c.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Build Query
	query := bson.M{
		"id": key,
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "clientMongoManager",
		Method:  "Delete",
		Query:   query,
	})
	defer span.Finish()

	collection := c.db.C(entityName).With(mgoSession)
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
