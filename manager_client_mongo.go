package storage

import (
	"github.com/MatthewHartstonge/storage/model"
	"github.com/ory/fosite"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ClientManager struct {
	Hasher fosite.Hasher
	DB     *mgo.Database
}

func (c *ClientManager) GetConcreteClient(id string) (*model.Client, error) {
	client := &model.Client{}
	collection := c.DB.C("clients").With(c.DB.Session.Copy())
	err := collection.Find(bson.M{"_id": id}).One(client)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return client, nil
}

func (c *ClientManager) GetClient(id string) (fosite.Client, error) {
	return c.GetConcreteClient(id)
}
