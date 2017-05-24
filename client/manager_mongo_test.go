package client_test

import (
	"github.com/MatthewHartstonge/storage"
	"github.com/MatthewHartstonge/storage/client"
	"github.com/ory/fosite"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var clientMongoDB = connectToMongo()
var secret = "foobarbaz"
var hash = generateHash(secret)

// connectToMongo generates a default mongo config and returns a connection to Mongo.
func connectToMongo() *client.MongoManager {
	cfg := storage.DefaultConfig()
	dbConnection, err := storage.NewDatastore(cfg)
	if err != nil {
		panic(err)
	}
	return &client.MongoManager{
		DB: dbConnection,
		Hasher: &fosite.BCrypt{
			WorkFactor: 10,
		},
	}
}

// setup creates a connection to Mongo.
func setup() {
	connectToMongo()
}

// teardown removes any left over created database and closes the underlying Mongo session.
func teardown() {
	clientMongoDB.DB.DropDatabase()
	clientMongoDB.DB.Session.Close()
}

// TestMain enables set up and teardown to ensure immutable test environments.
func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	teardown()
	os.Exit(retCode)
}

// setupTestCase resets the database to ensure idempotent tests and then returns a Teardown function which can be
// deferred.
func setupTestCase(t *testing.T) func(t *testing.T) {
	clientMongoDB.DB.DropDatabase()
	collection := clientMongoDB.DB.C("clients")
	c := expectedClient()
	err := collection.Insert(c)
	if err != nil {
		panic(err)
	}

	// Return the teardown case
	return func(t *testing.T) {
		clientMongoDB.DB.DropDatabase()
	}
}

// generateHash creates a single hash that wil be used for all tests.
func generateHash(pw string) string {
	h, err := clientMongoDB.Hasher.Hash([]byte(pw))
	if err != nil {
		panic(err)
	}
	return string(h)
}

// expectedClient returns an idempotent version of an expected client each time it's called.
func expectedClient() *client.Client {
	return &client.Client{
		ID:                "foo",
		Name:              "Foo bar App",
		Secret:            hash,
		RedirectURIs:      []string{"https://app.foo.example.com/callback", "https://dev-app.foo.example.com/callback"},
		GrantTypes:        []string{"client_credentials", "implicit"},
		ResponseTypes:     []string{"code", "token"},
		Scope:             "urn.foo.bar urn.foo.baz",
		Owner:             "FooBar Baz inc.",
		PolicyURI:         "https://foo.example.com/policy",
		TermsOfServiceURI: "https://foo.example.com/tos",
		ClientURI:         "https://app.foo.example.com/about",
		LogoURI:           "https://logos.example.com/happy-kitten.jpg",
		Contacts:          []string{"foo@example.com", "bar@example.com"},
		Public:            true,
	}
}

// TestClientManager_GetClientNotExist ensures that a error is raised if a client cannot be found by ID.
func TestClientManager_GetClientNotExist(t *testing.T) {
	teardown := setupTestCase(t)
	defer teardown(t)

	got, err := clientMongoDB.GetClient("notAnId")
	assert.Nil(t, got)
	assert.NotNil(t, err)
	assert.Error(t, err)
}

// TestClientManager_GetClient ensures that a client will be returned if the ID is found.
func TestClientManager_GetClient(t *testing.T) {
	teardown := setupTestCase(t)
	defer teardown(t)

	got, err := clientMongoDB.GetClient("foo")
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.ObjectsAreEqualValues(expectedClient(), got)
}

// TestMongoManager_UpdateClient ensures updating errors if a client can't be found with the provided ID
func TestMongoManager_UpdateClientNotExist(t *testing.T) {
	teardown := setupTestCase(t)
	defer teardown(t)

	expected := &client.Client{ID: "notanid", Name: "Updated1 Client Name"}

	err := clientMongoDB.UpdateClient(expected)
	assert.NotNil(t, err)
	assert.Error(t, err)
}

// TestMongoManager_UpdateClient ensures that a client will be updated
func TestMongoManager_UpdateClient(t *testing.T) {
	teardown := setupTestCase(t)
	defer teardown(t)

	expected := &client.Client{ID: "foo", Name: "Updated2 Client Name"}
	err := clientMongoDB.UpdateClient(expected)
	assert.Nil(t, err)

	// ensure update verifies against expected
	got, err := clientMongoDB.GetClient(expected.ID)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.ObjectsAreEqualValues(expected, got)
}

// TestMongoManager_UpdateClientDoesntRehashStoredPassword ensures the client secret doesn't get hashed if it has not
// been updated.
func TestMongoManager_UpdateClientDoesntRehashStoredPassword(t *testing.T) {
	teardown := setupTestCase(t)
	defer teardown(t)

	expected := &client.Client{ID: "foo", Name: "Updated3 Client Name"}
	err := clientMongoDB.UpdateClient(expected)
	assert.Nil(t, err)

	got, err := clientMongoDB.GetConcreteClient(expectedClient().ID)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.EqualValues(t, expectedClient().Secret, got.Secret)
}

// TestMongoManager_UpdateClientRehashsStoredPasswordIfUpdated ensures a new password gets correctly hashed if updated
func TestMongoManager_UpdateClientRehashsStoredPasswordIfUpdated(t *testing.T) {
	teardown := setupTestCase(t)
	defer teardown(t)

	// Update password
	newSecret := "bazbarfoo"
	expected := &client.Client{ID: "foo", Name: "Updated4 Client Name", Secret: newSecret}
	err := clientMongoDB.UpdateClient(expected)
	assert.Nil(t, err)

	// Obtain update from database
	got, err := clientMongoDB.GetConcreteClient(expectedClient().ID)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.NotEqual(t, expectedClient().Secret, got.Secret)

	// Check that new password compares with the new hash stored in the database
	err = clientMongoDB.Hasher.Compare(got.GetHashedSecret(), []byte(newSecret))
	assert.Nil(t, err)
}

// TestMongoManager_AuthenticateNotExist ensures Authenticate errors if a client can't be found with the provided ID
func TestMongoManager_AuthenticateNotExist(t *testing.T) {
	teardown := setupTestCase(t)
	defer teardown(t)

	expected := &client.Client{ID: "notanid", Name: "Updated3 Client Name"}

	_, err := clientMongoDB.Authenticate(expected.ID, []byte(secret))
	assert.NotNil(t, err)
	assert.Error(t, err)
}

// TestMongoManager_AuthenticateNoMatch ensures Authenticate errors if the secret doesn't match
func TestMongoManager_AuthenticateNoMatch(t *testing.T) {
	teardown := setupTestCase(t)
	defer teardown(t)

	secret := []byte("notreallythepassword")
	c, err := clientMongoDB.Authenticate(expectedClient().ID, secret)
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.Error(t, err)
}

// TestMongoManager_Authenticate ensures that a client can authenticate successfully
func TestMongoManager_Authenticate(t *testing.T) {
	teardown := setupTestCase(t)
	defer teardown(t)

	expected := expectedClient()

	got, err := clientMongoDB.Authenticate(expected.ID, []byte(secret))
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.ObjectsAreEqualValues(expected, got)
}
