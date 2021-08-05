package mongo_test

import (
	// Standard Library Imports
	"context"
	"reflect"
	"testing"
	"time"

	// External Imports
	"github.com/google/uuid"
	"github.com/ory/fosite"

	// Internal Imports
	"github.com/matthewhartstonge/storage"
	"github.com/matthewhartstonge/storage/mongo"
)

func expectedUser() storage.User {
	return storage.User{
		ID:         uuid.NewString(),
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix() + 600,
		AllowedTenantAccess: []string{
			uuid.NewString(),
			uuid.NewString(),
		},
		AllowedPersonAccess: []string{
			uuid.NewString(),
			uuid.NewString(),
		},
		Scopes: []string{
			"urn:test:cats:write",
			"urn:test:dogs:read",
		},
		Roles: []string{
			"user",
			"printer",
		},
		PersonID:   uuid.NewString(),
		Disabled:   false,
		Username:   "j.doe@example.com",
		Password:   "foobar",
		FirstName:  "John",
		LastName:   "Doe",
		ProfileURI: "https://profiles.example.com/j.doe@example.com",
	}
}

func createUser(ctx context.Context, t *testing.T, store *mongo.Store) storage.User {
	expected := expectedUser()
	got, err := store.UserManager.Create(ctx, expected)
	if err != nil {
		AssertError(t, err, nil, "create should return no database errors")
		t.FailNow()
	}

	if got.Password == "" || got.Password == expected.Password {
		AssertError(t, got.Password, "bcrypt encoded secret", "create should hash the secret")
		t.FailNow()
	}

	expected.Password = got.Password
	if !reflect.DeepEqual(got, expected) {
		AssertError(t, got, expected, "client not equal")
		t.FailNow()
	}

	return expected
}

func TestUserManager_Create(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	createUser(ctx, t, store)
}

func TestUserManager_Create_ShouldConflict(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := createUser(ctx, t, store)
	_, err := store.UserManager.Create(ctx, expected)
	if err == nil {
		AssertError(t, err, nil, "create should return an error on conflict")
	}
	if err != storage.ErrResourceExists {
		AssertError(t, err, nil, "create should return conflict")
	}
}

func TestUserManager_Create_ShouldConflictOnUsername(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := createUser(ctx, t, store)
	expected.ID = uuid.NewString()
	_, err := store.UserManager.Create(ctx, expected)
	if err == nil {
		AssertError(t, err, nil, "create should return an error on conflict")
	}
	if err != storage.ErrResourceExists {
		AssertError(t, err, nil, "create should return conflict")
	}
}

func TestUserManager_Get(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := createUser(ctx, t, store)
	got, err := store.UserManager.Get(ctx, expected.ID)
	if err != nil {
		AssertError(t, err, nil, "get should return no database errors")
	}
	if !reflect.DeepEqual(got, expected) {
		AssertError(t, got, expected, "user not equal")
	}
}

func TestUserManager_Get_ShouldReturnNotFound(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := fosite.ErrNotFound
	got, err := store.UserManager.Get(ctx, "lolNotFound")
	if err != expected {
		AssertError(t, got, expected, "get should return not found")
	}
}

func TestUserManager_Update(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := createUser(ctx, t, store)
	// Perform an update..
	expected.FirstName = "Bob"
	expected.LastName = "Marley"
	expected.Username = "b.marley@example.com"
	expected.ProfileURI = "https://profiles.example.com/"

	got, err := store.UserManager.Update(ctx, expected.ID, expected)
	if err != nil {
		AssertError(t, err, nil, "update should return no database errors")
	}
	if expected.UpdateTime == 0 {
		AssertError(t, got.UpdateTime, time.Now().Unix(), "update time was not set")
	}

	if expected.Password != got.Password {
		AssertError(t, got.Password, expected.Password, "password should not change on update unless explicitly changed")
	}

	// override update time on expected with got. The time stamp received
	// should match time.Now().Unix() but due to the nature of time based
	// testing against time.Now().Unix(), it can fail on crossing over the
	// second boundary.
	expected.UpdateTime = got.UpdateTime
	if !reflect.DeepEqual(got, expected) {
		AssertError(t, got, expected, "user update object not equal")
	}
}

func TestUserManager_Update_ShouldChangePassword(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	newPassword := "s0methingElse!"
	expected := createUser(ctx, t, store)
	oldHash := expected.Password

	// Perform a password update..
	expected.Password = newPassword

	got, err := store.UserManager.Update(ctx, expected.ID, expected)
	if err != nil {
		AssertError(t, err, nil, "update should return no database errors")
	}

	if expected.UpdateTime == 0 {
		AssertError(t, got.UpdateTime, time.Now().Unix(), "update time was not set")
	}

	if got.Password == oldHash {
		AssertError(t, got.Password, "new bcrypt hash", "password was not updated")
	}

	if got.Password == newPassword {
		AssertError(t, got.Password, "new bcrypt hash", "password was not hashed")
	}

	// Should authenticate against the new hash
	if err := got.Authenticate(newPassword, store.Hasher); err != nil {
		AssertError(t, got.Password, "bcrypt authenticate-able hash", "unable to authenticate with updated hash")
	}

	// override update time on expected with got. The time stamp received
	// should match time.Now().Unix() but due to the nature of time based
	// testing against time.Now().Unix(), it can fail on crossing over the
	// second boundary.
	expected.UpdateTime = got.UpdateTime
	// override expected password as the assertions have passed above.
	expected.Password = got.Password

	if !reflect.DeepEqual(got, expected) {
		AssertError(t, got, expected, "user update object not equal")
	}
}

func TestUserManager_Update_ShouldConflictUsername(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	user := createUser(ctx, t, store)

	// Create 2nd user
	newUser := user
	newUser.ID = uuid.NewString()
	newUser.FirstName = "Bob"
	newUser.LastName = "Marley"
	newUser.Password = "barbaz"
	newUser.Username = "b.marley@example.com"
	newUser.ProfileURI = "https://profiles.example.com/"

	newUser, err := store.UserManager.Create(ctx, newUser)
	if err != nil {
		AssertError(t, err, nil, "create should return no database errors")
		t.FailNow()
	}

	// Perform an update where the username matches an existing username..
	newUser.Username = "j.doe@example.com"

	_, err = store.UserManager.Update(ctx, newUser.ID, newUser)
	if err == nil {
		AssertError(t, err, nil, "update should return an error on username conflict")
	}
	if err != storage.ErrResourceExists {
		AssertError(t, err, nil, "update should return conflict on username")
	}
}

func TestUserManager_Update_ShouldReturnNotFound(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	_, err := store.UserManager.Update(ctx, uuid.NewString(), expectedUser())
	if err == nil {
		AssertError(t, err, nil, "update should return an error on not found")
	}
	if err != fosite.ErrNotFound {
		AssertError(t, err, nil, "update should return not found")
	}
}

func TestUserManager_Delete(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := createUser(ctx, t, store)

	err := store.UserManager.Delete(ctx, expected.ID)
	if err != nil {
		AssertError(t, err, nil, "delete should return no database errors")
	}

	// Double check that the original reference was deleted
	expectedErr := fosite.ErrNotFound
	got, err := store.UserManager.Get(ctx, expected.ID)
	if err != expectedErr {
		AssertError(t, got, expectedErr, "get should return not found")
	}
}

func TestUserManager_Delete_ShouldReturnNotFound(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	err := store.UserManager.Delete(ctx, expectedUser().ID)
	if err == nil {
		AssertError(t, err, nil, "delete should return an error on not found")
	}
	if err != fosite.ErrNotFound {
		AssertError(t, err, nil, "delete should return not found")
	}
}
