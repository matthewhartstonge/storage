package mongo

import (
	// Standard Library Imports
	"context"

	// External Imports
	"github.com/globalsign/mgo"
)

const (
	// IdxCacheRequestId provides a mongo index based on request id.
	IdxCacheRequestId = "idxRequestId"

	// IdxCacheRequestSignature provides a mongo index based on token
	// signature.
	IdxCacheRequestSignature = "idxSignature"

	// IdxClientId provides a mongo index based on clientId
	IdxClientId = "IdxClientId"

	// IdxUserId provides a mongo index based on userId
	IdxUserId = "IdxUserId"

	// IdxUsername provides a mongo index based on username
	IdxUsername = "IdxUsername"
)

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

// ctxMgoKey is an unexported type for context keys defined for mgo in this
// package. This prevents collisions with keys defined in other packages.
type ctxMgoKey int

const (
	// mgoSessionKey is the key for *mgo.Session values in Contexts. It is
	// unexported; clients use datastore.MgoSessionToContext and
	// datastore.ContextToMgoSession instead of using this key directly.
	mgoSessionKey ctxMgoKey = iota
)

// MgoSessionToContext provides a way to push a Mgo datastore session into the
// current session, which can then be passed on to other routes or functions.
func MgoSessionToContext(ctx context.Context, session *mgo.Session) context.Context {
	return context.WithValue(ctx, mgoSessionKey, session)
}

func ContextToMgoSession(ctx context.Context) (sess *mgo.Session, ok bool) {
	sess, ok = ctx.Value(mgoSessionKey).(*mgo.Session)
	return
}
