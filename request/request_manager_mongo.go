package request

import (
	"context"
	"encoding/json"

	"github.com/MatthewHartstonge/storage/cache"
	"github.com/MatthewHartstonge/storage/client"
	"github.com/MatthewHartstonge/storage/mongo"
	"github.com/MatthewHartstonge/storage/user"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/ory/fosite"
	"github.com/pkg/errors"
)

// MongoManager manages the main Mongo Session for a Request.
type MongoManager struct {
	// DB is the Mongo connection that holds the base session that can be copied and closed.
	DB *mgo.Database

	// In order to create, read, update and delete from the caching database, a CacheManager is required.
	Cache *cache.MongoManager

	// Due to the nature of an OAuth request, it will need to cross reference the Client collections.
	Clients *client.MongoManager

	// For the Password Credentials Grant, A user MongoManager is required in order to find and authenticate users.
	Users *user.MongoManager
}

// Given a request from fosite, marshals to a form that enables storing the request in mongo
func mongoCollectionFromRequest(signature string, r fosite.Requester) (*MongoRequest, error) {
	session, err := json.Marshal(r.GetSession())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &MongoRequest{
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
	c := m.DB.C(collectionName).With(m.DB.Session.Copy())
	defer c.Database.Session.Close()
	mongoData := &MongoRequest{}
	if err := c.Find(bson.M{"signature": signature}).One(mongoData); err == mgo.ErrNotFound {
		return nil, fosite.ErrNotFound
	} else if err != nil {
		return nil, errors.WithStack(err)
	}
	return mongoData.toRequest(session, m.Clients)
}

// deleteSession removes a session document from a specific mongo collection
func (m *MongoManager) deleteSession(signature string, collectionName string) error {
	c := m.DB.C(collectionName).With(m.DB.Session.Copy())
	defer c.Database.Session.Close()
	if err := c.Remove(bson.M{"signature": signature}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// RevokeRefreshToken finds a token stored in cache based on request ID and deletes the session by signature.
func (m *MongoManager) RevokeRefreshToken(ctx context.Context, requestID string) error {
	o, err := m.Cache.Get(requestID, mongo.CollectionCacheRefreshTokens)
	if err != nil {
		return err
	}
	err = m.DeleteRefreshTokenSession(ctx, o.GetValue())
	if err != nil {
		return err
	}
	return m.Cache.Delete(o.GetKey(), mongo.CollectionCacheRefreshTokens)
}

// RevokeAccessToken finds a token stored in cache based on request ID and deletes the session by signature.
func (m *MongoManager) RevokeAccessToken(ctx context.Context, requestID string) error {
	o, err := m.Cache.Get(requestID, mongo.CollectionCacheAccessTokens)
	if err != nil {
		return err
	}
	err = m.DeleteAccessTokenSession(ctx, o.GetValue())
	if err != nil {
		return err
	}
	return m.Cache.Delete(o.GetKey(), mongo.CollectionCacheAccessTokens)
}
