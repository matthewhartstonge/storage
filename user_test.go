package storage

import (
	// Standard Library Imports
	"testing"

	// External Imports
	"github.com/stretchr/testify/assert"
)

func expectedUser() User {
	return User{
		ID:         "cc935033-d1b0-4bd8-b209-e6fbffe6b624",
		CreateTime: 123,
		UpdateTime: 987,
		AllowedTenantAccess: []string{
			"29c78d37-a555-4d90-a038-bdb67a82b461",
			"5253ee1a-aaac-49b1-ab7c-85b6d0571366",
		},
		AllowedPersonAccess: []string{
			"7f6dfb7d-a6b0-442e-aab0-ad54c917f506",
			"794a55bd-69d4-4668-b319-62bfa0cd59ac",
		},
		Scopes: []string{
			"cats:read",
			"cats:delete",
		},
		Roles: []string{
			"user",
			"printer",
		},
		PersonID:   "123",
		Disabled:   false,
		Username:   "kitteh@example.com",
		Password:   "i<3kittehs",
		FirstName:  "Fluffy",
		LastName:   "McKittison",
		ProfileURI: "https://kittehs-unite.meow",
	}
}

func TestUser_FullName(t *testing.T) {
	type fields struct {
		FirstName string
		LastName  string
	}
	tests := []struct {
		name         string
		fields       fields
		wantFullName string
	}{
		{
			name:         "Should return an empty full name",
			fields:       fields{},
			wantFullName: "",
		},
		{
			name: "Should return a full name, with only a first name",
			fields: fields{
				FirstName: "Jimmy",
			},
			wantFullName: "Jimmy",
		},
		{
			name: "Should return a full name, with only a last name",
			fields: fields{
				LastName: "Barnes",
			},
			wantFullName: "Barnes",
		},
		{
			name: "Should return a full name",
			fields: fields{
				FirstName: "Jimmy",
				LastName:  "Barnes",
			},
			wantFullName: "Jimmy Barnes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := User{
				FirstName: tt.fields.FirstName,
				LastName:  tt.fields.LastName,
			}
			if gotFullName := u.FullName(); gotFullName != tt.wantFullName {
				t.Errorf(
					"FullName()\ngot:  %#+v\nwant: %#+v",
					gotFullName,
					tt.wantFullName,
				)
			}
		})
	}
}

func TestUser_EnableScopeAccess_None(t *testing.T) {
	u := expectedUser()

	expectedScopes := []string{
		"cats:read",
		"cats:delete",
	}

	u.EnableScopeAccess("cats:read")
	assert.EqualValues(t, expectedScopes, u.Scopes)

	u.EnableScopeAccess("cats:delete")
	assert.EqualValues(t, expectedScopes, u.Scopes)
}

func TestUser_EnableScopeAccess_One(t *testing.T) {
	u := expectedUser()

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

func TestUser_EnableScopeAccess_Many(t *testing.T) {
	u := expectedUser()

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

func TestUser_DisableScopeAccess_None(t *testing.T) {
	u := expectedUser()

	expectedScopes := []string{
		"cats:read",
		"cats:delete",
	}

	u.DisableScopeAccess("cats:hug")
	assert.EqualValues(t, expectedScopes, u.Scopes)
}

func TestUser_DisableScopeAccess_One(t *testing.T) {
	u := expectedUser()
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

func TestUser_DisableScopeAccess_Many(t *testing.T) {
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

	u.DisableScopeAccess("cats:hug", "cats:purr", "cats:delete", "cats:meow")
	assert.EqualValues(t, expectedScopes, u.Scopes)

	u.DisableScopeAccess("cats:hug", "cats:purr", "cats:delete", "cats:meow")
	assert.EqualValues(t, expectedScopes, u.Scopes)
}

func TestUser_EnableTenantAccess_None(t *testing.T) {
	u := expectedUser()

	expectedTenants := []string{
		"29c78d37-a555-4d90-a038-bdb67a82b461",
		"5253ee1a-aaac-49b1-ab7c-85b6d0571366",
	}

	u.EnableTenantAccess("29c78d37-a555-4d90-a038-bdb67a82b461")
	assert.EqualValues(t, expectedTenants, u.AllowedTenantAccess)

	u.EnableTenantAccess("5253ee1a-aaac-49b1-ab7c-85b6d0571366")
	assert.EqualValues(t, expectedTenants, u.AllowedTenantAccess)
}

func TestUser_EnableTenantAccess_One(t *testing.T) {
	u := expectedUser()

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

func TestUser_EnableTenantAccess_Many(t *testing.T) {
	u := expectedUser()

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

func TestUser_DisableTenantAccess_None(t *testing.T) {
	u := expectedUser()

	expectedTenantIDs := []string{
		"29c78d37-a555-4d90-a038-bdb67a82b461",
		"5253ee1a-aaac-49b1-ab7c-85b6d0571366",
	}

	u.DisableTenantAccess("bc7f5c05-3698-4855-8244-b0aac80a3ec1")
	assert.EqualValues(t, expectedTenantIDs, u.AllowedTenantAccess)
}

func TestUser_DisableTenantAccess_One(t *testing.T) {
	u := expectedUser()
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

func TestUser_DisableTenantAccess_Many(t *testing.T) {
	u := expectedUser()
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

func TestUser_EnablePeopleAccess_None(t *testing.T) {
	u := expectedUser()

	expectedPeopleIDs := []string{
		"7f6dfb7d-a6b0-442e-aab0-ad54c917f506",
		"794a55bd-69d4-4668-b319-62bfa0cd59ac",
	}

	u.EnablePeopleAccess("7f6dfb7d-a6b0-442e-aab0-ad54c917f506")
	assert.EqualValues(t, expectedPeopleIDs, u.AllowedPersonAccess)

	u.EnablePeopleAccess("794a55bd-69d4-4668-b319-62bfa0cd59ac")
	assert.EqualValues(t, expectedPeopleIDs, u.AllowedPersonAccess)
}

func TestUser_EnablePeopleAccess_One(t *testing.T) {
	u := expectedUser()

	expectedPeopleIDs := []string{
		"7f6dfb7d-a6b0-442e-aab0-ad54c917f506",
		"794a55bd-69d4-4668-b319-62bfa0cd59ac",
		"bc7f5c05-3698-4855-8244-b0aac80a3ec1",
	}

	u.EnablePeopleAccess("bc7f5c05-3698-4855-8244-b0aac80a3ec1")
	assert.EqualValues(t, expectedPeopleIDs, u.AllowedPersonAccess)

	u.EnablePeopleAccess("bc7f5c05-3698-4855-8244-b0aac80a3ec1")
	assert.EqualValues(t, expectedPeopleIDs, u.AllowedPersonAccess)

	u.EnablePeopleAccess("794a55bd-69d4-4668-b319-62bfa0cd59ac")
	assert.EqualValues(t, expectedPeopleIDs, u.AllowedPersonAccess)
}

func TestUser_EnablePeopleAccess_Many(t *testing.T) {
	u := expectedUser()

	expectedPeopleIDs := []string{
		"7f6dfb7d-a6b0-442e-aab0-ad54c917f506",
		"794a55bd-69d4-4668-b319-62bfa0cd59ac",
		"bc7f5c05-3698-4855-8244-b0aac80a3ec1",
		"b1f8c420-81a0-4980-9bb0-432b2860fd05",
		"c3414224-c98b-42f7-a017-ee0549cca762",
	}

	u.EnablePeopleAccess(
		"bc7f5c05-3698-4855-8244-b0aac80a3ec1",
		"b1f8c420-81a0-4980-9bb0-432b2860fd05",
		"c3414224-c98b-42f7-a017-ee0549cca762",
	)
	assert.EqualValues(t, expectedPeopleIDs, u.AllowedPersonAccess)

	u.EnablePeopleAccess(
		"bc7f5c05-3698-4855-8244-b0aac80a3ec1",
		"b1f8c420-81a0-4980-9bb0-432b2860fd05",
		"c3414224-c98b-42f7-a017-ee0549cca762",
	)
	assert.EqualValues(t, expectedPeopleIDs, u.AllowedPersonAccess)
}

func TestUser_DisablePeopleAccess_None(t *testing.T) {
	u := expectedUser()

	expectedPeopleIDs := []string{
		"7f6dfb7d-a6b0-442e-aab0-ad54c917f506",
		"794a55bd-69d4-4668-b319-62bfa0cd59ac",
	}

	u.DisablePeopleAccess("bc7f5c05-3698-4855-8244-b0aac80a3ec1")
	assert.EqualValues(t, expectedPeopleIDs, u.AllowedPersonAccess)
}

func TestUser_DisablePeopleAccess_One(t *testing.T) {
	u := expectedUser()

	expectedPeopleIDs := []string{
		"7f6dfb7d-a6b0-442e-aab0-ad54c917f506",
	}

	u.DisablePeopleAccess("794a55bd-69d4-4668-b319-62bfa0cd59ac")
	assert.EqualValues(t, expectedPeopleIDs, u.AllowedPersonAccess)

	u.DisablePeopleAccess("794a55bd-69d4-4668-b319-62bfa0cd59ac")
	assert.EqualValues(t, expectedPeopleIDs, u.AllowedPersonAccess)

	u.DisablePeopleAccess("7f6dfb7d-a6b0-442e-aab0-ad54c917f506")
	assert.EqualValues(t, expectedPeopleIDs[:len(expectedPeopleIDs)-1], u.AllowedPersonAccess)

	u.DisablePeopleAccess("b1f8c420-81a0-4980-9bb0-432b2860fd05")
	assert.EqualValues(t, expectedPeopleIDs[:len(expectedPeopleIDs)-1], u.AllowedPersonAccess)

	u.DisablePeopleAccess("c3414224-c98b-42f7-a017-ee0549cca762")
	assert.EqualValues(t, expectedPeopleIDs[:len(expectedPeopleIDs)-1], u.AllowedPersonAccess)
}

func TestUser_DisablePeopleAccess_Many(t *testing.T) {
	u := expectedUser()

	expectedPeopleIDs := []string{
		"7f6dfb7d-a6b0-442e-aab0-ad54c917f506",
		"794a55bd-69d4-4668-b319-62bfa0cd59ac",
	}

	u.AllowedPersonAccess = []string{
		"7f6dfb7d-a6b0-442e-aab0-ad54c917f506",
		"794a55bd-69d4-4668-b319-62bfa0cd59ac",
		"bc7f5c05-3698-4855-8244-b0aac80a3ec1",
		"b1f8c420-81a0-4980-9bb0-432b2860fd05",
		"c3414224-c98b-42f7-a017-ee0549cca762",
	}

	u.DisablePeopleAccess(
		"bc7f5c05-3698-4855-8244-b0aac80a3ec1",
		"b1f8c420-81a0-4980-9bb0-432b2860fd05",
		"c3414224-c98b-42f7-a017-ee0549cca762",
	)
	assert.EqualValues(t, expectedPeopleIDs, u.AllowedPersonAccess)

	u.DisablePeopleAccess(
		"bc7f5c05-3698-4855-8244-b0aac80a3ec1",
		"b1f8c420-81a0-4980-9bb0-432b2860fd05",
		"c3414224-c98b-42f7-a017-ee0549cca762",
	)
	assert.EqualValues(t, expectedPeopleIDs, u.AllowedPersonAccess)
}

func TestUser_EnableRoles_None(t *testing.T) {
	u := expectedUser()

	expectedUserRoles := []string{
		"user",
		"printer",
	}

	u.EnableRoles("user")
	assert.EqualValues(t, expectedUserRoles, u.Roles)

	u.EnableRoles("printer")
	assert.EqualValues(t, expectedUserRoles, u.Roles)
}

func TestUser_EnableRoles_One(t *testing.T) {
	u := expectedUser()

	expectedRoles := []string{
		"user",
		"printer",
		"administrator",
	}

	u.EnableRoles("administrator")
	assert.EqualValues(t, expectedRoles, u.Roles)

	u.EnableRoles("administrator")
	assert.EqualValues(t, expectedRoles, u.Roles)

	u.EnableRoles("printer")
	assert.EqualValues(t, expectedRoles, u.Roles)
}

func TestUser_EnableRoles_Many(t *testing.T) {
	u := expectedUser()

	expectedRoles := []string{
		"user",
		"printer",
		"administrator",
		"finance",
		"groups",
	}

	u.EnableRoles("administrator", "finance", "groups")
	assert.EqualValues(t, expectedRoles, u.Roles)

	u.EnableRoles("administrator", "finance", "groups")
	assert.EqualValues(t, expectedRoles, u.Roles)
}

func TestUser_DisableRoles_None(t *testing.T) {
	u := expectedUser()

	expectedRoles := []string{
		"user",
		"printer",
	}

	u.DisableRoles("administrator")
	assert.EqualValues(t, expectedRoles, u.Roles)
}

func TestUser_DisableRoles_One(t *testing.T) {
	u := expectedUser()

	expectedRoles := []string{
		"user",
	}

	u.DisableRoles("printer")
	assert.EqualValues(t, expectedRoles, u.Roles)

	u.DisableRoles("printer")
	assert.EqualValues(t, expectedRoles, u.Roles)

	u.DisableRoles("user")
	assert.EqualValues(t, expectedRoles[:len(expectedRoles)-1], u.Roles)

	u.DisableRoles("administrator")
	assert.EqualValues(t, expectedRoles[:len(expectedRoles)-1], u.Roles)

	u.DisableRoles("finance")
	assert.EqualValues(t, expectedRoles[:len(expectedRoles)-1], u.Roles)
}

func TestUser_DisableRoles_Many(t *testing.T) {
	u := expectedUser()

	expectedPeopleIDs := []string{
		"user",
		"printer",
	}

	u.Roles = []string{
		"user",
		"printer",
		"administrator",
		"finance",
		"groups",
	}

	u.DisableRoles(
		"administrator",
		"finance",
		"groups",
	)
	assert.EqualValues(t, expectedPeopleIDs, u.Roles)

	u.DisableRoles(
		"administrator",
		"finance",
		"groups",
	)
	assert.EqualValues(t, expectedPeopleIDs, u.Roles)
}

func TestUser_Equal(t *testing.T) {
	tests := []struct {
		description string
		x           User
		y           User
		expected    bool
	}{
		{
			description: "empty should be equal",
			x:           User{},
			y:           User{},
			expected:    true,
		},
		{
			description: "non-empty should not be equal",
			x: User{
				ID: "lol",
			},
			y:        User{},
			expected: false,
		},
		{
			description: "ID should be equal",
			x: User{
				ID: "1",
			},
			y: User{
				ID: "1",
			},
			expected: true,
		},
		{
			description: "ID should not be equal",
			x: User{
				ID: "1",
			},
			y: User{
				ID: "2",
			},
			expected: false,
		},
		{
			description: "Create time should be equal",
			x: User{
				CreateTime: 123456789,
			},
			y: User{
				CreateTime: 123456789,
			},
			expected: true,
		},
		{
			description: "Create time should not be equal",
			x: User{
				CreateTime: 123456789,
			},
			y: User{
				CreateTime: 1234567890,
			},
			expected: false,
		},
		{
			description: "Update time should be equal",
			x: User{
				UpdateTime: 123456789,
			},
			y: User{
				UpdateTime: 123456789,
			},
			expected: true,
		},
		{
			description: "Update time should not be equal",
			x: User{
				UpdateTime: 123456789,
			},
			y: User{
				UpdateTime: 1234567890,
			},
			expected: false,
		},
		{
			description: "Tenant IDs should be equal",
			x: User{
				AllowedTenantAccess: []string{"ten", "ants"},
			},
			y: User{
				AllowedTenantAccess: []string{"ten", "ants"},
			},
			expected: true,
		},
		{
			description: "Tenant IDs should not be equal",
			x: User{
				AllowedTenantAccess: []string{"ten", "ants"},
			},
			y: User{
				AllowedTenantAccess: []string{"nine", "ants"},
			},
			expected: false,
		},
		{
			description: "Allowed Person IDs should be equal",
			x: User{
				AllowedPersonAccess: []string{"adam", "eve"},
			},
			y: User{
				AllowedPersonAccess: []string{"adam", "eve"},
			},
			expected: true,
		},
		{
			description: "Allowed Person IDs should not be equal",
			x: User{
				AllowedPersonAccess: []string{"adam", "eve"},
			},
			y: User{
				AllowedPersonAccess: []string{"adam", "cat"},
			},
			expected: false,
		},
		{
			description: "scopes should be equal",
			x: User{
				Scopes: []string{"x2", "10x", "1x red-dot"},
			},
			y: User{
				Scopes: []string{"x2", "10x", "1x red-dot"},
			},
			expected: true,
		},
		{
			description: "scopes length should not be equal",
			x: User{
				Scopes: []string{"1x red-dot"},
			},
			y: User{
				Scopes: []string{"1x red-dot", "x2", "10x"},
			},
			expected: false,
		},
		{
			description: "scopes should not be equal",
			x: User{
				Scopes: []string{"x2", "10x", "1x red-dot"},
			},
			y: User{
				Scopes: []string{"10x", "1x red-dot", "x2"},
			},
			expected: false,
		},
		{
			description: "roles should be equal",
			x: User{
				Roles: []string{"cheese", "marmite"},
			},
			y: User{
				Roles: []string{"cheese", "marmite"},
			},
			expected: true,
		},
		{
			description: "roles should not be equal",
			x: User{
				Roles: []string{"cheese", "marmite"},
			},
			y: User{
				Roles: []string{"cheese", "chicken"},
			},
			expected: false,
		},
		{
			description: "roles length should not be equal",
			x: User{
				Roles: []string{"cheese"},
			},
			y: User{
				Roles: []string{"cheese", "chicken and bacon"},
			},
			expected: false,
		},
		{
			description: "personid should be equal",
			x: User{
				PersonID: "socialsecuritynumber",
			},
			y: User{
				PersonID: "socialsecuritynumber",
			},
			expected: true,
		},
		{
			description: "personid should not be equal",
			x: User{
				PersonID: "socialsecuritynumber",
			},
			y: User{
				PersonID: "lol dont ever use a person's social security number",
			},
			expected: false,
		},
		{
			description: "username should be equal",
			x: User{
				Username: "timmy",
			},
			y: User{
				Username: "timmy",
			},
			expected: true,
		},
		{
			description: "username should not be equal",
			x: User{
				Username: "timmy",
			},
			y: User{
				Username: "jimmy",
			},
			expected: false,
		},
		{
			description: "password should be equal",
			x: User{
				Password: "salty",
			},
			y: User{
				Password: "salty",
			},
			expected: true,
		},
		{
			description: "password should not be equal",
			x: User{
				Password: "salty",
			},
			y: User{
				Password: "not-very-salty",
			},
			expected: false,
		},
		{
			description: "firstname should be equal",
			x: User{
				FirstName: "bob lee",
			},
			y: User{
				FirstName: "bob lee",
			},
			expected: true,
		},
		{
			description: "firstname should not be equal",
			x: User{
				FirstName: "bob lee",
			},
			y: User{
				FirstName: "bobby lee",
			},
			expected: false,
		},
		{
			description: "lastname should be equal",
			x: User{
				LastName: "swagger",
			},
			y: User{
				LastName: "swagger",
			},
			expected: true,
		},
		{
			description: "lastname should not be equal",
			x: User{
				LastName: "swagger",
			},
			y: User{
				LastName: "swaggerz",
			},
			expected: false,
		},
		{
			description: "profile uri should be equal",
			x: User{
				ProfileURI: "https://cats.example.com/cat1.jpg",
			},
			y: User{
				ProfileURI: "https://cats.example.com/cat1.jpg",
			},
			expected: true,
		},
		{
			description: "profile uri should not be equal",
			x: User{
				ProfileURI: "https://cats.example.com/cat1.jpg",
			},
			y: User{
				ProfileURI: "https://dogs.example.com/dog1.jpg",
			},
			expected: false,
		},
		{
			description: "disabled should be equal",
			x: User{
				Disabled: false,
			},
			y: User{
				Disabled: false,
			},
			expected: true,
		},
		{
			description: "disabled should not be equal",
			x: User{
				Disabled: false,
			},
			y: User{
				Disabled: true,
			},
			expected: false,
		},
		{
			description: "user should be equal",
			x: User{
				ID:                  "1",
				CreateTime:          123,
				UpdateTime:          321,
				AllowedTenantAccess: []string{"apple", "lettuce"},
				AllowedPersonAccess: []string{"elvis"},
				Scopes:              []string{"10x", "2x"},
				Roles:               []string{"cheese"},
				PersonID:            "123",
				Disabled:            false,
				Username:            "boblee@auth.example.com",
				Password:            "saltypa@ssw0rd",
				FirstName:           "Bob Lee",
				LastName:            "Swagger",
				ProfileURI:          "https://marines.example.com/boblee.png",
			},
			y: User{
				ID:                  "1",
				CreateTime:          123,
				UpdateTime:          321,
				AllowedTenantAccess: []string{"apple", "lettuce"},
				AllowedPersonAccess: []string{"elvis"},
				Scopes:              []string{"10x", "2x"},
				Roles:               []string{"cheese"},
				PersonID:            "123",
				Disabled:            false,
				Username:            "boblee@auth.example.com",
				Password:            "saltypa@ssw0rd",
				FirstName:           "Bob Lee",
				LastName:            "Swagger",
				ProfileURI:          "https://marines.example.com/boblee.png",
			},
			expected: true,
		},
		{
			description: "user should not be equal",
			x: User{
				ID:                  "1",
				CreateTime:          123,
				UpdateTime:          321,
				AllowedTenantAccess: []string{"apple", "lettuce"},
				AllowedPersonAccess: []string{"elvis"},
				Scopes:              []string{"10x", "2x"},
				Roles:               []string{"cheese"},
				PersonID:            "123",
				Disabled:            false,
				Username:            "boblee@auth.example.com",
				Password:            "saltypa@ssw0rd",
				FirstName:           "Bob Lee",
				LastName:            "Swagger",
				ProfileURI:          "https://marines.example.com/boblee.png",
			},
			y: User{
				ID:                  "1",
				CreateTime:          123,
				UpdateTime:          321,
				AllowedTenantAccess: []string{"apple", "lettuce"},
				AllowedPersonAccess: []string{"elvis"},
				Scopes:              []string{"10x"},
				Roles:               []string{"cheese"},
				PersonID:            "123",
				Disabled:            false,
				Username:            "boblee@auth.example.com",
				Password:            "saltypa@ssw0rd",
				FirstName:           "Bob Lee",
				LastName:            "Swagger",
				ProfileURI:          "https://marines.example.com/boblee.png",
			},
			expected: false,
		},
	}

	for _, testcase := range tests {
		assert.Equal(t, testcase.expected, testcase.x.Equal(testcase.y), testcase.description)
	}
}

func TestUser_IsEmpty(t *testing.T) {
	notEmptyUser := User{
		ID: "lol-not-empty",
	}
	assert.Equal(t, notEmptyUser.IsEmpty(), false)

	emptyUser := User{}
	assert.Equal(t, emptyUser.IsEmpty(), true)
}
