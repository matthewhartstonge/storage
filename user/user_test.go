package user_test

import (
	"testing"

	"github.com/MatthewHartstonge/storage/user"
	"github.com/stretchr/testify/assert"
)

func expectedUser() user.User {
	return user.User{
		ID: "cc935033-d1b0-4bd8-b209-e6fbffe6b624",
		TenantIDs: []string{
			"29c78d37-a555-4d90-a038-bdb67a82b461",
			"5253ee1a-aaac-49b1-ab7c-85b6d0571366",
		},
		Username: "kitteh@example.com",
		Password: "i<3kittehs",
		Scopes: []string{
			"cats:read",
			"cats:delete",
		},
		FirstName:  "Fluffy",
		LastName:   "McKittison",
		ProfileURI: "",
	}
}

func TestUser_AddScopes_None(t *testing.T) {
	u := expectedUser()

	expectedScopes := []string{
		"cats:read",
		"cats:delete",
	}

	u.AddScopes("cats:read")
	assert.EqualValues(t, expectedScopes, u.Scopes)

	u.AddScopes("cats:delete")
	assert.EqualValues(t, expectedScopes, u.Scopes)
}

func TestUser_AddScopes_One(t *testing.T) {
	u := expectedUser()

	expectedScopes := []string{
		"cats:read",
		"cats:delete",
		"cats:hug",
	}

	u.AddScopes("cats:hug")
	assert.EqualValues(t, expectedScopes, u.Scopes)

	u.AddScopes("cats:hug")
	assert.EqualValues(t, expectedScopes, u.Scopes)

	u.AddScopes("cats:read")
	assert.EqualValues(t, expectedScopes, u.Scopes)
}

func TestUser_AddScopes_Many(t *testing.T) {
	u := expectedUser()

	expectedScopes := []string{
		"cats:read",
		"cats:delete",
		"cats:hug",
		"cats:purr",
		"cats:meow",
	}

	u.AddScopes("cats:hug", "cats:purr", "cats:meow")
	assert.EqualValues(t, expectedScopes, u.Scopes)

	u.AddScopes("cats:hug", "cats:purr", "cats:meow")
	assert.EqualValues(t, expectedScopes, u.Scopes)
}

func TestUser_RemoveScopes_None(t *testing.T) {
	u := expectedUser()

	expectedScopes := []string{
		"cats:read",
		"cats:delete",
	}

	u.RemoveScopes("cats:hug")
	assert.EqualValues(t, expectedScopes, u.Scopes)
}

func TestUser_RemoveScopes_One(t *testing.T) {
	u := expectedUser()
	expectedScopes := []string{
		"cats:delete",
	}

	u.RemoveScopes("cats:read")
	assert.EqualValues(t, expectedScopes, u.Scopes)

	u.RemoveScopes("cats:read")
	assert.EqualValues(t, expectedScopes, u.Scopes)

	u.RemoveScopes("cats:delete")
	assert.EqualValues(t, expectedScopes[:len(expectedScopes)-1], u.Scopes)

	u.RemoveScopes("cats:read")
	assert.EqualValues(t, expectedScopes[:len(expectedScopes)-1], u.Scopes)

	u.RemoveScopes("cats:mug")
	assert.EqualValues(t, expectedScopes[:len(expectedScopes)-1], u.Scopes)
}

func TestUser_RemoveScopes_Many(t *testing.T) {
	u := expectedUser()
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

	u.RemoveScopes("cats:hug", "cats:purr", "cats:delete", "cats:meow")
	assert.EqualValues(t, expectedScopes, u.Scopes)

	u.RemoveScopes("cats:hug", "cats:purr", "cats:delete", "cats:meow")
	assert.EqualValues(t, expectedScopes, u.Scopes)
}

func TestUser_AddTenantIDs_None(t *testing.T) {
	u := expectedUser()

	expectedTenants := []string{
		"29c78d37-a555-4d90-a038-bdb67a82b461",
		"5253ee1a-aaac-49b1-ab7c-85b6d0571366",
	}

	u.AddTenantIDs("29c78d37-a555-4d90-a038-bdb67a82b461")
	assert.EqualValues(t, expectedTenants, u.TenantIDs)

	u.AddTenantIDs("5253ee1a-aaac-49b1-ab7c-85b6d0571366")
	assert.EqualValues(t, expectedTenants, u.TenantIDs)
}

func TestUser_AddTenantIDs_One(t *testing.T) {
	u := expectedUser()

	expectedTenantIDs := []string{
		"29c78d37-a555-4d90-a038-bdb67a82b461",
		"5253ee1a-aaac-49b1-ab7c-85b6d0571366",
		"bc7f5c05-3698-4855-8244-b0aac80a3ec1",
	}

	u.AddTenantIDs("bc7f5c05-3698-4855-8244-b0aac80a3ec1")
	assert.EqualValues(t, expectedTenantIDs, u.TenantIDs)

	u.AddTenantIDs("bc7f5c05-3698-4855-8244-b0aac80a3ec1")
	assert.EqualValues(t, expectedTenantIDs, u.TenantIDs)

	u.AddTenantIDs("5253ee1a-aaac-49b1-ab7c-85b6d0571366")
	assert.EqualValues(t, expectedTenantIDs, u.TenantIDs)
}

func TestUser_AddTenantIDs_Many(t *testing.T) {
	u := expectedUser()

	expectedTenantIDs := []string{
		"29c78d37-a555-4d90-a038-bdb67a82b461",
		"5253ee1a-aaac-49b1-ab7c-85b6d0571366",
		"bc7f5c05-3698-4855-8244-b0aac80a3ec1",
		"b1f8c420-81a0-4980-9bb0-432b2860fd05",
		"c3414224-c98b-42f7-a017-ee0549cca762",
	}

	u.AddTenantIDs(
		"bc7f5c05-3698-4855-8244-b0aac80a3ec1",
		"b1f8c420-81a0-4980-9bb0-432b2860fd05",
		"c3414224-c98b-42f7-a017-ee0549cca762",
	)
	assert.EqualValues(t, expectedTenantIDs, u.TenantIDs)

	u.AddTenantIDs(
		"bc7f5c05-3698-4855-8244-b0aac80a3ec1",
		"b1f8c420-81a0-4980-9bb0-432b2860fd05",
		"c3414224-c98b-42f7-a017-ee0549cca762",
	)
	assert.EqualValues(t, expectedTenantIDs, u.TenantIDs)
}

func TestUser_RemoveTenantIDs_None(t *testing.T) {
	u := expectedUser()

	expectedTenantIDs := []string{
		"29c78d37-a555-4d90-a038-bdb67a82b461",
		"5253ee1a-aaac-49b1-ab7c-85b6d0571366",
	}

	u.RemoveTenantIDs("bc7f5c05-3698-4855-8244-b0aac80a3ec1")
	assert.EqualValues(t, expectedTenantIDs, u.TenantIDs)
}

func TestUser_RemoveTenantIDs_One(t *testing.T) {
	u := expectedUser()
	expectedTenants := []string{
		"29c78d37-a555-4d90-a038-bdb67a82b461",
	}

	u.RemoveTenantIDs("5253ee1a-aaac-49b1-ab7c-85b6d0571366")
	assert.EqualValues(t, expectedTenants, u.TenantIDs)

	u.RemoveTenantIDs("5253ee1a-aaac-49b1-ab7c-85b6d0571366")
	assert.EqualValues(t, expectedTenants, u.TenantIDs)

	u.RemoveTenantIDs("29c78d37-a555-4d90-a038-bdb67a82b461")
	assert.EqualValues(t, expectedTenants[:len(expectedTenants)-1], u.TenantIDs)

	u.RemoveTenantIDs("b1f8c420-81a0-4980-9bb0-432b2860fd05")
	assert.EqualValues(t, expectedTenants[:len(expectedTenants)-1], u.TenantIDs)

	u.RemoveTenantIDs("c3414224-c98b-42f7-a017-ee0549cca762")
	assert.EqualValues(t, expectedTenants[:len(expectedTenants)-1], u.TenantIDs)
}

func TestUser_RemoveTenantIDs_Many(t *testing.T) {
	u := expectedUser()
	expectedTenants := []string{
		"29c78d37-a555-4d90-a038-bdb67a82b461",
		"5253ee1a-aaac-49b1-ab7c-85b6d0571366",
	}

	u.TenantIDs = []string{
		"29c78d37-a555-4d90-a038-bdb67a82b461",
		"5253ee1a-aaac-49b1-ab7c-85b6d0571366",
		"bc7f5c05-3698-4855-8244-b0aac80a3ec1",
		"b1f8c420-81a0-4980-9bb0-432b2860fd05",
		"c3414224-c98b-42f7-a017-ee0549cca762",
	}

	u.RemoveTenantIDs(
		"bc7f5c05-3698-4855-8244-b0aac80a3ec1",
		"b1f8c420-81a0-4980-9bb0-432b2860fd05",
		"c3414224-c98b-42f7-a017-ee0549cca762",
	)
	assert.EqualValues(t, expectedTenants, u.TenantIDs)

	u.RemoveTenantIDs(
		"bc7f5c05-3698-4855-8244-b0aac80a3ec1",
		"b1f8c420-81a0-4980-9bb0-432b2860fd05",
		"c3414224-c98b-42f7-a017-ee0549cca762",
	)
	assert.EqualValues(t, expectedTenants, u.TenantIDs)
}

func TestUser_Equal(t *testing.T) {
	tests := []struct {
		description string
		x           user.User
		y           user.User
		expected    bool
	}{
		{
			description: "empty should be equal",
			x:           user.User{},
			y:           user.User{},
			expected:    true,
		},
		{
			description: "non-empty should not be equal",
			x: user.User{
				ID: "lol",
			},
			y:        user.User{},
			expected: false,
		},
		{
			description: "ID should be equal",
			x: user.User{
				ID: "1",
			},
			y: user.User{
				ID: "1",
			},
			expected: true,
		},
		{
			description: "ID should not be equal",
			x: user.User{
				ID: "1",
			},
			y: user.User{
				ID: "2",
			},
			expected: false,
		},
		{
			description: "username should be equal",
			x: user.User{
				Username: "timmy",
			},
			y: user.User{
				Username: "timmy",
			},
			expected: true,
		},
		{
			description: "username should not be equal",
			x: user.User{
				Username: "timmy",
			},
			y: user.User{
				Username: "jimmy",
			},
			expected: false,
		},
		{
			description: "password should be equal",
			x: user.User{
				Password: "salty",
			},
			y: user.User{
				Password: "salty",
			},
			expected: true,
		},
		{
			description: "password should not be equal",
			x: user.User{
				Username: "salty",
			},
			y: user.User{
				Username: "not-very-salty",
			},
			expected: false,
		},
		{
			description: "scopes should be equal",
			x: user.User{
				Scopes: []string{"x2", "10x", "1x red-dot"},
			},
			y: user.User{
				Scopes: []string{"x2", "10x", "1x red-dot"},
			},
			expected: true,
		},
		{
			description: "scopes length should not be equal",
			x: user.User{
				Scopes: []string{"1x red-dot"},
			},
			y: user.User{
				Scopes: []string{"1x red-dot", "x2", "10x"},
			},
			expected: false,
		},
		{
			description: "scopes should not be equal",
			x: user.User{
				Scopes: []string{"x2", "10x", "1x red-dot"},
			},
			y: user.User{
				Scopes: []string{"10x", "1x red-dot", "x2"},
			},
			expected: false,
		},
		{
			description: "firstname should be equal",
			x: user.User{
				FirstName: "bob lee",
			},
			y: user.User{
				FirstName: "bob lee",
			},
			expected: true,
		},
		{
			description: "firstname should not be equal",
			x: user.User{
				LastName: "bob lee",
			},
			y: user.User{
				LastName: "bobby lee",
			},
			expected: false,
		},
		{
			description: "lastname should be equal",
			x: user.User{
				FirstName: "swagger",
			},
			y: user.User{
				FirstName: "swagger",
			},
			expected: true,
		},
		{
			description: "lastname should not be equal",
			x: user.User{
				LastName: "swagger",
			},
			y: user.User{
				LastName: "swaggerz",
			},
			expected: false,
		},
		{
			description: "profile uri should be equal",
			x: user.User{
				ProfileURI: "https://cats.example.com/cat1.jpg",
			},
			y: user.User{
				ProfileURI: "https://cats.example.com/cat1.jpg",
			},
			expected: true,
		},
		{
			description: "profile uri should not be equal",
			x: user.User{
				ProfileURI: "https://cats.example.com/cat1.jpg",
			},
			y: user.User{
				ProfileURI: "https://dogs.example.com/dog1.jpg",
			},
			expected: false,
		},
		{
			description: "disabled should be equal",
			x: user.User{
				Disabled: false,
			},
			y: user.User{
				Disabled: false,
			},
			expected: true,
		},
		{
			description: "disabled should not be equal",
			x: user.User{
				Disabled: false,
			},
			y: user.User{
				Disabled: true,
			},
			expected: false,
		},
		{
			description: "user should be equal",
			x: user.User{
				ID:         "1",
				TenantIDs:  []string{"apple", "lettuce"},
				Username:   "boblee@auth.example.com",
				Password:   "saltypa@ssw0rd",
				Scopes:     []string{"10x", "2x"},
				FirstName:  "Bob Lee",
				LastName:   "Swagger",
				ProfileURI: "https://marines.example.com/boblee.png",
				Disabled:   false,
			},
			y: user.User{
				ID:         "1",
				TenantIDs:  []string{"apple", "lettuce"},
				Username:   "boblee@auth.example.com",
				Password:   "saltypa@ssw0rd",
				Scopes:     []string{"10x", "2x"},
				FirstName:  "Bob Lee",
				LastName:   "Swagger",
				ProfileURI: "https://marines.example.com/boblee.png",
				Disabled:   false,
			},
			expected: true,
		},
		{
			description: "user should not be equal",
			x: user.User{
				ID:         "1",
				TenantIDs:  []string{"apple", "lettuce"},
				Username:   "boblee@auth.example.com",
				Password:   "saltypa@ssw0rd",
				Scopes:     []string{"10x", "2x"},
				FirstName:  "Bob Lee",
				LastName:   "Swagger",
				ProfileURI: "https://marines.example.com/boblee.png",
				Disabled:   false,
			},
			y: user.User{
				ID:         "1",
				TenantIDs:  []string{"apple", "lettuce"},
				Username:   "boblee@auth.example.com",
				Password:   "saltypa@ssw0rd",
				Scopes:     []string{"10x"},
				FirstName:  "Bob Lee",
				LastName:   "Swagger",
				ProfileURI: "https://marines.example.com/boblee.png",
				Disabled:   false,
			},
			expected: false,
		},
	}

	for _, testcase := range tests {
		assert.Equal(t, testcase.expected, testcase.x.Equal(testcase.y), testcase.description)
	}
}

func TestUser_IsEmpty(t *testing.T) {
	notEmptyUser := user.User{
		ID: "lol-not-empty",
	}
	assert.Equal(t, notEmptyUser.IsEmpty(), false)

	emptyUser := user.User{}
	assert.Equal(t, emptyUser.IsEmpty(), true)
}
