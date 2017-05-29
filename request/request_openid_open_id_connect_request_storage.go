package request

import (
	"context"
	"github.com/MatthewHartstonge/storage/mongo"
	"github.com/ory/fosite"
)

/* These functions provide a concrete implementation of openid.OpenIDConnectRequestStorage */

// CreateOpenIDConnectSession creates an open id connect session for a given authorize code in mongo. This is relevant
// for explicit open id connect flow.
func (m *MongoManager) CreateOpenIDConnectSession(ctx context.Context, authorizeCode string, requester fosite.Requester) (err error) {
	return m.createSession(authorizeCode, requester, mongo.CollectionOpenIDSessions)
}

// GetOpenIDConnectSession gets a session based off the Authorize Code and returns a fosite.Requester which contains a
// session or an error.
func (m *MongoManager) GetOpenIDConnectSession(ctx context.Context, authorizeCode string, requester fosite.Requester) (req fosite.Requester, err error) {
	session := requester.GetSession()
	if session == nil {
		return nil, fosite.ErrNotFound
	}
	return m.findSessionBySignature(authorizeCode, session, mongo.CollectionOpenIDSessions)
}

// DeleteOpenIDConnectSession removes an open id connect session from mongo.
func (m *MongoManager) DeleteOpenIDConnectSession(ctx context.Context, authorizeCode string) (err error) {
	return m.deleteSession(authorizeCode, mongo.CollectionOpenIDSessions)
}
