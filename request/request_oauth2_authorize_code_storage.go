package request

import (
	"context"
	"github.com/ory/fosite"
)

/* These functions provide a concrete implementation of fosite.AuthorizeCodeStorage */

func (m *MongoManager) CreateAuthorizeCodeSession(ctx context.Context, code string, request fosite.Requester) (err error) {
	return
}

func (m MongoManager) GetAuthorizeCodeSession(ctx context.Context, code string, session fosite.Session) (request fosite.Requester, err error) {
	return
}

func (m *MongoManager) DeleteAuthorizeCodeSession(ctx context.Context, code string) (err error) {
	return
}
