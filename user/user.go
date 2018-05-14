package user

import (
	// Standard Library Imports
	"fmt"

	// External Imports
	"github.com/ory/fosite"
)

// User provides the specific types for storing, editing, deleting and retrieving a User record in mongo.
type User struct {
	// User Meta
	// ID is the uniquely assigned uuid that references the user
	ID string `bson:"_id" json:"id" xml:"id"`

	// AllowedTenantAccess contains the Tenant IDs that the user has been given rights to access.
	// This helps in multi-tenanted situations where a user can be given explicit cross-tenant access.
	AllowedTenantAccess []string `bson:"allowedTenantAccess,omitempty" json:"allowedTenantAccess,omitempty" xml:"allowedTenantAccess,omitempty"`

	// AllowedPeopleAccess contains People IDs that users are allowed access to.
	// This helps in multi-tenanted situations where a user can be given explicit access to other people accounts, for
	// example, parents to children records.
	AllowedPeopleAccess []string `bson:"allowedPeopleAccess" json:"allowedPeopleAccess" xml:"allowedPeopleAccess"`

	// Scopes contains the scopes that have been granted to
	Scopes []string `bson:"scopes" json:"scopes" xml:"scopes"`

	// PersonID is a uniquely assigned uuid that references a person within the system.
	// This enables applications where an external person data store is present. This helps in multi-tenanted
	// situations where the person is unique, but the underlying user accounts can exist per tenant.
	PersonID string `bson:"personID" json:"personID" xml:"personID"`

	// User Content
	// Username is used to authenticate a user
	Username string `bson:"username" json:"username" xml:"username"`

	// Password of the user - will be a hash based on your fosite selected hasher
	// If using this model directly in an API, be sure to clear the password out when marshaling to json/xml
	Password string `bson:"password,omitempty" json:"password,omitempty" xml:"password,omitempty"`

	// FirstName stores the user's Last Name
	FirstName string `bson:"firstName" json:"firstName" xml:"firstName"`

	// LastName stores the user's Last Name
	LastName string `bson:"lastName" json:"lastName" xml:"lastName"`

	// ProfileURI is a pointer to where their profile picture lives
	ProfileURI string `bson:"profileUri" json:"profileUri,omitempty" xml:"profileUri,omitempty"`

	// Disabled specifies whether the user has been disallowed from signing in
	Disabled bool `bson:"disabled" json:"disabled" xml:"disabled"`
}

// GetFullName concatenates the User's First Name and Last Name for templating purposes
func (u User) GetFullName() (fn string) {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

// SetPassword takes a cleartext secret, hashes it with a hasher and sets it as the user's password
func (u *User) SetPassword(cleartext string, hasher fosite.Hasher) (err error) {
	h, err := hasher.Hash([]byte(cleartext))
	if err != nil {
		return err
	}
	u.Password = string(h)
	return nil
}

// GetHashedSecret returns the Users's Hashed Secret as a byte array
func (u *User) GetHashedSecret() []byte {
	return []byte(u.Password)
}

// Authenticate compares a cleartext string against the user's
func (u User) Authenticate(cleartext string, hasher fosite.Hasher) error {
	return hasher.Compare(u.GetHashedSecret(), []byte(cleartext))
}

// EnableTenantAccess enables user access to one or many tenants.
func (u *User) EnableTenantAccess(tenantIDs ...string) {
	for i := range tenantIDs {
		found := false
		for j := range u.AllowedTenantAccess {
			if tenantIDs[i] == u.AllowedTenantAccess[j] {
				found = true
				break
			}
		}
		if !found {
			u.AllowedTenantAccess = append(u.AllowedTenantAccess, tenantIDs[i])
		}
	}
}

// DisableTenantAccess disables user access to one or many tenants.
func (u *User) DisableTenantAccess(tenantIDs ...string) {
	for i := range tenantIDs {
		for j := range u.AllowedTenantAccess {
			if tenantIDs[i] == u.AllowedTenantAccess[j] {
				copy(u.AllowedTenantAccess[j:], u.AllowedTenantAccess[j+1:])
				u.AllowedTenantAccess[len(u.AllowedTenantAccess)-1] = ""
				u.AllowedTenantAccess = u.AllowedTenantAccess[:len(u.AllowedTenantAccess)-1]
				break
			}
		}
	}
}

// EnablePeopleAccess enables user access to the provided people
func (u *User) EnablePeopleAccess(peopleIDs ...string) {
	for i := range peopleIDs {
		found := false
		for j := range u.AllowedPeopleAccess {
			if peopleIDs[i] == u.AllowedPeopleAccess[j] {
				found = true
				break
			}
		}
		if !found {
			u.AllowedPeopleAccess = append(u.AllowedPeopleAccess, peopleIDs[i])
		}
	}
}

// DisablePeopleAccess disables user access to the provided people.
func (u *User) DisablePeopleAccess(peopleIDs ...string) {
	for i := range peopleIDs {
		for j := range u.AllowedPeopleAccess {
			if peopleIDs[i] == u.AllowedPeopleAccess[j] {
				copy(u.AllowedPeopleAccess[j:], u.AllowedPeopleAccess[j+1:])
				u.AllowedPeopleAccess[len(u.AllowedPeopleAccess)-1] = ""
				u.AllowedPeopleAccess = u.AllowedPeopleAccess[:len(u.AllowedPeopleAccess)-1]
				break
			}
		}
	}
}

// EnableScopeAccess enables user access to one or many scopes.
func (u *User) EnableScopeAccess(addScopes ...string) {
	for i := range addScopes {
		found := false
		for j := range u.Scopes {
			if addScopes[i] == u.Scopes[j] {
				found = true
				break
			}
		}
		if !found {
			u.Scopes = append(u.Scopes, addScopes[i])
		}
	}
}

// DisableScopeAccess disables user access to one or many scopes.
func (u *User) DisableScopeAccess(removeScopes ...string) {
	for i := range removeScopes {
		for j := range u.Scopes {
			if removeScopes[i] == u.Scopes[j] {
				copy(u.Scopes[j:], u.Scopes[j+1:])
				u.Scopes[len(u.Scopes)-1] = ""
				u.Scopes = u.Scopes[:len(u.Scopes)-1]
				break
			}
		}
	}
}

// Equal enables checking equality as having a byte array in a struct stop allowing equality checks.
func (u User) Equal(x User) bool {
	if u.ID != x.ID {
		return false
	}

	if !strArrEquals(u.AllowedTenantAccess, x.AllowedTenantAccess) {
		return false
	}

	if !strArrEquals(u.AllowedPeopleAccess, x.AllowedPeopleAccess) {
		return false
	}

	if !strArrEquals(u.Scopes, x.Scopes) {
		return false
	}

	if u.PersonID != x.PersonID {
		return false
	}

	if u.Username != x.Username {
		return false
	}

	if u.Password != x.Password {
		return false
	}

	if u.FirstName != x.FirstName {
		return false
	}

	if u.LastName != x.LastName {
		return false
	}

	if u.ProfileURI != x.ProfileURI {
		return false
	}

	if u.Disabled != x.Disabled {
		return false
	}

	return true
}

func (u User) IsEmpty() bool {
	return u.Equal(User{})
}

func strArrEquals(arr1 []string, arr2 []string) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	for i := range arr1 {
		if arr1[i] != arr2[i] {
			return false
		}
	}

	return true
}
