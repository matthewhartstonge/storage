package request

import (
	"context"
	"github.com/ory/fosite"
)

/* These functions provide a concrete implementation of fosite.handler.oauth2.PersistRefreshTokenGrantSession */

// PersistRefreshTokenGrantSession stores a refresh token grant session in mongo
func (m *MongoManager) PersistRefreshTokenGrantSession(ctx context.Context, requestRefreshSignature, accessSignature, refreshSignature string, request fosite.Requester) (err error) {
	if err := m.DeleteRefreshTokenSession(ctx, requestRefreshSignature); err != nil {
		return err
	} else if err := m.CreateAccessTokenSession(ctx, accessSignature, request); err != nil {
		return err
	} else if err := m.CreateRefreshTokenSession(ctx, refreshSignature, request); err != nil {
		return err
	}

	return nil
}
