package request

import (
	"context"
	"github.com/ory/fosite"
)

/* These functions provide a concrete implementation of fosite.ImplicitGrantStorage */

// CreateImplicitAccessTokenSession stores an implicit access token based session in mongo
func (m *MongoManager) CreateImplicitAccessTokenSession(ctx context.Context, token string, request fosite.Requester) (err error) {
	return m.CreateAccessTokenSession(ctx, token, request)
}
