package request

import (
	// Standard Library Imports
	"context"

	// External Imports
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
)

// Manager provides a generic interface to clients in order to build a DataStore
type Manager interface {
	Storer
}

// Storer implements all fosite interfaces required to be a storage driver.
type Storer interface {
	//fosite.Requester

	// OAuth2 storage interfaces.
	oauth2.CoreStorage

	// OpenID storage interfaces.
	openid.OpenIDConnectRequestStorage

	// provides the storage implementation as specified in: fosite.handler.oauth2.TokenRevocationStorage
	RevokeRefreshToken(ctx context.Context, requestID string) error
	RevokeAccessToken(ctx context.Context, requestID string) error

	// Authenticate is required to implement the oauth2.ResourceOwnerPasswordCredentialsGrantStorage interface
	Authenticate(ctx context.Context, name string, secret string) error
}
