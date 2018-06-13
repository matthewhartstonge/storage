package mongo

import (
	// Standard Library Imports
	"testing"

	// Internal Imports
	"github.com/matthewhartstonge/storage"
)

func TestCacheMongoManager_ImplementsStorageConfigurer(t *testing.T) {
	c := &cacheMongoManager{}

	var i interface{} = c
	if _, ok := i.(storage.Configurer); !ok {
		t.Error("cacheMongoManager does not implement interface storage.Configurer")
	}
}

func TestCacheMongoManager_ImplementsStorageCacheStorer(t *testing.T) {
	c := &cacheMongoManager{}

	var i interface{} = c
	if _, ok := i.(storage.CacheStorer); !ok {
		t.Error("cacheMongoManager does not implement interface storage.CacheStorer")
	}
}

func TestCacheMongoManager_ImplementsStorageCacheManager(t *testing.T) {
	c := &cacheMongoManager{}

	var i interface{} = c
	if _, ok := i.(storage.CacheManager); !ok {
		t.Error("cacheMongoManager does not implement interface storage.CacheManager")
	}
}
