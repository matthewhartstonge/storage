package storage

import "context"

// Store brings all the interfaces together as a way to be composable into
// storage backend implementations
type Store struct {
	ClientManager
	DeniedJTIManager
	RequestManager
	UserManager
}

// Authenticate provides a top level pointer to UserManager to implement
// fosite.ResourceOwnerPasswordCredentialsGrantStorage at the top level.
//
// You can still access either the RequestManager API, or UserManager API by
// calling the methods on store direct depending on if you want the User
// resource returned as well via:
// - `store.RequestManager.Authenticate(ctx, username, secret) error`
// - `store.UserManager.Authenticate(ctx, username, secret) (User, error)`
func (s *Store) Authenticate(ctx context.Context, username string, secret string) error {
	_, err := s.UserManager.Authenticate(ctx, username, secret)
	return err
}

// AuthClientFunc enables developers to supply their own authentication
// function, to check old hashes that need to be upgraded for clients.
//
// For example, you may have passwords in MD5 and wanting them to be
// migrated to fosite's default hasher, bcrypt. Therefore, if you do a mass
// data migration, the function you supply would have to:
//   +- Shortcut logic if the hash string prefix matches what you expect from
// 		the new hash
//   - Get the current client record (return nil, false if not found)
//   - Authenticate the current DB secret against MD5
//   - Return the Client record and if the client authenticated.
//   - if true, the AuthenticateMigration function will upgrade the hash.
type AuthClientFunc func(ctx context.Context) (Client, bool)

// AuthUserFunc enables developers to supply their own authentication
// function, to check old hashes that need to be upgraded for users.
//
// See AuthClientFunc for example usage.
type AuthUserFunc func(ctx context.Context) (User, bool)

// AuthClientMigrator provides an interface to enable storage backends to
// implement functionality to upgrade hashes currently stored in the datastore.
type AuthClientMigrator interface {
	// Migrate is provided solely for the case where you want to migrate clients
	// and push in their old hash. This should perform an upsert, either
	// creating or overwriting the record with the newly provided record.
	// If Client.ID is passed in empty, a new ID will be generated for you.
	// Use with caution, be secure, don't be dumb.
	Migrate(ctx context.Context, migratedClient Client) (Client, error)

	// AuthenticateMigration enables developers to supply your own
	// authentication function, which in turn, if true, will migrate the secret
	// to the hasher implemented within fosite.
	AuthenticateMigration(ctx context.Context, currentAuth AuthClientFunc, clientID string, secret string) (Client, error)
}

// AuthUserMigrator provides an interface to enable storage backends to
// implement functionality to upgrade hashes currently stored in the datastore.
type AuthUserMigrator interface {
	// Migrate is provided solely for the case where you want to migrate users
	// and push in their old hash. This should perform an upsert, either
	// creating or overwriting the current record with the newly provided
	// record.
	// If User.ID is passed in empty, a new ID will be generated for you.
	// Use with caution, be secure, don't be dumb.
	Migrate(ctx context.Context, migratedUser User) (User, error)

	// AuthenticateMigration enables developers to supply your own
	// authentication function, which in turn, if true, will migrate the secret
	// to the hasher implemented within fosite.
	AuthenticateMigration(ctx context.Context, currentAuth AuthUserFunc, userID string, password string) (User, error)
}

// Configurer enables an implementer to configure required migrations, indexing
// and is called when the datastore connects.
type Configurer interface {
	// Configure configures the underlying database engine to match
	// requirements.
	// Configure will be called each time a service is started, so ensure this
	// function maintains idempotency.
	// The main use here is to apply creation of tables, collections, schemas
	// any needed migrations and configuration of indexes as required.
	Configure(ctx context.Context) error
}
