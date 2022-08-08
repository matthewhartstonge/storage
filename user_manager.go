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
	List(ctx context.Context, filter ListUsersRequest) ([]User, error)
	Create(ctx context.Context, user User) (User, error)
	Get(ctx context.Context, userID string) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	Update(ctx context.Context, userID string, user User) (User, error)
	Delete(ctx context.Context, userID string) error

	// Utility Functions
	Authenticate(ctx context.Context, username string, password string) (User, error)
	AuthenticateByID(ctx context.Context, userID string, password string) (User, error)
	AuthenticateByUsername(ctx context.Context, username string, password string) (User, error)
	GrantScopes(ctx context.Context, userID string, scopes []string) (User, error)
	RemoveScopes(ctx context.Context, userID string, scopes []string) (User, error)
}

// ListUsersRequest enables filtering stored User entities.
type ListUsersRequest struct {
	// AllowedTenantAccess filters users based on an Allowed Tenant Access.
	AllowedTenantAccess string `json:"allowedTenantAccess" xml:"allowedTenantAccess"`
	// AllowedPersonAccess filters users based on Allowed Person Access.
	AllowedPersonAccess string `json:"allowedPersonAccess" xml:"allowedPersonAccess"`
	// PersonID filters users based on a person's identifier.
	PersonID string `json:"personId" xml:"personId"`
	// PersonIDs filters users based on a list of person identifiers.
	// Setting PersonID will take precedence over this filter.
	PersonIDs []string `json:"personIds" xml:"personIds"`
	// Username filters users based on username.
	Username string `json:"username" xml:"username"`
	// ScopesUnion filters users that have at least one of the listed scopes.
	// ScopesUnion performs an OR operation.
	// If ScopesUnion is provided, a union operation will be performed as it
	// returns the wider selection.
	ScopesUnion []string `json:"scopesUnion" xml:"scopesUnion"`
	// ScopesIntersection filters users that have all the listed scopes.
	// ScopesIntersection performs an AND operation.
	ScopesIntersection []string `json:"scopesIntersection" xml:"scopesIntersection"`
	// FirstName filters users based on their First Name.
	FirstName string `json:"firstName" xml:"firstName"`
	// LastName filters users based on their Last Name.
	LastName string `json:"lastName" xml:"lastName"`
	// Disabled filters users to those with disabled accounts.
	Disabled bool `json:"disabled" xml:"disabled"`
}
