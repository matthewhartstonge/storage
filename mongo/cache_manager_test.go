package mongo

import (
	// Standard Library Imports
	"testing"

	// Internal Imports
	"github.com/matthewhartstonge/storage"
)

func TestCacheMongoManager_ImplementsStorageCacheStorer(t *testing.T) {
	c := &cacheMongoManager{}

	var i interface{} = c
	_, ok := i.(storage.CacheStorer)
	if ok != true {
		t.Error("cacheMongoManager does not implement interface storage.CacheStorer")
	}
}

func TestCacheMongoManager_ImplementsStorageCacheManager(t *testing.T) {
	c := &cacheMongoManager{}

	var i interface{} = c
	_, ok := i.(storage.CacheManager)
	if ok != true {
		t.Error("cacheMongoManager does not implement interface storage.CacheManager")
	}
}
