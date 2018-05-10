package user

import (
	// Standard Library Imports
	"strings"
	// External Imports
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/ory/fosite"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	// Internal Imports
	"github.com/MatthewHartstonge/storage/mongo"
)

var (
	ErrUserExists = errors.New("user already exists")
)

// MongoManager manages the Mongo Session instance of a User. Implements user.Manager.
type MongoManager struct {
	// DB is the Mongo connection that holds the base session that can be copied and closed.
	DB     *mgo.Database
	Hasher fosite.Hasher
}

// GetUser gets a user document that has been previously stored in mongo
func (m *MongoManager) GetUser(id string) (*User, error) {
	c := m.DB.C(mongo.CollectionUsers).With(m.DB.Session.Copy())
	defer c.Database.Session.Close()

	var user *User
	var q bson.M
	q = bson.M{"_id": id}
	if err := c.Find(q).One(&user); err == mgo.ErrNotFound {
		return nil, fosite.ErrNotFound
	} else if err != nil {
		return nil, errors.WithStack(err)
	}
	return user, nil
}

// GetUserByUsername gets a user document by searching for a username that has been previously stored in mongo
func (m *MongoManager) GetUserByUsername(username string) (*User, error) {
	c := m.DB.C(mongo.CollectionUsers).With(m.DB.Session.Copy())
	defer c.Database.Session.Close()

	var user *User
	var q bson.M
	q = bson.M{"username": strings.ToLower(username)}
	if err := c.Find(q).One(&user); err != nil {
		if err == mgo.ErrNotFound {
			return nil, fosite.ErrNotFound
		}
		return nil, errors.WithStack(err)
	}
	return user, nil
}

// Filter enables querying MongoDB for specific user accounts.
type Filter struct {
	// AllowedTenantAccess filters users based on Tenant Access.
	AllowedTenantAccess string
	// AllowedPeopleAccess filters users based on People Access.
	AllowedPeopleAccess string
	// PersonID filters users based on People ID.
	PersonID string
	// Username filters users based on username.
	Username string
	// Scopes filters users based on scopes users must have.
	// Scopes performs an AND operation. To obtain OR, do multiple requests with a single scope.
	Scopes []string
	// FirstName filters users based on their First Name.
	FirstName string
	// LastName filters users based on their Last Name.
	LastName string
	// Disabled filters users to those with disabled accounts.
	Disabled bool
}

// GetUsers returns a map of IDs mapped to a User object that are stored in mongo
func (m *MongoManager) GetUsers(filters Filter) (map[string]User, error) {
	c := m.DB.C(mongo.CollectionUsers).With(m.DB.Session.Copy())
	defer c.Database.Session.Close()

	var user *User
	var q bson.M
	q = bson.M{}
	if filters.AllowedTenantAccess != "" {
		q["allowedTenantAccess"] = filters.AllowedTenantAccess
	}
	if filters.AllowedPeopleAccess != "" {
		q["allowedPeopleAccess"] = filters.AllowedPeopleAccess
	}
	if len(filters.Scopes) > 0 {
		q["scopes"] = bson.M{"$all": filters.Scopes}
	}
	if filters.PersonID != "" {
		q["personID"] = filters.PersonID
	}
	if filters.Username != "" {
		q["username"] = filters.Username
	}
	if filters.FirstName != "" {
		q["firstName"] = filters.FirstName
	}
	if filters.LastName != "" {
		q["lastName"] = filters.LastName
	}
	if filters.Disabled {
		q["disabled"] = filters.Disabled
	}

	users := make(map[string]User)
	iter := c.Find(q).Limit(100).Iter()
	for iter.Next(&user) {
		users[user.ID] = *user
	}
	if iter.Err() != nil {
		return nil, errors.WithStack(iter.Err())
	}
	return users, nil
}

// CreateUser stores a new user into mongo
func (m *MongoManager) CreateUser(u *User) error {
	// Ensure unique user
	usr, err := m.GetUserByUsername(strings.ToLower(u.Username))
	if err == nil && !usr.IsEmpty() {
		return ErrUserExists
	}
	if err != fosite.ErrNotFound {
		return err
	}

	if u.ID == "" || uuid.Parse(u.ID) == nil {
		u.ID = uuid.New()
	}

	// Hash incoming secret
	h, err := m.Hasher.Hash([]byte(u.Password))
	if err != nil {
		return err
	}

	u.Password = string(h)
	u.Username = strings.ToLower(u.Username)

	// Insert new user into mongo
	c := m.DB.C(mongo.CollectionUsers).With(m.DB.Session.Copy())
	defer c.Database.Session.Close()
	if err := c.Insert(u); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateUser updates a user record. This is done using the equivalent of an object replace.
func (m *MongoManager) UpdateUser(user *User) error {
	old, err := m.GetUser(user.ID)
	if err != nil {
		return err
	}

	// If the password isn't updated, grab it from the stored object
	if user.Password == "" {
		user.Password = string(old.GetHashedSecret())
	} else {
		h, err := m.Hasher.Hash([]byte(user.Password))
		if err != nil {
			return err
		}
		user.Password = string(h)
	}

	// Update Mongo reference with the updated object
	collection := m.DB.C(mongo.CollectionUsers).With(m.DB.Session.Copy())
	defer collection.Database.Session.Close()
	selector := bson.M{"_id": user.ID}
	if err := collection.Update(selector, user); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// DeleteUser removes a user from mongo
func (m *MongoManager) DeleteUser(id string) error {
	collection := m.DB.C(mongo.CollectionUsers).With(m.DB.Session.Copy())
	defer collection.Database.Session.Close()
	if err := collection.Remove(bson.M{"_id": id}); err != nil {
		if err == mgo.ErrNotFound {
			return fosite.ErrNotFound
		}
		return errors.WithStack(err)
	}
	return nil
}

// GrantScopeToUser adds a scope to a user if it doesn't already exist in the mongo record
func (m *MongoManager) GrantScopeToUser(id string, scope string) error {
	c := m.DB.C(mongo.CollectionUsers).With(m.DB.Session.Copy())
	defer c.Database.Session.Close()
	u, err := m.GetUser(id)
	if err != nil {
		return err
	}
	isExist := fosite.StringInSlice(scope, u.Scopes)
	if !(isExist) {
		u.Scopes = append(u.Scopes, scope)
		selector := bson.M{"_id": u.ID}
		c.Update(selector, u)
	}
	return nil
}

// RemoveScopeFromUser takes a scoped right away from the given user.
func (m *MongoManager) RemoveScopeFromUser(id string, scope string) error {
	c := m.DB.C(mongo.CollectionUsers).With(m.DB.Session.Copy())
	defer c.Database.Session.Close()
	u, err := m.GetUser(id)
	if err != nil {
		return err
	}
	for i, s := range u.Scopes {
		if scope == s {
			u.Scopes = append(u.Scopes[:i], u.Scopes[i+1:]...)
			selector := bson.M{"_id": u.ID}
			c.Update(selector, u)
			break
		}
	}
	return nil
}

// Authenticate wraps AuthenticateByUsername to allow users to be found via their username. Returns a user record
// if authentication is successful.
func (m *MongoManager) Authenticate(username string, secret []byte) (*User, error) {
	return m.AuthenticateByUsername(username, secret)
}

// AuthenticateByID gets the stored user by ID and authenticates it using a hasher
func (m *MongoManager) AuthenticateByID(id string, secret []byte) (*User, error) {
	u, err := m.GetUser(id)
	if err != nil {
		return nil, err
	}

	if err := m.Hasher.Compare(u.GetHashedSecret(), secret); err != nil {
		return nil, err
	}
	return u, nil
}

// AuthenticateByUsername gets the stored user by username and authenticates it using a hasher
func (m *MongoManager) AuthenticateByUsername(username string, secret []byte) (*User, error) {
	u, err := m.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if err := m.Hasher.Compare(u.GetHashedSecret(), secret); err != nil {
		return nil, err
	}
	return u, nil
}
