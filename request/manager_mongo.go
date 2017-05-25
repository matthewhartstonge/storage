package request

import (
	"gopkg.in/mgo.v2"
)

const (
	mongoCollectionOpenIDSessions    = "OpenIDConnectSessions"
	mongoCollectionAccessTokens      = "AccessTokens"
	mongoCollectionRefreshTokens     = "RefreshTokens"
	mongoCollectionAuthorizationCode = "AuthorizationCode"
)

// MongoManager manages the main Mongo Session for a Request.
type MongoManager struct {
	// DB is the Mongo connection that holds the base session that can be copied and closed.
	DB *mgo.Database

	// TODO: Add AES cipher for Token Encryption?
}
