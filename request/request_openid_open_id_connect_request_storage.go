package request

import (
	"context"
	"github.com/ory/fosite"
)

/* These functions provide a concrete implementation of openid.OpenIDConnectRequestStorage */

// CreateOpenIDConnectSession creates an open id connect session
// for a given authorize code. This is relevant for explicit open id connect flow.
func (m *MongoManager) CreateOpenIDConnectSession(ctx context.Context, authorizeCode string, requester fosite.Requester) (err error) {
	return
}

// IsOpenIDConnectSession returns error
// - nil if a session was found,
// - ErrNoSessionFound if no session was found
// - or an arbitrary error if an error occurred.
func (m *MongoManager) GetOpenIDConnectSession(ctx context.Context, authorizeCode string, requester fosite.Requester) (req fosite.Requester, err error) {
	return
}

// DeleteOpenIDConnectSession removes an open id connect session from the store.
func (m *MongoManager) DeleteOpenIDConnectSession(ctx context.Context, authorizeCode string) (err error) {
	return
}
