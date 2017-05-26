package user

import (
	"fmt"
	"github.com/ory/fosite"
)

type User struct {
	// ID is the uniquely assigned uuid that references the user
	ID string `bson:"_id" json:"id" xml:"id"`

	// The organisation the user belongs to
	OrganisationID string `bson:"organisation_id,omitempty" json:"organisation_id,omitempty" xml:"organisation_id,omitempty"`

	// Username is used to authenticate a user
	Username string `bson:"username" json:"username" xml:"username"`

	//
	Password string `bson:"password" json:"-" xml:"-"`

	// Scopes contains the scopes that have been granted to
	Scopes []string `bson:"scopes" json:"scopes" xml:"scopes"`

	// FirstName stores the user's Last Name
	FirstName string `bson:"first_name" json:"first_name" xml:"first_name"`

	// LastName stores the user's Last Name
	LastName string `bson:"last_name" json:"last_name" xml:"last_name"`

	// ProfileURI is a pointer to where their profile picture lives
	ProfileURI string `bson:"profile_uri" json:"profile_uri,omitempty" xml:"profile_uri,omitempty"`
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
