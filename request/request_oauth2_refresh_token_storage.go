package request

import (
	"context"
	"github.com/ory/fosite"
)

/* These functions provide a concrete implementation of fosite.RefreshTokenStorage */

// CreateRefreshTokenSession stores a new Refresh Token Session in mongo
func (m *MongoManager) CreateRefreshTokenSession(_ context.Context, signature string, request fosite.Requester) (err error) {
	return m.createSession(signature, request, mongoCollectionRefreshTokens)
}

func (m *MongoManager) GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	return
}

func (m *MongoManager) DeleteRefreshTokenSession(ctx context.Context, signature string) (err error) {
	return
}
