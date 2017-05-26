package request

import (
	"encoding/json"
	"github.com/MatthewHartstonge/storage/client"
	"github.com/MatthewHartstonge/storage/user"
	"github.com/ory/fosite"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoManager manages the main Mongo Session for a Request.
type MongoManager struct {
	// DB is the Mongo connection that holds the base session that can be copied and closed.
	DB *mgo.Database

	// Due to the nature of an OAuth request, it will need to cross reference the Client collections.
	Clients *client.MongoManager

	// For the Password Credentials Grant, A user MongoManager is required in order to find and authenticate users.
	Users *user.MongoManager
}

// Given a request from fosite, marshals to a form that enables storing the request in mongo
func mongoCollectionFromRequest(signature string, r fosite.Requester) (*mongoRequestData, error) {
	session, err := json.Marshal(r.GetSession())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &mongoRequestData{
		ID:            r.GetID(),
		RequestedAt:   r.GetRequestedAt(),
		Signature:     signature,
		ClientID:      r.GetClient().GetID(),
		Scopes:        r.GetRequestedScopes(),
		GrantedScopes: r.GetGrantedScopes(),
		Form:          r.GetRequestForm().Encode(),
		Session:       session,
	}, nil

}

// createSession stores a session to a specific mongo collection
func (m *MongoManager) createSession(signature string, requester fosite.Requester, collectionName string) error {
	data, err := mongoCollectionFromRequest(signature, requester)
	if err != nil {
		return err
	}

	c := m.DB.C(collectionName).With(m.DB.Session.Copy())
	defer c.Database.Session.Close()
	if err := c.Insert(data); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// findSessionBySignature finds a session within a specific mongo collection
func (m *MongoManager) findSessionBySignature(signature string, session fosite.Session, collectionName string) (fosite.Requester, error) {
	var d *mongoRequestData
	c := m.DB.C(collectionName).With(m.DB.Session.Copy())
	defer c.Database.Session.Close()
	if err := c.Find(bson.M{"signature": signature}).One(d); err == mgo.ErrNotFound {
		return nil, fosite.ErrNotFound
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	return d.toRequest(session, m.Clients)
}

// deleteSession removes a session document from a specfic mongo collection
func (m *MongoManager) deleteSession(signature string, collectionName string) error {
	c := m.DB.C(collectionName).With(m.DB.Session.Copy())
	defer c.Database.Session.Close()
	if err := c.Remove(bson.M{"signature": signature}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
