package storage

import (
	// Standard Library Imports
	"context"
)

// DeniedJTIManager provides a generic interface to clients in order to build a
// Datastore backend.
type DeniedJTIManager interface {
	Configurer
	DeniedJTIStorer
}

// DeniedJTIStorer enables storing denied JWT Tokens, by ID.
type DeniedJTIStorer interface {
	// Standard CRUD Storage API
	Create(ctx context.Context, deniedJti DeniedJTI) (DeniedJTI, error)
	Get(ctx context.Context, jti string) (DeniedJTI, error)
	Delete(ctx context.Context, jti string) error

	// DeleteBefore removes all denied JTIs before the given unix time.
	DeleteBefore(ctx context.Context, expBefore int64) error
}
