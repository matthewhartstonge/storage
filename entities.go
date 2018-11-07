package storage

const (
	// EntityOpenIDSessions provides the name of the entity to use in order to
	// create, read, update and delete OpenID Sessions.
	EntityOpenIDSessions = "openIDConnectSessions"

	// EntityAccessTokens provides the name of the entity to use in order to
	// create, read, update and delete Access Token sessions.
	EntityAccessTokens = "accessTokens"

	// EntityRefreshTokens provides the name of the entity to use in order to
	// create, read, update and delete Refresh Token sessions.
	EntityRefreshTokens = "refreshTokens"

	// EntityAuthorizationCodes provides the name of the entity to use in order
	// to create, read, update and delete Authorization Code sessions.
	EntityAuthorizationCodes = "authorizationCodes"

	// EntityPKCESessions provides the name of the entity to use in order to
	// create, read, update and delete Proof Key for Code Exchange sessions.
	EntityPKCESessions = "pkceSessions"

	// EntityClients provides the name of the entity to use in order to create,
	// read, update and delete Clients.
	EntityClients = "clients"

	// EntityUsers provides the name of the entity to use in order to create,
	// read, update and delete Users.
	EntityUsers = "users"

	// EntityCacheAccessTokens provides the name of the entity to use in order
	// to create, read, update and delete Cache Access Tokens.
	EntityCacheAccessTokens = "cacheAccessTokens"

	// EntityCacheRefreshTokens provides the name of the entity to use in
	// order to create, read, update and delete Cache Refresh Tokens.
	EntityCacheRefreshTokens = "cacheRefreshTokens"
)
