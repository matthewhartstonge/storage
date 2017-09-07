package user_test

import (
	"testing"

	"github.com/MatthewHartstonge/storage/user"
	"github.com/stretchr/testify/assert"
)

func expectedUser(t *testing.T) user.User {
	t.Helper()
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
	u := expectedUser(t)

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
	u := expectedUser(t)

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
	u := expectedUser(t)

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
	u := expectedUser(t)

	expectedScopes := []string{
		"cats:read",
		"cats:delete",
	}

	u.RemoveScopes("cats:hug")
	assert.EqualValues(t, expectedScopes, u.Scopes)
}

func TestUser_RemoveScopes_One(t *testing.T) {
	u := expectedUser(t)
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
	u := expectedUser(t)
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
	u := expectedUser(t)

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
	u := expectedUser(t)

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
	u := expectedUser(t)

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
	u := expectedUser(t)

	expectedTenantIDs := []string{
		"29c78d37-a555-4d90-a038-bdb67a82b461",
		"5253ee1a-aaac-49b1-ab7c-85b6d0571366",
	}

	u.RemoveTenantIDs("bc7f5c05-3698-4855-8244-b0aac80a3ec1")
	assert.EqualValues(t, expectedTenantIDs, u.TenantIDs)
}

func TestUser_RemoveTenantIDs_One(t *testing.T) {
	u := expectedUser(t)
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
	u := expectedUser(t)
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
