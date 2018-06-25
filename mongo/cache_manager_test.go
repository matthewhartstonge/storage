package mongo_test

import (
	// Standard Library Imports
	"testing"
	"time"

	// External Imports
	"github.com/ory/fosite"
	"github.com/pborman/uuid"

	// Public Imports
	"github.com/matthewhartstonge/storage"
)

var expected = storage.SessionCache{
	ID:         uuid.New(),
	CreateTime: time.Now().Unix(),
	UpdateTime: time.Now().Unix() + 600,
	Signature:  "Yhte@ensa#ei!+suu$re%sta^viik&oss*aha(joaisiaut)ta-is+ie%to_n==",
}

func TestCacheManager_Create(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	got, err := store.CacheManager.Create(ctx, storage.EntityCacheAccessTokens, expected)
	if err != nil {
		AssertError(t, err, nil, "create should return no database errors")
	}
	if got != expected {
		AssertError(t, got, expected, "cache object not equal")
	}
}

func TestCacheManager_Create_ShouldConflict(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	got, err := store.CacheManager.Create(ctx, storage.EntityCacheAccessTokens, expected)
	if err != nil {
		AssertError(t, err, nil, "create should return no database errors")
	}
	if got != expected {
		AssertError(t, got, expected, "cache object not equal")
	}

	got, err = store.CacheManager.Create(ctx, storage.EntityCacheAccessTokens, expected)
	if err == nil {
		AssertError(t, err, nil, "create duplicate should return conflict")
	}
	if err != storage.ErrResourceExists {
		AssertError(t, err, nil, "create duplicate should return conflict")
	}
}

func TestCacheManager_Get(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	created, err := store.CacheManager.Create(ctx, storage.EntityCacheAccessTokens, expected)
	if err != nil {
		AssertError(t, err, nil, "create should return no database errors")
	}
	if created != expected {
		AssertError(t, created, expected, "cache object not equal")
	}

	got, err := store.CacheManager.Get(ctx, storage.EntityCacheAccessTokens, expected.Key())
	if err != nil {
		AssertError(t, err, nil, "get should return no database errors")
	}
	if got != expected {
		AssertError(t, got, expected, "cache object not equal")
	}
}

func TestCacheManager_Get_ShouldReturnNotFound(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := fosite.ErrNotFound

	got, err := store.CacheManager.Get(ctx, storage.EntityCacheAccessTokens, "lolNotFound")
	if err != expected {
		AssertError(t, got, expected, "get should return not found")
	}
}
