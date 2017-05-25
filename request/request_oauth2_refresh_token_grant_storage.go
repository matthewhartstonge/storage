package request

import (
	"context"
	"github.com/ory/fosite"
)

func (m *MongoManager) PersistRefreshTokenGrantSession(ctx context.Context, requestRefreshSignature, accessSignature, refreshSignature string, request fosite.Requester) (err error) {
	return
}
