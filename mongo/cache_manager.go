package mongo

import (
	// Standard Library Imports
	"context"
	"time"

	// External Imports
	"github.com/ory/fosite"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// Internal Imports
	"github.com/matthewhartstonge/storage"
)

// CacheManager provides a cache implementation in MongoDB for auth
// sessions.
type CacheManager struct {
	DB *mongo.Database
}

// Configure sets up the Mongo collection for cache resources.
func (c *CacheManager) Configure(ctx context.Context) error {
	if err := c.configure(ctx, storage.EntityCacheAccessTokens); err != nil {
		return err
	}

	if err := c.configure(ctx, storage.EntityCacheRefreshTokens); err != nil {
		return err
	}

	return nil
}

func (c *CacheManager) configure(ctx context.Context, entityName string) (err error) {
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"driver":     "mongo",
		"collection": entityName,
		"method":     "configure",
	})

	// Copy a new DB session if none specified
	_, ok := ContextToSession(ctx)
	if !ok {
		var closer func()
		ctx, _, closer, err = newSession(ctx, c.DB)
		if err != nil {
			log.WithError(err).Debug("error starting session")
			return err
		}
		defer closer()
	}

	// Ensure index on request id
	indices := []mongo.IndexModel{
		{
			Keys: bson.D{
				{
					Key:   "id",
					Value: int32(1),
				},
			},
			Options: options.Index().
				SetBackground(true).
				SetName(IdxCacheRequestID).
				SetSparse(true).
				SetUnique(true),
		},
		{
			Keys: bson.D{
				{
					Key:   "signature",
					Value: int32(1),
				},
			},
			Options: options.Index().
				SetBackground(true).
				SetName(IdxCacheRequestSignature).
				SetSparse(true).
				SetUnique(true),
		},
	}

	collection := c.DB.Collection(entityName)
	_, err = collection.Indexes().CreateMany(ctx, indices)
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
	_, ok := ContextToSession(ctx)
	if !ok {
		var closer func()
		ctx, _, closer, err = newSession(ctx, c.DB)
		if err != nil {
			log.WithError(err).Debug("error starting session")
			return result, err
		}
		defer closer()
	}

	// Build Query
	query := bson.M{
		"id": key,
	}

	// Trace how long the Mongo operation takes to complete.
	span, _ := traceMongoCall(ctx, dbTrace{
		Manager: "CacheManager",
		Method:  "getConcrete",
		Query:   query,
	})
	defer span.Finish()

	var sessionCache storage.SessionCache
	collection := c.DB.Collection(entityName)
	if err := collection.FindOne(ctx, query).Decode(&sessionCache); err != nil {
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
	_, ok := ContextToSession(ctx)
	if !ok {
		var closer func()
		ctx, _, closer, err = newSession(ctx, c.DB)
		if err != nil {
			log.WithError(err).Debug("error starting session")
			return result, err
		}
		defer closer()
	}

	if cacheObject.CreateTime == 0 {
		cacheObject.CreateTime = time.Now().Unix()
	}

	// Trace how long the Mongo operation takes to complete.
	span, _ := traceMongoCall(ctx, dbTrace{
		Manager: "CacheManager",
		Method:  "Create",
	})
	defer span.Finish()

	// Create resource
	collection := c.DB.Collection(entityName)
	_, err = collection.InsertOne(ctx, cacheObject)
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
	_, ok := ContextToSession(ctx)
	if !ok {
		var closer func()
		ctx, _, closer, err = newSession(ctx, c.DB)
		if err != nil {
			log.WithError(err).Debug("error starting session")
			return result, err
		}
		defer closer()
	}

	// Update modified time
	updatedCacheObject.UpdateTime = time.Now().Unix()

	// Build Query
	selector := bson.M{
		"id": updatedCacheObject.ID,
	}

	// Trace how long the Mongo operation takes to complete.
	span, _ := traceMongoCall(ctx, dbTrace{
		Manager:  "CacheManager",
		Method:   "Update",
		Selector: selector,
	})
	defer span.Finish()

	collection := c.DB.Collection(entityName)
	res, err := collection.ReplaceOne(ctx, selector, updatedCacheObject)
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
		otLogQuery(span, updatedCacheObject)
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

	return updatedCacheObject, nil
}

// Delete deletes the specified Cache resource.
func (c *CacheManager) Delete(ctx context.Context, entityName string, key string) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": entityName,
		"method":     "Delete",
		"key":        key,
	})

	// Copy a new DB session if none specified
	_, ok := ContextToSession(ctx)
	if !ok {
		var closer func()
		ctx, _, closer, err = newSession(ctx, c.DB)
		if err != nil {
			log.WithError(err).Debug("error starting session")
			return err
		}
		defer closer()
	}

	// Build Query
	query := bson.M{
		"id": key,
	}

	// Trace how long the Mongo operation takes to complete.
	span, _ := traceMongoCall(ctx, dbTrace{
		Manager: "CacheManager",
		Method:  "Delete",
		Query:   query,
	})
	defer span.Finish()

	collection := c.DB.Collection(entityName)
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

// DeleteByValue deletes a Cache resource by matching on value.
func (c *CacheManager) DeleteByValue(ctx context.Context, entityName string, value string) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": entityName,
		"method":     "DeleteByValue",
	})

	// Copy a new DB session if none specified
	_, ok := ContextToSession(ctx)
	if !ok {
		var closer func()
		ctx, _, closer, err = newSession(ctx, c.DB)
		if err != nil {
			log.WithError(err).Debug("error starting session")
			return err
		}
		defer closer()
	}

	// Build Query
	query := bson.M{
		"signature": value,
	}

	// Trace how long the Mongo operation takes to complete.
	span, _ := traceMongoCall(ctx, dbTrace{
		Manager: "CacheManager",
		Method:  "DeleteByValue",
		Query:   query,
	})
	defer span.Finish()

	collection := c.DB.Collection(entityName)
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
