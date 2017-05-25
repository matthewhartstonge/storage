package request

import (
	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
)

// Manager provides a generic interface to clients in order to build a DataStore
type Manager interface {
	Storer
}

// Storer conforms to fosite.Requester and provides methods
type Storer interface {
	fosite.Requester

	// OAuth2 Required Storage interfaces
	oauth2.AuthorizeCodeGrantStorage
	oauth2.ClientCredentialsGrantStorage
	oauth2.ImplicitGrantStorage
	oauth2.RefreshTokenGrantStorage
	oauth2.ResourceOwnerPasswordCredentialsGrantStorage

	// OpenID Required Storage Interfaces
	openid.OpenIDConnectRequestStorage

	// Enable revoking of tokens
	// see: https://github.com/ory/hydra/blob/master/pkg/fosite_storer.go
	//RevokeRefreshToken(ctx context.Context, requestID string) error
	//RevokeAccessToken(ctx context.Context, requestID string) error
}
