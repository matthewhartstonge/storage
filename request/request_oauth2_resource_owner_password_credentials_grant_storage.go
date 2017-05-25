package request

import (
	"context"
)

/* These functions provide a concrete implementation of fosite.ResourceOwnerPasswordCredentialsGrantStorage */
/* fosite.ResourceOwnerPasswordCredentialsGrantStorage also implements fosite.AccessTokenStorage */
/* fosite.ResourceOwnerPasswordCredentialsGrantStorage also implements fosite.RefreshTokenStorage */

func (m *MongoManager) Authenticate(ctx context.Context, name string, secret string) (err error) {
	return
}
