package client_test

import (
	"testing"

	"github.com/ory/fosite"
	"github.com/stretchr/testify/assert"

	"github.com/matthewhartstonge/storage/client"
)

// expectedClient returns an idempotent version of an expected client each time it's called.
func expectedClient() *client.Client {
	return &client.Client{
		ID:                  "foo",
		Name:                "Foo bar App",
		AllowedTenantAccess: []string{"29c78d37-a555-4d90-a038-bdb67a82b461", "5253ee1a-aaac-49b1-ab7c-85b6d0571366"},
		Secret:              []byte(Hash),
		RedirectURIs:        []string{"https://app.foo.example.com/callback", "https://dev-app.foo.example.com/callback"},
		GrantTypes:          []string{"client_credentials", "implicit"},
		ResponseTypes:       []string{"code", "token"},
		Scopes:              []string{"cats:read", "cats:delete"},
		Owner:               "FooBar Baz inc.",
		PolicyURI:           "https://foo.example.com/policy",
		TermsOfServiceURI:   "https://foo.example.com/tos",
		ClientURI:           "https://app.foo.example.com/about",
		LogoURI:             "https://logos.example.com/happy-kitten.jpg",
		Contacts:            []string{"foo@example.com", "bar@example.com"},
		Public:              true,
		Disabled:            false,
	}
}

func TestClientMongoManager_ImplementsFositeClientManagerInterface(t *testing.T) {
	c := &client.MongoManager{}

	var i interface{} = c
	_, ok := i.(fosite.ClientManager)
	assert.Equal(t, true, ok, "client.MongoManager does not implement interface fosite.ClientManager")
}

// TestClientManager_GetClientNotExist ensures that a error is raised if a client cannot be found by ID.
func TestClientManager_GetClientNotExist(t *testing.T) {
	teardown := SetupTestCase(t)
	defer teardown(t)

	got, err := ClientMongoDB.GetClient(nil, "notAnId")
	assert.Nil(t, got)
	assert.NotNil(t, err)
	assert.Error(t, err)
}

// TestClientManager_GetClient ensures that a client will be returned if the ID is found.
func TestClientManager_GetClient(t *testing.T) {
	teardown := SetupTestCase(t)
	defer teardown(t)

	got, err := ClientMongoDB.GetClient(nil, "foo")
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.ObjectsAreEqualValues(expectedClient(), got)
}

// TestMongoManager_UpdateClient ensures updating errors if a client can't be found with the provided ID
func TestMongoManager_UpdateClientNotExist(t *testing.T) {
	teardown := SetupTestCase(t)
	defer teardown(t)

	expected := &client.Client{ID: "notanid", Name: "Updated1 Client Name"}

	err := ClientMongoDB.UpdateClient(expected)
	assert.NotNil(t, err)
	assert.Error(t, err)
}

// TestMongoManager_UpdateClient ensures that a client will be updated
func TestMongoManager_UpdateClient(t *testing.T) {
	teardown := SetupTestCase(t)
	defer teardown(t)

	expected := &client.Client{ID: "foo", Name: "Updated2 Client Name"}
	err := ClientMongoDB.UpdateClient(expected)
	assert.Nil(t, err)

	// ensure update verifies against expected
	got, err := ClientMongoDB.GetClient(nil, expected.ID)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.ObjectsAreEqualValues(expected, got)
}

// TestMongoManager_UpdateClientDoesntRehashStoredPassword ensures the client Secret doesn't get hashed if it has not
// been updated.
func TestMongoManager_UpdateClientDoesntRehashStoredPassword(t *testing.T) {
	teardown := SetupTestCase(t)
	defer teardown(t)

	expected := &client.Client{ID: "foo", Name: "Updated3 Client Name"}
	err := ClientMongoDB.UpdateClient(expected)
	assert.Nil(t, err)

	got, err := ClientMongoDB.GetConcreteClient(expectedClient().ID)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.EqualValues(t, expectedClient().Secret, got.Secret)
}

// TestMongoManager_UpdateClientRehashsStoredPasswordIfUpdated ensures a new password gets correctly hashed if updated
func TestMongoManager_UpdateClientRehashsStoredPasswordIfUpdated(t *testing.T) {
	teardown := SetupTestCase(t)
	defer teardown(t)

	// Update password
	newSecret := []byte("bazbarfoo")
	expected := &client.Client{ID: "foo", Name: "Updated4 Client Name", Secret: newSecret}
	err := ClientMongoDB.UpdateClient(expected)
	assert.Nil(t, err)

	// Obtain update from database
	got, err := ClientMongoDB.GetConcreteClient(expectedClient().ID)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.NotEqual(t, expectedClient().Secret, got.Secret)

	// Check that new password compares with the new Hash stored in the database
	err = ClientMongoDB.Hasher.Compare(got.GetHashedSecret(), []byte(newSecret))
	assert.Nil(t, err)
}

// TestMongoManager_AuthenticateNotExist ensures Authenticate errors if a client can't be found with the provided ID
func TestMongoManager_AuthenticateNotExist(t *testing.T) {
	teardown := SetupTestCase(t)
	defer teardown(t)

	expected := &client.Client{ID: "notanid", Name: "Updated3 Client Name"}

	_, err := ClientMongoDB.Authenticate(expected.ID, []byte(Secret))
	assert.NotNil(t, err)
	assert.Error(t, err)
}

// TestMongoManager_AuthenticateNoMatch ensures Authenticate errors if the Secret doesn't match
func TestMongoManager_AuthenticateNoMatch(t *testing.T) {
	teardown := SetupTestCase(t)
	defer teardown(t)

	secret := []byte("notreallythepassword")
	c, err := ClientMongoDB.Authenticate(expectedClient().ID, secret)
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.Error(t, err)
}

// TestMongoManager_Authenticate ensures that a client can authenticate successfully
func TestMongoManager_Authenticate(t *testing.T) {
	teardown := SetupTestCase(t)
	defer teardown(t)

	expected := expectedClient()
	got, err := ClientMongoDB.Authenticate(expected.ID, []byte(Secret))
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.ObjectsAreEqualValues(expected, got)
}

// TODO: Unit tests for CreateClient
// TODO: Unit tests for DeleteClient
// TODO: Unit tests for GetClients
