package request

import (
	"context"
	"github.com/ory/fosite"
)

/* These functions provide a concrete implementation of fosite.ImplicitGrantStorage */

func (m *MongoManager) CreateImplicitAccessTokenSession(ctx context.Context, token string, request fosite.Requester) (err error) {
	return
}
