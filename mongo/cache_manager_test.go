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

func expectedSessionCache() storage.SessionCache {
	return storage.SessionCache{
		ID:         uuid.New(),
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix() + 600,
		Signature:  "Yhte@ensa#ei!+suu$re%sta^viik&oss*aha(joaisiaut)ta-is+ie%to_n==",
	}
}

func TestCacheManager_Create(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := expectedSessionCache()
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

	expected := expectedSessionCache()
	got, err := store.CacheManager.Create(ctx, storage.EntityCacheAccessTokens, expected)
	if err != nil {
		AssertError(t, err, nil, "create should return no database errors")
	}
	if got != expected {
		AssertError(t, got, expected, "cache object not equal")
	}

	got, err = store.CacheManager.Create(ctx, storage.EntityCacheAccessTokens, expected)
	if err == nil {
		AssertError(t, err, nil, "create should return an error on conflict")
	}
	if err != storage.ErrResourceExists {
		AssertError(t, err, nil, "create should return conflict")
	}
}

func TestCacheManager_Get(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := expectedSessionCache()
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

func TestCacheManager_Update(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := expectedSessionCache()
	created, err := store.CacheManager.Create(ctx, storage.EntityCacheAccessTokens, expected)
	if err != nil {
		AssertError(t, err, nil, "create should return no database errors")
	}
	if created != expected {
		AssertError(t, created, expected, "cache object not equal")
	}

	// Perform an update..
	updatedSignature := "something completely different!"
	created.Signature = updatedSignature

	got, err := store.CacheManager.Update(ctx, storage.EntityCacheAccessTokens, created)
	if err != nil {
		AssertError(t, err, nil, "update should return no database errors")
	}
	if expected.UpdateTime == 0 {
		AssertError(t, got.UpdateTime, time.Now().Unix(), "update time was not set")
	}

	// override update time on expected with got. The time stamp received
	// should match time.Now().Unix() but due to the nature of time based
	// testing against time.Now().Unix(), it can fail on crossing over the
	// second boundary.
	expected.UpdateTime = got.UpdateTime
	expected.Signature = updatedSignature
	if got != expected {
		AssertError(t, got, expected, "cache update object not equal")
	}
}

func TestCacheManager_Update_ShouldReturnNotFound(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	_, err := store.CacheManager.Update(ctx, storage.EntityCacheAccessTokens, expectedSessionCache())
	if err == nil {
		AssertError(t, err, nil, "update should return an error on not found")
	}
	if err != fosite.ErrNotFound {
		AssertError(t, err, nil, "update should return not found")
	}
}

func TestCacheManager_Delete(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := expectedSessionCache()
	created, err := store.CacheManager.Create(ctx, storage.EntityCacheAccessTokens, expected)
	if err != nil {
		AssertError(t, err, nil, "create should return no database errors")
	}
	if created != expected {
		AssertError(t, created, expected, "cache object not equal")
	}

	err = store.CacheManager.Delete(ctx, storage.EntityCacheAccessTokens, created.Key())
	if err != nil {
		AssertError(t, err, nil, "delete should return no database errors")
	}

	// Double check that the original reference was deleted
	expectedErr := fosite.ErrNotFound
	got, err := store.CacheManager.Get(ctx, storage.EntityCacheAccessTokens, created.Key())
	if err != expectedErr {
		AssertError(t, got, expectedErr, "get should return not found")
	}
}

func TestCacheManager_Delete_ShouldReturnNotFound(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := expectedSessionCache()
	err := store.CacheManager.Delete(ctx, storage.EntityCacheAccessTokens, expected.Key())
	if err == nil {
		AssertError(t, err, nil, "delete should return an error on not found")
	}
	if err != fosite.ErrNotFound {
		AssertError(t, err, nil, "delete should return not found")
	}
}

func TestCacheManager_DeleteByValue(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := expectedSessionCache()
	created, err := store.CacheManager.Create(ctx, storage.EntityCacheAccessTokens, expected)
	if err != nil {
		AssertError(t, err, nil, "create should return no database errors")
	}
	if created != expected {
		AssertError(t, created, expected, "cache object not equal")
	}

	err = store.CacheManager.DeleteByValue(ctx, storage.EntityCacheAccessTokens, created.Value())
	if err != nil {
		AssertError(t, err, nil, "DeleteByValue should return no database errors")
	}

	// Double check that the original reference was deleted
	expectedErr := fosite.ErrNotFound
	got, err := store.CacheManager.Get(ctx, storage.EntityCacheAccessTokens, created.Key())
	if err != expectedErr {
		AssertError(t, got, expectedErr, "get should return not found")
	}
}

func TestCacheManager_DeleteByValue_ShouldReturnNotFound(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := expectedSessionCache()
	err := store.CacheManager.DeleteByValue(ctx, storage.EntityCacheAccessTokens, expected.Value())
	if err == nil {
		AssertError(t, err, nil, "DeleteByValue should return an error on not found")
	}
	if err != fosite.ErrNotFound {
		AssertError(t, err, nil, "DeleteByValue should return not found")
	}
}
