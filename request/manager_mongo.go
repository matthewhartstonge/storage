package request

import (
	"encoding/json"
	"github.com/ory/fosite"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

const (
	mongoCollectionOpenIDSessions    = "OpenIDConnectSessions"
	mongoCollectionAccessTokens      = "AccessTokens"
	mongoCollectionRefreshTokens     = "RefreshTokens"
	mongoCollectionAuthorizationCode = "AuthorizationCode"
)

// MongoManager manages the main Mongo Session for a Request.
type MongoManager struct {
	// DB is the Mongo connection that holds the base session that can be copied and closed.
	DB *mgo.Database

	// TODO: Add AES cipher for Token Encryption?
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
	if err := m.DB.C(collectionName).Insert(data); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
