package client

import (
	"github.com/ory/fosite"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoManager cares for the managing of the Mongo Session instance of a Client
type MongoManager struct {
	DB     *mgo.Database
	Hasher fosite.Hasher
}

// GetConcreteClient finds a Client based on ID and returns it, if found in Mongo.
func (m *MongoManager) GetConcreteClient(id string) (*Client, error) {
	collection := m.DB.C("clients").With(m.DB.Session.Copy())
	defer collection.Database.Session.Close()

	client := &Client{}
	err := collection.Find(bson.M{"_id": id}).One(client)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return client, nil
}

// GetClient returns a Client if found by an ID lookup.
func (m *MongoManager) GetClient(id string) (fosite.Client, error) {
	return m.GetConcreteClient(id)
}
