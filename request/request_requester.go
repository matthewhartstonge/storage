package request

import (
	"github.com/ory/fosite"
	"net/url"
	"time"
)

/* These functions implement fosite.Requester */

// GetID returns a unique identifier.
func (m *MongoManager) GetID() string { return "" }

// GetRequestedAt returns the time the request was created.
func (m *MongoManager) GetRequestedAt() (requestedAt time.Time) { return }

// GetClient returns the requests client.
func (m *MongoManager) GetClient() (client fosite.Client) { return }

// GetRequestedScopes returns the request's scopes.
func (m *MongoManager) GetRequestedScopes() (scopes fosite.Arguments) { return }

// SetRequestedScopes sets the request's scopes.
func (m *MongoManager) SetRequestedScopes(scopes fosite.Arguments) { return }

// AppendRequestedScope appends a scope to the request.
func (m *MongoManager) AppendRequestedScope(scope string) { return }

// GetGrantScopes returns all granted scopes.
func (m *MongoManager) GetGrantedScopes() (grantedScopes fosite.Arguments) { return }

// GrantScope marks a request's scope as granted.
func (m *MongoManager) GrantScope(scope string) {}

// GetSession returns a pointer to the request's session or nil if none is set.
func (m *MongoManager) GetSession() (session fosite.Session) { return }

// GetSession sets the request's session pointer.
func (m *MongoManager) SetSession(session fosite.Session) {}

// GetRequestForm returns the request's form input.
func (m *MongoManager) GetRequestForm() url.Values { return url.Values{} }

func (m *MongoManager) Merge(requester fosite.Requester) {}
