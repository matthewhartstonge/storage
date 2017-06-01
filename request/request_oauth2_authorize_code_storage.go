package request

import (
	"context"
	"github.com/MatthewHartstonge/storage/mongo"
	"github.com/ory/fosite"
)

/* These functions provide a concrete implementation of fosite.AuthorizeCodeStorage */

// CreateAuthorizeCodeSession creates a new session for an authorize code grant in mongo
func (m *MongoManager) CreateAuthorizeCodeSession(_ context.Context, code string, request fosite.Requester) (err error) {
	return m.createSession(code, request, mongo.CollectionAuthorizationCodes)
}

// GetAuthorizeCodeSession finds an authorize code grant session in mongo
func (m MongoManager) GetAuthorizeCodeSession(_ context.Context, code string, session fosite.Session) (request fosite.Requester, err error) {
	return m.findSessionBySignature(code, session, mongo.CollectionAuthorizationCodes)
}

// DeleteAuthorizeCodeSession removes an authorize code session from mongo
func (m *MongoManager) DeleteAuthorizeCodeSession(_ context.Context, code string) (err error) {
	return m.deleteSession(code, mongo.CollectionAuthorizationCodes)
}
