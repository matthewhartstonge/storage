package request

import (
	"context"
	"github.com/ory/fosite"
)

/* These functions provide a concrete implementation of fosite.AuthorizeCodeGrantStorage */

func (m *MongoManager) PersistAuthorizeCodeGrantSession(ctx context.Context, authorizeCode, accessSignature, refreshSignature string, request fosite.Requester) error {
	return nil
}
