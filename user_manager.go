package storage

import "context"

// UserManager provides a generic interface to users in order to build a DataStore
type UserManager interface {
	Configurer
	UserStorer
	AuthUserMigrator
}

// UserStorer provides a definition of specific methods that are required to store a User in a data store.
type UserStorer interface {
	List(ctx context.Context, filter ListClientsRequest) ([]User, error)
	Create(ctx context.Context, user User) (User, error)
	Get(ctx context.Context, userID string) (User, error)
	Update(ctx context.Context, userID string, user User) (User, error)
	Delete(ctx context.Context, userID string) error

	// Utility Functions
	Authenticate(ctx context.Context, username string, secret []byte) (User, error)
	GrantScopes(ctx context.Context, userID string, scopes []string) error
	RemoveScopes(ctx context.Context, userID string, scopes []string) error
	AuthenticateByID(ctx context.Context, userID string, secret []byte) (User, error)
	AuthenticateByUsername(ctx context.Context, username string, secret []byte) (User, error)
}

type ListUsersRequest struct {
	// TenantID filters users based on Tenant Access.
	TenantID string `json:"tenantId" xml:"tenantId"`
	// PersonID filters users based on Allowed Person Access.
	PersonID string `json:""`
	// PersonID filters users based on Person Access.
	PID string
	// Username filters users based on username.
	Username string
	// Scopes filters users based on scopes users must have.
	// Scopes performs an AND operation. To obtain OR, do multiple requests with a single scope.
	Scopes []string
	// FirstName filters users based on their First Name.
	FirstName string
	// LastName filters users based on their Last Name.
	LastName string
	// Disabled filters users to those with disabled accounts.
	Disabled bool
}
