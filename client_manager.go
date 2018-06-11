package storage

import (
	// Standard Library Imports
	"context"

	// External Imports
	"github.com/ory/fosite"
)

// ClientManager provides a generic interface to clients in order to build a
// Datastore backend.
type ClientManager interface {
	Configurer
	ClientStorer
	AuthClientMigrator
}

// ClientStorer conforms to fosite.Storage and provides methods
type ClientStorer interface {
	// fosite.Storage provides get client.
	fosite.Storage

	List(ctx context.Context, filter ListClientsRequest) ([]Client, error)
	Create(ctx context.Context, client Client) (Client, error)
	Get(ctx context.Context, clientID string) (Client, error)
	Update(ctx context.Context, clientID string, client Client) (Client, error)
	Delete(ctx context.Context, clientID string) error

	// Utility Functions
	Authenticate(ctx context.Context, clientID string, secret string) (Client, error)
	GrantScopes(ctx context.Context, clientID string, scopes []string) (Client, error)
	RemoveScopes(ctx context.Context, clientID string, scopes []string) (Client, error)
}

// ListClientsRequest enables listing and filtering client records.
type ListClientsRequest struct {
	// AllowedTenantAccess filters clients based on an Allowed Tenant Access.
	AllowedTenantAccess string `json:"allowedTenantAccess" xml:"allowedTenantAccess"`
	// RedirectURI filters clients based on redirectURI.
	RedirectURI string `json:"redirectURI" xml:"redirectURI"`
	// GrantType filters clients based on GrantType.
	GrantType string `json:"grantType" xml:"grantType"`
	// ResponseType filters clients based on ResponseType.
	ResponseType string `json:"responseType" xml:"responseType"`
	// ScopesIntersection filters clients that have all of the listed scopes.
	// ScopesIntersection performs an AND operation.
	// If ScopesUnion is provided, a union operation will be performed as it
	// returns the wider selection.
	ScopesIntersection []string `json:"scopesIntersection" xml:"scopesIntersection"`
	// ScopesUnion filters users that have at least one of of the listed scopes.
	// ScopesUnion performs an OR operation.
	ScopesUnion []string `json:"scopesUnion" xml:"scopesUnion"`
	// Contact filters clients based on Contact.
	Contact string `json:"contact" xml:"contact"`
	// Public filters clients based on Public status.
	Public bool `json:"public" xml:"public"`
	// Disabled filters clients based on denied access.
	Disabled bool `json:"disabled" xml:"disabled"`
}
