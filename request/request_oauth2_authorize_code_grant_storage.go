package request

import (
	"context"
	"github.com/ory/fosite"
)

/* These functions provide a concrete implementation of fosite.AuthorizeCodeGrantStorage */

// PersistAuthorizeCodeGrantSession creates an Authorise Code Grant session in mongo
func (m *MongoManager) PersistAuthorizeCodeGrantSession(ctx context.Context, authorizeCode, accessSignature, refreshSignature string, request fosite.Requester) error {
	if err := m.DeleteAuthorizeCodeSession(ctx, authorizeCode); err != nil {
		return err
	} else if err := m.CreateAccessTokenSession(ctx, accessSignature, request); err != nil {
		return err
	}

	if refreshSignature == "" {
		return nil
	}

	if err := m.CreateRefreshTokenSession(ctx, refreshSignature, request); err != nil {
		return err
	}

	return nil
}
