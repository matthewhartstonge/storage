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
	ClientStorer

	// Authenticate
	Authenticate(ctx context.Context, clientID string, secret []byte) (Client, error)
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
}

// ListClientsRequest enables listing and filtering client records.
type ListClientsRequest struct {
	// TenantID filters clients based on Tenant.
	TenantID string
	// RedirectURI filters clients based on redirectURI.
	RedirectURI string
	// GrantType filters clients based on GrantType.
	GrantType string
	// ResponseType filters clients based on ResponseType.
	ResponseType string
	// Scope filters clients based on Scope.
	Scope string
	// Contact filters clients based on Contact.
	Contact string
	// Public filters clients based on Public status.
	Public bool
	// Disabled filters clients based on denied access.
	Disabled bool
}
