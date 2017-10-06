package client

import (
	// Standard Library Imports
	"context"

	// External Imports
	"github.com/imdario/mergo"
	"github.com/ory/fosite"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	// Internal Imports
	"github.com/MatthewHartstonge/storage/mongo"
)

var (
	ErrClientExists = errors.New("client already exists")
)

// MongoManager cares for the managing of the Mongo Session instance of a Client.
type MongoManager struct {
	DB     *mgo.Database
	Hasher fosite.Hasher
}

// GetConcreteClient finds a Client based on ID and returns it, if found in Mongo.
func (m MongoManager) GetConcreteClient(id string) (*Client, error) {
	collection := m.DB.C(mongo.CollectionClients).With(m.DB.Session.Copy())
	defer collection.Database.Session.Close()

	client := &Client{}
	if err := collection.Find(bson.M{"_id": id}).One(client); err != nil {
		if err == mgo.ErrNotFound {
			return nil, fosite.ErrNotFound
		}

		return nil, errors.WithStack(err)
	}
	return client, nil
}

// GetClient returns a Client if found by an ID lookup.
func (m MongoManager) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	return m.GetConcreteClient(id)
}

// GetClients returns a map of clients mapped by client ID
func (m MongoManager) GetClients() (clients map[string]Client, err error) {
	clients = make(map[string]Client)
	collection := m.DB.C(mongo.CollectionClients).With(m.DB.Session.Copy())
	defer collection.Database.Session.Close()

	var result *Client
	iter := collection.Find(bson.M{}).Limit(100).Iter()
	for iter.Next(&result) {
		clients[result.ID] = *result
	}
	if iter.Err() != nil {
		return nil, iter.Err()
	}
	return
}

// CreateClient adds a new OAuth2.0 Client to the client store.
func (m *MongoManager) CreateClient(c *Client) error {
	if c.ID == "" {
		c.ID = uuid.New()
	}

	// Hash incoming secret
	h, err := m.Hasher.Hash(c.Secret)
	if err != nil {
		return errors.WithStack(err)
	}
	c.Secret = h

	// Insert to Mongo
	collection := m.DB.C(mongo.CollectionClients).With(m.DB.Session.Copy())
	defer collection.Database.Session.Close()
	if err := collection.Insert(c); err != nil {
		if mgo.IsDup(err) {
			return ErrClientExists
		}

		return errors.WithStack(err)
	}
	return nil
}

// UpdateClient updates an OAuth 2.0 Client record. This is done using the equivalent of an object replace.
func (m *MongoManager) UpdateClient(c *Client) error {
	o, err := m.GetConcreteClient(c.ID)
	if err != nil {
		if err == fosite.ErrNotFound {
			return err
		}
		return errors.WithStack(err)
	}

	// If the password isn't updated, grab it from the stored object
	if string(c.Secret) == "" {
		c.Secret = o.GetHashedSecret()
	} else {
		h, err := m.Hasher.Hash(c.Secret)
		if err != nil {
			return errors.WithStack(err)
		}
		c.Secret = h
	}

	// Otherwise, update the object with the new updates
	if err := mergo.Merge(c, o); err != nil {
		return errors.WithStack(err)
	}

	// Update Mongo reference with the updated object
	collection := m.DB.C(mongo.CollectionClients).With(m.DB.Session.Copy())
	defer collection.Database.Session.Close()
	selector := bson.M{"_id": c.ID}
	if err := collection.Update(selector, c); err != nil {
		if err == mgo.ErrNotFound {
			return fosite.ErrNotFound
		}

		return errors.WithStack(err)
	}
	return nil
}

// DeleteClient removes an OAuth 2.0 Client from the client store
func (m *MongoManager) DeleteClient(id string) error {
	collection := m.DB.C(mongo.CollectionClients).With(m.DB.Session.Copy())
	defer collection.Database.Session.Close()
	if err := collection.Remove(bson.M{"_id": id}); err != nil {
		if err == mgo.ErrNotFound {
			return fosite.ErrNotFound
		}

		return errors.WithStack(err)
	}
	return nil
}

// Authenticate compares a client secret with the client's stored hashed secret
func (m MongoManager) Authenticate(id string, secret []byte) (*Client, error) {
	c, err := m.GetConcreteClient(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := m.Hasher.Compare(c.GetHashedSecret(), secret); err != nil {
		return nil, errors.WithStack(err)
	}
	return c, nil
}
