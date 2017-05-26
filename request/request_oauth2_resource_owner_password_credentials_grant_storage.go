package request

import (
	"context"
	"github.com/ory/fosite"
	"github.com/pkg/errors"
)

/* These functions provide a concrete implementation of fosite.ResourceOwnerPasswordCredentialsGrantStorage */
/* fosite.ResourceOwnerPasswordCredentialsGrantStorage also implements fosite.AccessTokenStorage */
/* fosite.ResourceOwnerPasswordCredentialsGrantStorage also implements fosite.RefreshTokenStorage */

func (m *MongoManager) Authenticate(ctx context.Context, username string, secret string) (err error) {
	user, err := m.Users.GetUserByUsername(username)
	if err != nil {
		return fosite.ErrNotFound
	}
	if err := m.Users.Hasher.Compare(user.GetHashedSecret(), []byte(secret)); err != nil {
		return errors.New("Invalid credentials")
	}
	return nil
}
