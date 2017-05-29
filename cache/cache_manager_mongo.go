package cache

import (
	"github.com/ory/fosite"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoManager cares for the managing of the Mongo Session instance for a mongo cache.
type MongoManager struct {
	DB *mgo.Database
}

func (m *MongoManager) getConcreteCacheObject(k string, collectionName string) (value *Cacher, err error) {
	c := m.DB.C(collectionName).With(m.DB.Session.Copy())
	defer c.Database.Session.Close()
	if err := c.Find(bson.M{"_id": k}).One(&value); err == mgo.ErrNotFound {
		return nil, fosite.ErrNotFound
	} else if err != nil {
		return nil, errors.WithStack(err)
	}
	return value, nil
}

// Create provides a way to Create a cache object that has been stored in Mongo. Assumes the struct passed in has bson
// parsing parameters provided in the incoming struct. Namely `_id` must be mapped.
func (m *MongoManager) Create(data Cacher, collectionName string) error {
	c := m.DB.C(collectionName).With(m.DB.Session.Copy())
	defer c.Database.Session.Close()
	if err := c.Insert(data); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Get provides a way to Get a cache object that has been stored in Mongo
func (m *MongoManager) Get(k string, collectionName string) (*Cacher, error) {
	return m.getConcreteCacheObject(k, collectionName)
}

// Update provides a way to Update an old cache object that has been stored in Mongo
func (m *MongoManager) Update(u Cacher, collectionName string) error {
	// Update Mongo reference with the updated object
	collection := m.DB.C(collectionName).With(m.DB.Session.Copy())
	defer collection.Database.Session.Close()
	selector := bson.M{"_id": u.GetKey()}
	if err := collection.Update(selector, u); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete provides a way to Delete a cache object that has been stored in Mongo
func (m *MongoManager) Delete(k string, collectionName string) error {
	collection := m.DB.C(collectionName).With(m.DB.Session.Copy())
	defer collection.Database.Session.Close()
	selector := bson.M{"_id": k}
	if err := collection.Remove(selector); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
