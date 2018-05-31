package mongo

const (
	// CollectionOpenIDSessions provides the name of the mongo collection to use in order to create, read, update and delete OpenID Sessions
	CollectionOpenIDSessions = "openIDConnectSessions"

	// CollectionOpenIDSessions provides the name of the mongo collection to use in order to create, read, update and delete Access Tokens
	CollectionAccessTokens = "accessTokens"

	// CollectionOpenIDSessions provides the name of the mongo collection to use in order to create, read, update and delete Refresh Tokens
	CollectionRefreshTokens = "refreshTokens"

	// CollectionAuthorizationCodes provides the name of the mongo collection to use in order to create, read, update and delete Authorization Codes
	CollectionAuthorizationCodes = "authorizationCodes"

	// CollectionClients provides the name of the mongo collection to use in order to create, read, update and delete Clients
	CollectionClients = "clients"

	// CollectionUsers provides the name of the mongo collection to use in order to create, read, update and delete Users
	CollectionUsers = "users"

	// CollectionCacheAccessTokens provides the name of the mongo collection to use in order to create, read, update and delete Cache Access Tokens
	CollectionCacheAccessTokens = "cacheAccessTokens"

	// CollectionCacheRefreshTokens provides the name of the mongo collection to use in order to create, read, update and delete Cache Refresh Tokens
	CollectionCacheRefreshTokens = "cacheRefreshTokens"
)
