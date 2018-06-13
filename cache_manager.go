package storage

import "context"

// Cacher provides a generic interface for storing data in a cache
type Cacher interface {
	Key() string
	Value() string
}

// CacheManager provides a generic interface to key value cache objects in
// order to build a cache datastore.
type CacheManager interface {
	CacheStorer
}

// Storer provides a way to create cache based objects in mongo
type CacheStorer interface {
	Create(ctx context.Context, entityName string, cacheObject SessionCache) (SessionCache, error)
	Get(ctx context.Context, entityName string, key string) (SessionCache, error)
	Update(ctx context.Context, entityName string, cacheObject SessionCache) (SessionCache, error)
	Delete(ctx context.Context, entityName string, key string) error
	DeleteByValue(ctx context.Context, entityName string, value string) error
}
