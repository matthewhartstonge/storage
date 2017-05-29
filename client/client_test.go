package client_test

import (
	"github.com/MatthewHartstonge/storage"
	"github.com/MatthewHartstonge/storage/client"
	"github.com/MatthewHartstonge/storage/mongo"
	"github.com/ory/fosite"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var ClientMongoDB = ConnectToMongo()
var Secret = "foobarbaz"
var Hash = GenerateHash(Secret)

// ConnectToMongo generates a default mongo config and returns a connection to Mongo.
func ConnectToMongo() *client.MongoManager {
	cfg := storage.DefaultConfig()
	db, err := storage.ConnectToMongo(cfg)
	if err != nil {
		panic(err)
	}
	return &client.MongoManager{
		DB: db,
		Hasher: &fosite.BCrypt{
			WorkFactor: 10,
		},
	}
}

// Setup creates a connection to Mongo.
func Setup() {
	ConnectToMongo()
}

// teardown removes any left over created database and closes the underlying Mongo session.
func Teardown() {
	ClientMongoDB.DB.DropDatabase()
	ClientMongoDB.DB.Session.Close()
}

// TestMain enables set up and teardown to ensure immutable test environments.
func TestMain(m *testing.M) {
	Setup()
	retCode := m.Run()
	Teardown()
	os.Exit(retCode)
}

// SetupTestCase resets the database to ensure idempotent tests and then returns a Teardown function which can be
// deferred.
func SetupTestCase(t *testing.T) func(t *testing.T) {
	ClientMongoDB.DB.DropDatabase()
	collection := ClientMongoDB.DB.C(mongo.CollectionClients)
	c := expectedClient()
	err := collection.Insert(c)
	if err != nil {
		panic(err)
	}

	// Return the teardown case
	return func(t *testing.T) {
		ClientMongoDB.DB.DropDatabase()
	}
}

// GenerateHash creates a single Hash that wil be used for all tests.
func GenerateHash(pw string) string {
	h, err := ClientMongoDB.Hasher.Hash([]byte(pw))
	if err != nil {
		panic(err)
	}
	return string(h)
}

// TestClient ensures that Client conforms to fosite interfaces and that inputs and outputs are formed correctly.
func TestClient(t *testing.T) {
	c := &client.Client{
		ID:           "foo",
		RedirectURIs: []string{"foo"},
		Scopes:       []string{"foo", "bar"},
	}

	assert.EqualValues(t, "foo", c.GetID())
	assert.EqualValues(t, c.RedirectURIs, c.GetRedirectURIs())
	assert.EqualValues(t, []byte(c.Secret), c.GetHashedSecret())
	assert.EqualValues(t, fosite.Arguments{"authorization_code"}, c.GetGrantTypes())
	assert.EqualValues(t, fosite.Arguments{"code"}, c.GetResponseTypes())
	assert.EqualValues(t, (c.Owner), c.GetOwner())
	assert.EqualValues(t, (c.Public), c.IsPublic())
	assert.Len(t, c.GetScopes(), 2)
	assert.EqualValues(t, c.RedirectURIs, c.GetRedirectURIs())

	// fosite.Argument logic branches
	expectedGrantTypes := fosite.Arguments{"foo", "bar"}
	c.GrantTypes = []string{"foo", "bar"}
	assert.EqualValues(t, expectedGrantTypes, c.GetGrantTypes())

	expectedResponseTypes := fosite.Arguments{"bar", "foo"}
	c.ResponseTypes = []string{"bar", "foo"}
	assert.EqualValues(t, expectedResponseTypes, c.GetResponseTypes())
}
