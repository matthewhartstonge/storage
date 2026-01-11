package mongo_test

import (
	// Standard Library Imports
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	// External Imports
	"github.com/google/uuid"
	"github.com/ory/fosite"
	"go.mongodb.org/mongo-driver/bson"

	// Internal Imports
	"github.com/matthewhartstonge/storage"
	"github.com/matthewhartstonge/storage/mongo"
)

func expectedUser() storage.User {
	now := time.Now().UTC()
	return storage.User{
		ID:         uuid.NewString(),
		CreateTime: now.Unix(),
		UpdateTime: now.Unix() + 600,
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

func createUser(ctx context.Context, t *testing.T, store *mongo.Store, expected storage.User) storage.User {
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
		AssertError(t, got, expected, "user not equal")
		t.FailNow()
	}

	return expected
}

func TestUserManager_Create(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	createUser(ctx, t, store, expectedUser())
}

func TestUserManager_Create_ShouldConflict(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := createUser(ctx, t, store, expectedUser())
	_, err := store.UserManager.Create(ctx, expected)
	if err == nil {
		AssertError(t, err, nil, "create should return an error on conflict")
	}
	if !errors.Is(err, storage.ErrResourceExists) {
		AssertError(t, err, nil, "create should return conflict")
	}
}

func TestUserManager_Create_ShouldConflictOnUsername(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := createUser(ctx, t, store, expectedUser())
	expected.ID = uuid.NewString()
	_, err := store.UserManager.Create(ctx, expected)
	if err == nil {
		AssertError(t, err, nil, "create should return an error on conflict")
	}
	if !errors.Is(err, storage.ErrResourceExists) {
		AssertError(t, err, nil, "create should return conflict")
	}
}

func TestUserManager_Create_ShouldStoreCustomData(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expectedEntity := expectedUser()
	expectedData := expectedCustomData()

	// push in custom data
	if err := expectedEntity.Data.Marshal(expectedData); err != nil {
		AssertError(t, err, nil, "failed to unmarshal custom data")
	}

	// save user to mongo
	expectedEntity = createUser(ctx, t, store, expectedEntity)

	// extract entity directly
	query := bson.M{
		"id": expectedEntity.ID,
	}
	var gotEntity storage.User
	if err := store.DB.Collection(storage.EntityUsers).FindOne(ctx, query).Decode(&gotEntity); err != nil {
		AssertError(t, err, nil, "expected user to exist")
	}

	// Test expectations
	AssertUserCustomData(t, gotEntity, expectedEntity, &TestCustomData{}, &expectedData)
}

func TestUserManager_Get(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := createUser(ctx, t, store, expectedUser())
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
	if !errors.Is(err, expected) {
		AssertError(t, got, expected, "get should return not found")
	}
}

func TestUserManager_Get_ShouldReturnCustomData(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expectedEntity := expectedUser()
	expectedData := expectedCustomData()

	// push in custom data
	if err := expectedEntity.Data.Marshal(expectedData); err != nil {
		AssertError(t, err, nil, "failed to unmarshal custom data")
	}
	// save user to mongo
	expectedEntity = createUser(ctx, t, store, expectedEntity)

	// Get entity
	gotEntity, err := store.UserManager.Get(ctx, expectedEntity.ID)
	if err != nil {
		AssertError(t, err, nil, "failed to get user")
	}

	// Test expectations
	AssertUserCustomData(t, gotEntity, expectedEntity, &TestCustomData{}, &expectedData)
}

func TestUserManager_Update(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := createUser(ctx, t, store, expectedUser())
	// Perform an update...
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
	expected := createUser(ctx, t, store, expectedUser())
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

	_ = createUser(ctx, t, store, expectedUser())

	// Create 2nd user
	newUser := expectedUser()
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
	if !errors.Is(err, storage.ErrResourceExists) {
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
	if !errors.Is(err, fosite.ErrNotFound) {
		AssertError(t, err, nil, "update should return not found")
	}
}

func TestUserManager_Update_ShouldUpdateCustomData(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expectedEntity := expectedUser()
	expectedData := expectedCustomData()

	// push in custom data
	if err := expectedEntity.Data.Marshal(expectedData); err != nil {
		AssertError(t, err, nil, "failed to unmarshal custom data")
	}
	// save user to mongo
	expectedEntity = createUser(ctx, t, store, expectedEntity)

	// Update custom data
	expectedData.Contact.Name = "John Doe"
	if err := expectedEntity.Data.Marshal(expectedData); err != nil {
		AssertError(t, err, nil, "failed to unmarshal custom data")
	}

	gotEntity, err := store.UserManager.Update(ctx, expectedEntity.ID, expectedEntity)
	if err != nil {
		AssertError(t, err, nil, "failed to update user")
	}

	// Test expectations
	// override update time on expected with got. The time stamp received
	// should match time.Now().Unix() but due to the nature of time based
	// testing against time.Now().Unix(), it can fail on crossing over the
	// second boundary.
	expectedEntity.UpdateTime = gotEntity.UpdateTime
	AssertUserCustomData(t, gotEntity, expectedEntity, &TestCustomData{}, &expectedData)

	// Get record directly to verify struct passed back is persisting correctly
	// on update
	query := bson.M{
		"id": expectedEntity.ID,
	}
	gotEntity = storage.User{}
	if err := store.DB.Collection(storage.EntityUsers).FindOne(ctx, query).Decode(&gotEntity); err != nil {
		AssertError(t, err, nil, "expected user to exist")
	}

	// Test expectations
	AssertUserCustomData(t, gotEntity, expectedEntity, &TestCustomData{}, &expectedData)
}

func TestUserManager_Delete(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := createUser(ctx, t, store, expectedUser())

	err := store.UserManager.Delete(ctx, expected.ID)
	if err != nil {
		AssertError(t, err, nil, "delete should return no database errors")
	}

	// Double check that the original reference was deleted
	expectedErr := fosite.ErrNotFound
	got, err := store.UserManager.Get(ctx, expected.ID)
	if !errors.Is(expectedErr, err) {
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
	if !errors.Is(err, fosite.ErrNotFound) {
		AssertError(t, err, nil, "delete should return not found")
	}
}
