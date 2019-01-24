package mongo_test

import (
	// Standard Library Imports
	"context"
	"reflect"
	"testing"
	"time"

	// External Imports
	"github.com/ory/fosite"
	"github.com/pborman/uuid"

	// Internal Imports
	"github.com/matthewhartstonge/storage"
	"github.com/matthewhartstonge/storage/mongo"
)

func expectedClient() storage.Client {
	return storage.Client{
		ID:         uuid.New(),
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix() + 600,
		AllowedAudiences: []string{
			uuid.New(),
			uuid.New(),
		},
		AllowedTenantAccess: []string{
			uuid.New(),
			uuid.New(),
		},
		GrantTypes: []string{
			string(fosite.AccessToken),
			string(fosite.RefreshToken),
			string(fosite.AuthorizeCode),
			string(fosite.IDToken),
		},
		ResponseTypes: []string{
			"code",
			"token",
		},
		Scopes: []string{
			"urn:test:cats:write",
			"urn:test:dogs:read",
		},
		Public:   true,
		Disabled: false,
		Name:     "Test Client",
		Secret:   "foobar",
		RedirectURIs: []string{
			"https://test.example.com",
		},
		Owner:             "Widgets Inc.",
		PolicyURI:         "https://test.example.com/policy",
		TermsOfServiceURI: "https://test.example.com/tos",
		ClientURI:         "https://app.example.com",
		LogoURI:           "https://app.example.com/favicon-128x128.png",
		Contacts: []string{
			"John Doe <j.doe@example.com>",
		},
	}
}

func createClient(t *testing.T, ctx context.Context, store *mongo.Store) storage.Client {
	expected := expectedClient()
	got, err := store.ClientManager.Create(ctx, expected)
	if err != nil {
		AssertError(t, err, nil, "create should return no database errors")
		t.FailNow()
	}

	if got.Secret == "" || got.Secret == expected.Secret {
		AssertError(t, got.Secret, "bcrypt encoded secret", "create should hash the secret")
		t.FailNow()
	}

	expected.Secret = got.Secret
	if !reflect.DeepEqual(got, expected) {
		AssertError(t, got, expected, "client not equal")
		t.FailNow()
	}

	return expected
}

func TestClientManager_Create(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	createClient(t, ctx, store)
}

func TestClientManager_Create_ShouldConflict(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := createClient(t, ctx, store)
	_, err := store.ClientManager.Create(ctx, expected)
	if err == nil {
		AssertError(t, err, nil, "create should return an error on conflict")
	}
	if err != storage.ErrResourceExists {
		AssertError(t, err, nil, "create should return conflict")
	}
}

func TestClientManager_Get(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := createClient(t, ctx, store)
	got, err := store.ClientManager.Get(ctx, expected.ID)
	if err != nil {
		AssertError(t, err, nil, "get should return no database errors")
	}
	if !reflect.DeepEqual(got, expected) {
		AssertError(t, got, expected, "client not equal")
	}
}

func TestClientManager_Get_ShouldReturnNotFound(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := fosite.ErrNotFound
	got, err := store.ClientManager.Get(ctx, "lolNotFound")
	if err != expected {
		AssertError(t, got, expected, "get should return not found")
	}
}

func TestClientManager_Update(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := createClient(t, ctx, store)
	// Perform an update..
	expected.Name = "something completely different!"

	got, err := store.ClientManager.Update(ctx, expected.ID, expected)
	if err != nil {
		AssertError(t, err, nil, "update should return no database errors")
	}

	if expected.UpdateTime == 0 {
		AssertError(t, got.UpdateTime, time.Now().Unix(), "update time was not set")
	}

	if expected.Secret != got.Secret {
		AssertError(t, got.Secret, expected.Secret, "secret should not change on update unless explicitly changed")
	}

	// override update time on expected with got. The time stamp received
	// should match time.Now().Unix() but due to the nature of time based
	// testing against time.Now().Unix(), it can fail on crossing over the
	// second boundary.
	expected.UpdateTime = got.UpdateTime
	if !reflect.DeepEqual(got, expected) {
		AssertError(t, got, expected, "client update object not equal")
	}
}

func TestClientManager_Update_ShouldChangePassword(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	newSecret := "s0methingElse!"
	expected := createClient(t, ctx, store)
	oldHash := expected.Secret

	// Perform a password update..
	expected.Secret = newSecret

	got, err := store.ClientManager.Update(ctx, expected.ID, expected)
	if err != nil {
		AssertError(t, err, nil, "update should return no database errors")
	}

	if expected.UpdateTime == 0 {
		AssertError(t, got.UpdateTime, time.Now().Unix(), "update time was not set")
	}

	if got.Secret == oldHash {
		AssertError(t, got.Secret, "new bcrypt hash", "secret was not updated")
	}

	if got.Secret == newSecret {
		AssertError(t, got.Secret, "new bcrypt hash", "secret was not hashed")
	}

	// Should authenticate against the new hash
	if err := store.Hasher.Compare(ctx, got.GetHashedSecret(), []byte(newSecret)); err != nil {
		AssertError(t, got.Secret, "bcrypt authenticate-able hash", "unable to authenticate with updated hash")
	}

	// override update time on expected with got. The time stamp received
	// should match time.Now().Unix() but due to the nature of time based
	// testing against time.Now().Unix(), it can fail on crossing over the
	// second boundary.
	expected.UpdateTime = got.UpdateTime
	// override expected secret as the assertions have passed above.
	expected.Secret = got.Secret

	if !reflect.DeepEqual(got, expected) {
		AssertError(t, got, expected, "client update object not equal")
	}
}

func TestClientManager_Update_ShouldReturnNotFound(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	_, err := store.ClientManager.Update(ctx, uuid.New(), expectedClient())
	if err == nil {
		AssertError(t, err, nil, "update should return an error on not found")
	}
	if err != fosite.ErrNotFound {
		AssertError(t, err, nil, "update should return not found")
	}
}

func TestClientManager_Delete(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := createClient(t, ctx, store)

	err := store.ClientManager.Delete(ctx, expected.ID)
	if err != nil {
		AssertError(t, err, nil, "delete should return no database errors")
	}

	// Double check that the original reference was deleted
	expectedErr := fosite.ErrNotFound
	got, err := store.ClientManager.Get(ctx, expected.ID)
	if err != expectedErr {
		AssertError(t, got, expectedErr, "get should return not found")
	}
}

func TestClientManager_Delete_ShouldReturnNotFound(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	err := store.ClientManager.Delete(ctx, expectedClient().ID)
	if err == nil {
		AssertError(t, err, nil, "delete should return an error on not found")
	}
	if err != fosite.ErrNotFound {
		AssertError(t, err, nil, "delete should return not found")
	}
}
