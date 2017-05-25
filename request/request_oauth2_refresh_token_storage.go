package request

import (
	"context"
	"github.com/ory/fosite"
)

/* These functions provide a concrete implementation of fosite.RefreshTokenStorage */

func (m *MongoManager) CreateRefreshTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	return
}

func (m *MongoManager) GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	return
}

func (m *MongoManager) DeleteRefreshTokenSession(ctx context.Context, signature string) (err error) {
	return
}
