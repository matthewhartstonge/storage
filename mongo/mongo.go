package mongo

const (
	// CollectionOpenIDSessions provides the name of the mongo collection to use in order to create, read, update and delete OpenID Sessions
	CollectionOpenIDSessions = "OpenIDConnectSessions"

	// CollectionOpenIDSessions provides the name of the mongo collection to use in order to create, read, update and delete Access Tokens
	CollectionAccessTokens = "AccessTokens"

	// CollectionOpenIDSessions provides the name of the mongo collection to use in order to create, read, update and delete Refresh Tokens
	CollectionRefreshTokens = "RefreshTokens"

	// CollectionAuthorizationCodes provides the name of the mongo collection to use in order to create, read, update and delete Authorization Codes
	CollectionAuthorizationCodes = "AuthorizationCodes"

	// CollectionClients provides the name of the mongo collection to use in order to create, read, update and delete Clients
	CollectionClients = "Clients"

	// CollectionUsers provides the name of the mongo collection to use in order to create, read, update and delete Users
	CollectionUsers = "Users"

	// CollectionCacheAccessTokens provides the name of the mongo collection to use in order to create, read, update and delete Cache Access Tokens
	CollectionCacheAccessTokens = "CacheAccessTokens"

	// CollectionCacheRefreshTokens provides the name of the mongo collection to use in order to create, read, update and delete Cache Refresh Tokens
	CollectionCacheRefreshTokens = "CacheRefreshTokens"
)
