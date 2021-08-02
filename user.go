package storage

import (
	// Standard Library Imports
	"context"
	"fmt"
	"strings"

	// External Imports
	"github.com/ory/fosite"
)

// User provides the specific types for storing, editing, deleting and
// retrieving a User record in mongo.
type User struct {
	//// User Meta
	// ID is the uniquely assigned uuid that references the user
	ID string `bson:"id" json:"id" xml:"id"`

	// createTime is when the resource was created in seconds from the epoch.
	CreateTime int64 `bson:"createTime" json:"createTime" xml:"createTime"`

	// updateTime is the last time the resource was modified in seconds from
	// the epoch.
	UpdateTime int64 `bson:"updateTime" json:"updateTime" xml:"updateTime"`

	// AllowedTenantAccess contains the Tenant IDs that the user has been given
	// rights to access.
	// This helps in multi-tenanted situations where a user can be given
	// explicit cross-tenant access.
	AllowedTenantAccess []string `bson:"allowedTenantAccess" json:"allowedTenantAccess,omitempty" xml:"allowedTenantAccess,omitempty"`

	// AllowedPersonAccess contains a list of Person IDs that the user is
	// allowed access to.
	// This helps in multi-tenanted situations where a user can be given
	// explicit access to other people accounts, for example, parents to
	// children records.
	AllowedPersonAccess []string `bson:"allowedPersonAccess" json:"allowedPersonAccess,omitempty" xml:"allowedPersonAccess,omitempty"`

	// Scopes contains the permissions that the user is entitled to request.
	Scopes []string `bson:"scopes" json:"scopes" xml:"scopes"`

	// PersonID is a uniquely assigned id that references a person within the
	// system.
	// This enables applications where an external person data store is present.
	// This helps in multi-tenanted situations where the person is unique, but
	// the underlying user accounts can exist per tenant.
	PersonID string `bson:"personId" json:"personId" xml:"personId"`

	// Disabled specifies whether the user has been disallowed from signing in
	Disabled bool `bson:"disabled" json:"disabled" xml:"disabled"`

	//// User Content
	// Username is used to authenticate a user
	Username string `bson:"username" json:"username" xml:"username"`

	// Password of the user - will be a hash based on your fosite selected
	// hasher.
	// If using this model directly in an API, be sure to clear the password
	// out when marshaling to json/xml.
	Password string `bson:"password,omitempty" json:"password,omitempty" xml:"password,omitempty"`

	// FirstName stores the user's Last Name
	FirstName string `bson:"firstName" json:"firstName" xml:"firstName"`

	// LastName stores the user's Last Name
	LastName string `bson:"lastName" json:"lastName" xml:"lastName"`

	// ProfileURI is a pointer to where their profile picture lives
	ProfileURI string `bson:"profileUri" json:"profileUri,omitempty" xml:"profileUri,omitempty"`
}

// FullName concatenates the User's First Name and Last Name for templating
// purposes
func (u User) FullName() (fullName string) {
	return strings.TrimSpace(fmt.Sprintf("%s %s", u.FirstName, u.LastName))
}

// SetPassword takes a cleartext secret, hashes it with a hasher and sets it as
// the user's password
func (u *User) SetPassword(cleartext string, hasher fosite.Hasher) (err error) {
	h, err := hasher.Hash(context.TODO(), []byte(cleartext))
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
	return hasher.Compare(context.TODO(), u.GetHashedSecret(), []byte(cleartext))
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
func (u *User) EnablePeopleAccess(personIDs ...string) {
	for i := range personIDs {
		found := false
		for j := range u.AllowedPersonAccess {
			if personIDs[i] == u.AllowedPersonAccess[j] {
				found = true
				break
			}
		}
		if !found {
			u.AllowedPersonAccess = append(u.AllowedPersonAccess, personIDs[i])
		}
	}
}

// DisablePeopleAccess disables user access to the provided people.
func (u *User) DisablePeopleAccess(personIDs ...string) {
	for i := range personIDs {
		for j := range u.AllowedPersonAccess {
			if personIDs[i] == u.AllowedPersonAccess[j] {
				copy(u.AllowedPersonAccess[j:], u.AllowedPersonAccess[j+1:])
				u.AllowedPersonAccess[len(u.AllowedPersonAccess)-1] = ""
				u.AllowedPersonAccess = u.AllowedPersonAccess[:len(u.AllowedPersonAccess)-1]
				break
			}
		}
	}
}

// EnableScopeAccess enables user access to one or many scopes.
func (u *User) EnableScopeAccess(scopes ...string) {
	for i := range scopes {
		found := false
		for j := range u.Scopes {
			if scopes[i] == u.Scopes[j] {
				found = true
				break
			}
		}
		if !found {
			u.Scopes = append(u.Scopes, scopes[i])
		}
	}
}

// DisableScopeAccess disables user access to one or many scopes.
func (u *User) DisableScopeAccess(scopes ...string) {
	for i := range scopes {
		for j := range u.Scopes {
			if scopes[i] == u.Scopes[j] {
				copy(u.Scopes[j:], u.Scopes[j+1:])
				u.Scopes[len(u.Scopes)-1] = ""
				u.Scopes = u.Scopes[:len(u.Scopes)-1]
				break
			}
		}
	}
}

// Equal enables checking equality as having a byte array in a struct stops
// allowing direct equality checks.
func (u User) Equal(x User) bool {
	if u.ID != x.ID {
		return false
	}

	if u.CreateTime != x.CreateTime {
		return false
	}

	if u.UpdateTime != x.UpdateTime {
		return false
	}

	if !stringArrayEquals(u.AllowedTenantAccess, x.AllowedTenantAccess) {
		return false
	}

	if !stringArrayEquals(u.AllowedPersonAccess, x.AllowedPersonAccess) {
		return false
	}

	if !stringArrayEquals(u.Scopes, x.Scopes) {
		return false
	}

	if u.PersonID != x.PersonID {
		return false
	}

	if u.Disabled != x.Disabled {
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

	return true
}

// IsEmpty returns true if the current user holds no data.
func (u User) IsEmpty() bool {
	return u.Equal(User{})
}
