package user

import (
	"github.com/MatthewHartstonge/storage/mongo"
	"github.com/imdario/mergo"
	"github.com/ory/fosite"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	if err := c.Find(q).One(&user); err != mgo.ErrNotFound {
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
	q = bson.M{"username": username}
	if err := c.Find(q).One(&user); err != mgo.ErrNotFound {
		return nil, fosite.ErrNotFound
	} else if err != nil {
		return nil, errors.WithStack(err)
	}
	return user, nil
}

// GetUsers returns a map of IDs mapped to a User object that are stored in mongo
func (m *MongoManager) GetUsers(orgid string) (map[string]User, error) {
	c := m.DB.C(mongo.CollectionUsers).With(m.DB.Session.Copy())
	defer c.Database.Session.Close()

	var user *User
	var q bson.M
	q = bson.M{}
	if orgid != "" {
		q = bson.M{"organisation_id": orgid}
	}
	users := make(map[string]User)
	iter := c.Find(q).Limit(100).Iter()
	for iter.Next(&user) {
		users[user.ID] = *user
	}
	if iter.Err() != nil {
		return nil, iter.Err()
	}
	return users, nil
}

// CreateUser stores a new user into mongo
func (m *MongoManager) CreateUser(u *User) error {
	// Ensure unique user
	_, err := m.GetUserByUsername(u.Username)
	if err == mgo.ErrNotFound {
		if u.ID == "" {
			u.ID = uuid.New()
		}
		// Hash incoming secret
		h, err := m.Hasher.Hash([]byte(u.Password))
		if err != nil {
			return errors.WithStack(err)
		}
		u.Password = string(h)
		// Insert new user into mongo
		c := m.DB.C(mongo.CollectionUsers).With(m.DB.Session.Copy())
		defer c.Database.Session.Close()
		if err := c.Insert(u); err != nil {
			return errors.WithStack(err)
		}
	}
	return err
}

// UpdateUser updates a user record. This is done using the equivalent of an object replace.
func (m *MongoManager) UpdateUser(u *User) error {
	o, err := m.GetUser(u.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	// If the password isn't updated, grab it from the stored object
	if u.Password == "" {
		u.Password = string(u.GetHashedSecret())
	} else {
		h, err := m.Hasher.Hash([]byte(u.Password))
		if err != nil {
			return errors.WithStack(err)
		}
		u.Password = string(h)
	}

	// Otherwise, update the object with the new updates
	if err := mergo.Merge(u, o); err != nil {
		return errors.WithStack(err)
	}

	// Update Mongo reference with the updated object
	collection := m.DB.C(mongo.CollectionUsers).With(m.DB.Session.Copy())
	defer collection.Database.Session.Close()
	selector := bson.M{"_id": u.ID}
	if err := collection.Update(selector, u); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// DeleteClient removes a user from mongo
func (m *MongoManager) DeleteUser(id string) error {
	collection := m.DB.C(mongo.CollectionUsers).With(m.DB.Session.Copy())
	defer collection.Database.Session.Close()
	if err := collection.Remove(bson.M{"_id": id}); err != nil {
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
		return errors.WithStack(err)
	}
	isExist := fosite.StringInSlice(scope, u.Scopes)
	if !(isExist) {
		u.Scopes = append(u.Scopes, scope)
		selector := bson.M{"_id": u.ID}
		c.Update(selector, u)
	}
	return nil
}

// RemoveScope gets
func (m *MongoManager) RemoveScopeFromUser(id string, scope string) error {
	c := m.DB.C(mongo.CollectionUsers).With(m.DB.Session.Copy())
	defer c.Database.Session.Close()
	u, err := m.GetUser(id)
	if err != nil {
		return errors.WithStack(err)
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

// Authenticate gets the stored user and authenticates it using a hasher
func (m *MongoManager) Authenticate(id string, secret []byte) (*User, error) {
	u, err := m.GetUser(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := m.Hasher.Compare(u.GetHashedSecret(), secret); err != nil {
		return nil, errors.WithStack(err)
	}
	return u, nil
}
