package client_test

import (
	"os"
	"testing"

	"github.com/matthewhartstonge/storage"
	"github.com/matthewhartstonge/storage/client"
	"github.com/matthewhartstonge/storage/mongo"
	"github.com/ory/fosite"
	"github.com/stretchr/testify/assert"
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
func SetupTestCase(_ *testing.T) func(t *testing.T) {
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

func TestClient_EnableScopeAccess_None(t *testing.T) {
	u := expectedClient()

	expectedScopes := []string{
		"cats:read",
		"cats:delete",
	}

	u.EnableScopeAccess("cats:read")
	assert.EqualValues(t, expectedScopes, u.Scopes)

	u.EnableScopeAccess("cats:delete")
	assert.EqualValues(t, expectedScopes, u.Scopes)
}

func TestClient_EnableScopeAccess_One(t *testing.T) {
	u := expectedClient()

	expectedScopes := []string{
		"cats:read",
		"cats:delete",
		"cats:hug",
	}

	u.EnableScopeAccess("cats:hug")
	assert.EqualValues(t, expectedScopes, u.Scopes)

	u.EnableScopeAccess("cats:hug")
	assert.EqualValues(t, expectedScopes, u.Scopes)

	u.EnableScopeAccess("cats:read")
	assert.EqualValues(t, expectedScopes, u.Scopes)
}

func TestClient_EnableScopeAccess_Many(t *testing.T) {
	u := expectedClient()

	expectedScopes := []string{
		"cats:read",
		"cats:delete",
		"cats:hug",
		"cats:purr",
		"cats:meow",
	}

	u.EnableScopeAccess("cats:hug", "cats:purr", "cats:meow")
	assert.EqualValues(t, expectedScopes, u.Scopes)

	u.EnableScopeAccess("cats:hug", "cats:purr", "cats:meow")
	assert.EqualValues(t, expectedScopes, u.Scopes)
}

func TestClient_DisableScopeAccess_None(t *testing.T) {
	u := expectedClient()

	expectedScopes := []string{
		"cats:read",
		"cats:delete",
	}

	u.DisableScopeAccess("cats:hug")
	assert.EqualValues(t, expectedScopes, u.Scopes)
}

func TestClient_DisableScopeAccess_One(t *testing.T) {
	u := expectedClient()
	expectedScopes := []string{
		"cats:delete",
	}

	u.DisableScopeAccess("cats:read")
	assert.EqualValues(t, expectedScopes, u.Scopes)

	u.DisableScopeAccess("cats:read")
	assert.EqualValues(t, expectedScopes, u.Scopes)

	u.DisableScopeAccess("cats:delete")
	assert.EqualValues(t, expectedScopes[:len(expectedScopes)-1], u.Scopes)

	u.DisableScopeAccess("cats:read")
	assert.EqualValues(t, expectedScopes[:len(expectedScopes)-1], u.Scopes)

	u.DisableScopeAccess("cats:mug")
	assert.EqualValues(t, expectedScopes[:len(expectedScopes)-1], u.Scopes)
}

func TestClient_DisableScopeAccess_Many(t *testing.T) {
	u := expectedClient()
	expectedScopes := []string{
		"cats:read",
	}

	u.Scopes = []string{
		"cats:read",
		"cats:delete",
		"cats:hug",
		"cats:purr",
		"cats:meow",
	}

	u.DisableScopeAccess("cats:hug", "cats:purr", "cats:delete", "cats:meow")
	assert.EqualValues(t, expectedScopes, u.Scopes)

	u.DisableScopeAccess("cats:hug", "cats:purr", "cats:delete", "cats:meow")
	assert.EqualValues(t, expectedScopes, u.Scopes)
}

func TestClient_EnableTenantAccess_None(t *testing.T) {
	u := expectedClient()

	expectedTenants := []string{
		"29c78d37-a555-4d90-a038-bdb67a82b461",
		"5253ee1a-aaac-49b1-ab7c-85b6d0571366",
	}

	u.EnableTenantAccess("29c78d37-a555-4d90-a038-bdb67a82b461")
	assert.EqualValues(t, expectedTenants, u.AllowedTenantAccess)

	u.EnableTenantAccess("5253ee1a-aaac-49b1-ab7c-85b6d0571366")
	assert.EqualValues(t, expectedTenants, u.AllowedTenantAccess)
}

func TestClient_EnableTenantAccess_One(t *testing.T) {
	u := expectedClient()

	expectedTenantIDs := []string{
		"29c78d37-a555-4d90-a038-bdb67a82b461",
		"5253ee1a-aaac-49b1-ab7c-85b6d0571366",
		"bc7f5c05-3698-4855-8244-b0aac80a3ec1",
	}

	u.EnableTenantAccess("bc7f5c05-3698-4855-8244-b0aac80a3ec1")
	assert.EqualValues(t, expectedTenantIDs, u.AllowedTenantAccess)

	u.EnableTenantAccess("bc7f5c05-3698-4855-8244-b0aac80a3ec1")
	assert.EqualValues(t, expectedTenantIDs, u.AllowedTenantAccess)

	u.EnableTenantAccess("5253ee1a-aaac-49b1-ab7c-85b6d0571366")
	assert.EqualValues(t, expectedTenantIDs, u.AllowedTenantAccess)
}

func TestClient_EnableTenantAccess_Many(t *testing.T) {
	u := expectedClient()

	expectedTenantIDs := []string{
		"29c78d37-a555-4d90-a038-bdb67a82b461",
		"5253ee1a-aaac-49b1-ab7c-85b6d0571366",
		"bc7f5c05-3698-4855-8244-b0aac80a3ec1",
		"b1f8c420-81a0-4980-9bb0-432b2860fd05",
		"c3414224-c98b-42f7-a017-ee0549cca762",
	}

	u.EnableTenantAccess(
		"bc7f5c05-3698-4855-8244-b0aac80a3ec1",
		"b1f8c420-81a0-4980-9bb0-432b2860fd05",
		"c3414224-c98b-42f7-a017-ee0549cca762",
	)
	assert.EqualValues(t, expectedTenantIDs, u.AllowedTenantAccess)

	u.EnableTenantAccess(
		"bc7f5c05-3698-4855-8244-b0aac80a3ec1",
		"b1f8c420-81a0-4980-9bb0-432b2860fd05",
		"c3414224-c98b-42f7-a017-ee0549cca762",
	)
	assert.EqualValues(t, expectedTenantIDs, u.AllowedTenantAccess)
}

func TestClient_DisableTenantAccess_None(t *testing.T) {
	u := expectedClient()

	expectedTenantIDs := []string{
		"29c78d37-a555-4d90-a038-bdb67a82b461",
		"5253ee1a-aaac-49b1-ab7c-85b6d0571366",
	}

	u.DisableTenantAccess("bc7f5c05-3698-4855-8244-b0aac80a3ec1")
	assert.EqualValues(t, expectedTenantIDs, u.AllowedTenantAccess)
}

func TestClient_DisableTenantAccess_One(t *testing.T) {
	u := expectedClient()
	expectedTenants := []string{
		"29c78d37-a555-4d90-a038-bdb67a82b461",
	}

	u.DisableTenantAccess("5253ee1a-aaac-49b1-ab7c-85b6d0571366")
	assert.EqualValues(t, expectedTenants, u.AllowedTenantAccess)

	u.DisableTenantAccess("5253ee1a-aaac-49b1-ab7c-85b6d0571366")
	assert.EqualValues(t, expectedTenants, u.AllowedTenantAccess)

	u.DisableTenantAccess("29c78d37-a555-4d90-a038-bdb67a82b461")
	assert.EqualValues(t, expectedTenants[:len(expectedTenants)-1], u.AllowedTenantAccess)

	u.DisableTenantAccess("b1f8c420-81a0-4980-9bb0-432b2860fd05")
	assert.EqualValues(t, expectedTenants[:len(expectedTenants)-1], u.AllowedTenantAccess)

	u.DisableTenantAccess("c3414224-c98b-42f7-a017-ee0549cca762")
	assert.EqualValues(t, expectedTenants[:len(expectedTenants)-1], u.AllowedTenantAccess)
}

func TestClient_DisableTenantAccess_Many(t *testing.T) {
	u := expectedClient()
	expectedTenants := []string{
		"29c78d37-a555-4d90-a038-bdb67a82b461",
		"5253ee1a-aaac-49b1-ab7c-85b6d0571366",
	}

	u.AllowedTenantAccess = []string{
		"29c78d37-a555-4d90-a038-bdb67a82b461",
		"5253ee1a-aaac-49b1-ab7c-85b6d0571366",
		"bc7f5c05-3698-4855-8244-b0aac80a3ec1",
		"b1f8c420-81a0-4980-9bb0-432b2860fd05",
		"c3414224-c98b-42f7-a017-ee0549cca762",
	}

	u.DisableTenantAccess(
		"bc7f5c05-3698-4855-8244-b0aac80a3ec1",
		"b1f8c420-81a0-4980-9bb0-432b2860fd05",
		"c3414224-c98b-42f7-a017-ee0549cca762",
	)
	assert.EqualValues(t, expectedTenants, u.AllowedTenantAccess)

	u.DisableTenantAccess(
		"bc7f5c05-3698-4855-8244-b0aac80a3ec1",
		"b1f8c420-81a0-4980-9bb0-432b2860fd05",
		"c3414224-c98b-42f7-a017-ee0549cca762",
	)
	assert.EqualValues(t, expectedTenants, u.AllowedTenantAccess)
}

func TestClient_Equal(t *testing.T) {
	tests := []struct {
		description string
		x           client.Client
		y           client.Client
		expected    bool
	}{
		{
			description: "empty should be equal",
			x:           client.Client{},
			y:           client.Client{},
			expected:    true,
		},
		{
			description: "non-empty should not be equal",
			x: client.Client{
				ID: "lol",
			},
			y:        client.Client{},
			expected: false,
		},
		{
			description: "ID should be equal",
			x: client.Client{
				ID: "1",
			},
			y: client.Client{
				ID: "1",
			},
			expected: true,
		},
		{
			description: "ID should not be equal",
			x: client.Client{
				ID: "1",
			},
			y: client.Client{
				ID: "2",
			},
			expected: false,
		},
		{
			description: "Secret should be equal",
			x: client.Client{
				Secret: []byte("Foo"),
			},
			y: client.Client{
				Secret: []byte("Foo"),
			},
			expected: true,
		},
		{
			description: "Secret should not be equal",
			x: client.Client{
				Secret: []byte("Foo"),
			},
			y: client.Client{
				Secret: []byte("Bar"),
			},
			expected: false,
		},
		{
			description: "RedirectURIs should be equal",
			x: client.Client{
				RedirectURIs: []string{"https://example.com/callback", "https://another.example.com/callback"},
			},
			y: client.Client{
				RedirectURIs: []string{"https://example.com/callback", "https://another.example.com/callback"},
			},
			expected: true,
		},
		{
			description: "RedirectURIs should not be equal",
			x: client.Client{
				RedirectURIs: []string{"https://example.com/callback", "https://another.example.com/callback"},
			},
			y: client.Client{
				RedirectURIs: []string{"https://example.com/callback", "https://yet.another.example.com/callback"},
			},
			expected: false,
		},
		{
			description: "GrantTypes should be equal",
			x: client.Client{
				GrantTypes: []string{"client_credentials", "implicit"},
			},
			y: client.Client{
				GrantTypes: []string{"client_credentials", "implicit"},
			},
			expected: true,
		},
		{
			description: "GrantTypes should not be equal",
			x: client.Client{
				GrantTypes: []string{"client_credentials", "implicit"},
			},
			y: client.Client{
				GrantTypes: []string{"client_credentials", "password"},
			},
			expected: false,
		},
		{
			description: "ResponseTypes should be equal",
			x: client.Client{
				ResponseTypes: []string{"code", "token"},
			},
			y: client.Client{
				ResponseTypes: []string{"code", "token"},
			},
			expected: true,
		},
		{
			description: "ResponseTypes should not be equal",
			x: client.Client{
				ResponseTypes: []string{"code", "token"},
			},
			y: client.Client{
				ResponseTypes: []string{"code", "unicorn"},
			},
			expected: false,
		},
		{
			description: "scopes should be equal",
			x: client.Client{
				Scopes: []string{"x2", "10x", "1x red-dot"},
			},
			y: client.Client{
				Scopes: []string{"x2", "10x", "1x red-dot"},
			},
			expected: true,
		},
		{
			description: "scopes length should not be equal",
			x: client.Client{
				Scopes: []string{"1x red-dot"},
			},
			y: client.Client{
				Scopes: []string{"1x red-dot", "x2", "10x"},
			},
			expected: false,
		},
		{
			description: "scopes should not be equal",
			x: client.Client{
				Scopes: []string{"x2", "10x", "1x red-dot"},
			},
			y: client.Client{
				Scopes: []string{"10x", "1x red-dot", "x2"},
			},
			expected: false,
		},
		{
			description: "Owner should be equal",
			x: client.Client{
				Owner: "Widgets Inc.",
			},
			y: client.Client{
				Owner: "Widgets Inc.",
			},
			expected: true,
		},
		{
			description: "Owner should not be equal",
			x: client.Client{
				Owner: "Widgets Inc.",
			},
			y: client.Client{
				Owner: "Fidgets Inc.",
			},
			expected: false,
		},
		{
			description: "Policy URI should be equal",
			x: client.Client{
				PolicyURI: "http://example.com/policy",
			},
			y: client.Client{
				PolicyURI: "http://example.com/policy",
			},
			expected: true,
		},
		{
			description: "Policy URI should not be equal",
			x: client.Client{
				PolicyURI: "http://example.com/policy",
			},
			y: client.Client{
				PolicyURI: "http://example.com/tos",
			},
			expected: false,
		},
		{
			description: "TermsOfServiceURI should be equal",
			x: client.Client{
				TermsOfServiceURI: "https://cats.example.com/tos",
			},
			y: client.Client{
				TermsOfServiceURI: "https://cats.example.com/tos",
			},
			expected: true,
		},
		{
			description: "TermsOfServiceURI should not be equal",
			x: client.Client{
				TermsOfServiceURI: "https://cats.example.com/tos",
			},
			y: client.Client{
				TermsOfServiceURI: "https://cats.example.com/meowmix",
			},
			expected: false,
		},
		{
			description: "ClientURI should be equal",
			x: client.Client{
				ClientURI: "https://myapp.example.com/about",
			},
			y: client.Client{
				ClientURI: "https://myapp.example.com/about",
			},
			expected: true,
		},
		{
			description: "ClientURI should not be equal",
			x: client.Client{
				ClientURI: "https://myapp.example.com/about",
			},
			y: client.Client{
				ClientURI: "https://myapp.example.com/mycats",
			},
			expected: false,
		},
		{
			description: "LogoURI should be equal",
			x: client.Client{
				LogoURI: "https://myapp.example.com/logo256x256.png",
			},
			y: client.Client{
				LogoURI: "https://myapp.example.com/logo256x256.png",
			},
			expected: true,
		},
		{
			description: "LogoURI should not be equal",
			x: client.Client{
				LogoURI: "https://myapp.example.com/logo256x256.png",
			},
			y: client.Client{
				LogoURI: "https://myapp.example.com/logrus.png",
			},
			expected: false,
		},
		{
			description: "Contacts should be equal",
			x: client.Client{
				Contacts: []string{"foo@example.com", "bar@example.com"},
			},
			y: client.Client{
				Contacts: []string{"foo@example.com", "bar@example.com"},
			},
			expected: true,
		},
		{
			description: "Contacts should not be equal",
			x: client.Client{
				Contacts: []string{"foo@example.com", "bar@example.com"},
			},
			y: client.Client{
				Contacts: []string{"bar@example.com", "foo@example.com"},
			},
			expected: false,
		},
		{
			description: "client should be equal",
			x: client.Client{
				ID:                  "foo",
				Name:                "Foo bar App",
				AllowedTenantAccess: []string{"78288f2c-4fd5-4f52-9e28-9d17e5524e83", "39d3f55e-3fa3-4d65-b2b2-18ef2904115b"},
				Secret:              []byte("S@ltyP@ssw0rd"),
				RedirectURIs:        []string{"https://app.foo.example.com/callback", "https://dev-app.foo.example.com/callback"},
				GrantTypes:          []string{"client_credentials", "implicit"},
				ResponseTypes:       []string{"code", "token"},
				Scopes:              []string{"urn.foo.bar", "urn.foo.baz"},
				Owner:               "FooBar Baz inc.",
				PolicyURI:           "https://foo.example.com/policy",
				TermsOfServiceURI:   "https://foo.example.com/tos",
				ClientURI:           "https://app.foo.example.com/about",
				LogoURI:             "https://logos.example.com/happy-kitten.jpg",
				Contacts:            []string{"foo@example.com", "bar@example.com"},
				Public:              true,
				Disabled:            false,
			},
			y: client.Client{
				ID:                  "foo",
				Name:                "Foo bar App",
				AllowedTenantAccess: []string{"78288f2c-4fd5-4f52-9e28-9d17e5524e83", "39d3f55e-3fa3-4d65-b2b2-18ef2904115b"},
				Secret:              []byte("S@ltyP@ssw0rd"),
				RedirectURIs:        []string{"https://app.foo.example.com/callback", "https://dev-app.foo.example.com/callback"},
				GrantTypes:          []string{"client_credentials", "implicit"},
				ResponseTypes:       []string{"code", "token"},
				Scopes:              []string{"urn.foo.bar", "urn.foo.baz"},
				Owner:               "FooBar Baz inc.",
				PolicyURI:           "https://foo.example.com/policy",
				TermsOfServiceURI:   "https://foo.example.com/tos",
				ClientURI:           "https://app.foo.example.com/about",
				LogoURI:             "https://logos.example.com/happy-kitten.jpg",
				Contacts:            []string{"foo@example.com", "bar@example.com"},
				Public:              true,
				Disabled:            false,
			},
			expected: true,
		},
		{
			description: "client should not be equal",
			x: client.Client{
				ID:                  "foo",
				Name:                "Foo bar App",
				AllowedTenantAccess: []string{"78288f2c-4fd5-4f52-9e28-9d17e5524e83", "39d3f55e-3fa3-4d65-b2b2-18ef2904115b"},
				Secret:              []byte("S@ltyP@ssw0rd"),
				RedirectURIs:        []string{"https://app.foo.example.com/callback", "https://dev-app.foo.example.com/callback"},
				GrantTypes:          []string{"client_credentials", "implicit"},
				ResponseTypes:       []string{"code", "token"},
				Scopes:              []string{"urn.foo.bar", "urn.foo.baz"},
				Owner:               "FooBar Baz inc.",
				PolicyURI:           "https://foo.example.com/policy",
				TermsOfServiceURI:   "https://foo.example.com/tos",
				ClientURI:           "https://app.foo.example.com/about",
				LogoURI:             "https://logos.example.com/happy-kitten.jpg",
				Contacts:            []string{"foo@example.com", "bar@example.com"},
				Public:              true,
				Disabled:            false,
			},
			y: client.Client{
				ID:                  "foo",
				Name:                "Foo bar App",
				AllowedTenantAccess: []string{"243763ae-c1ba-4988-863d-39d73884f17a", "78288f2c-4fd5-4f52-9e28-9d17e5524e83"},
				Secret:              []byte("hunter2"),
				RedirectURIs:        []string{"https://app.foo.example.com/callback", "https://dev-app.foo.example.com/callback"},
				GrantTypes:          []string{"client_credentials", "implicit"},
				ResponseTypes:       []string{"code", "token"},
				Scopes:              []string{"urn.foo.bar", "urn.foo.baz"},
				Owner:               "SomeCompany inc.",
				PolicyURI:           "https://foo.example.com/policy",
				TermsOfServiceURI:   "https://foo.example.com/tos",
				ClientURI:           "https://app.foo.example.com/about",
				LogoURI:             "https://logos.example.com/happy-kitten.jpg",
				Contacts:            []string{"foo.bar@example.com", "foo.baz@example.com"},
				Public:              true,
				Disabled:            false,
			},
			expected: false,
		},
	}

	for _, testcase := range tests {
		assert.Equal(t, testcase.expected, testcase.x.Equal(testcase.y), testcase.description)
	}
}

func TestClient_IsEmpty(t *testing.T) {
	notEmptyClient := client.Client{
		ID: "lol-not-empty",
	}
	assert.Equal(t, notEmptyClient.IsEmpty(), false)

	emptyClient := client.Client{}
	assert.Equal(t, emptyClient.IsEmpty(), true)
}
