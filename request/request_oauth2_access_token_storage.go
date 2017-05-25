package request

import (
	"context"
	"github.com/ory/fosite"
)

/* These functions provide a concrete implementation of fosite.AccessTokenStorage */

func (m *MongoManager) CreateAccessTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	return
}

func (m MongoManager) GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	return
}

func (m *MongoManager) DeleteAccessTokenSession(ctx context.Context, signature string) (err error) {
	return
}
