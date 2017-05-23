package storage_test

import (
	"github.com/MatthewHartstonge/storage"
	"github.com/MatthewHartstonge/storage/model"
	"github.com/ory/fosite"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var clientMongoDB = connectToMongo()
var expectedClient = &model.Client{
	ID:                "foo",
	Name:              "Foo bar App",
	Secret:            "foobarsecretbaz",
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

// connectToMongo generates a default mongo config and returns a connection to Mongo.
func connectToMongo() *storage.ClientManager {
	cfg := storage.DefaultConfig()
	dbConnection, err := storage.NewDatastore(cfg)
	if err != nil {
		panic(err)
	}
	return &storage.ClientManager{
		&fosite.BCrypt{},
		dbConnection,
	}
}

// setup creates a connection to Mongo, creates a database, collection and an expected client in the database.
func setup() {
	connectToMongo()
	collection := clientMongoDB.DB.C("clients")
	err := collection.Insert(expectedClient)
	if err != nil {
		panic(err)
	}
}

// teardown removes the created database and closes the underlying Mongo session.
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

// TestClientManager_GetClientNotExist ensures that a error is raised if a client cannot be found by ID.
func TestClientManager_GetClientNotExist(t *testing.T) {
	got, err := clientMongoDB.GetClient("notAnId")
	assert.Nil(t, got)
	assert.NotNil(t, err)
	assert.Error(t, err)
}

// TestClientManager_GetClient ensures that a client will be returned if the ID is found.
func TestClientManager_GetClient(t *testing.T) {
	expected := &model.Client{
		ID:                "foo",
		Name:              "Foo bar App",
		Secret:            "foobarsecretbaz",
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
	got, err := clientMongoDB.GetClient("foo")
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.EqualValues(t, expected, got)
}
