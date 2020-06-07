package mongo

import (
	// Standard Library Imports
	"context"

	// External Imports
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// IdxCacheRequestID provides a mongo index based on request id.
	IdxCacheRequestID = "idxRequestId"

	// IdxCacheRequestSignature provides a mongo index based on token
	// signature.
	IdxCacheRequestSignature = "idxSignature"

	// IdxClientID provides a mongo index based on clientId
	IdxClientID = "idxClientId"

	// IdxUserID provides a mongo index based on userId
	IdxUserID = "idxUserId"

	// IdxUsername provides a mongo index based on username
	IdxUsername = "idxUsername"

	// IdxSessionID provides a mongo index based on Session
	IdxSessionID = "idxSessionId"

	// IdxSignatureID provides a mongo index based on Signature
	IdxSignatureID = "idxSignatureId"

	// IdxCompoundRequester provides a mongo compound index based on Client ID
	// and User ID for when filtering request records.
	IdxCompoundRequester = "idxCompoundRequester"
)

// ctxMongoKey is an unexported type for context keys defined for mgo in this
// package. This prevents collisions with keys defined in other packages.
type ctxMongoKey int

const (
	// mongoSessionKey is the key for *mgo.Session values in Contexts. It is
	// unexported; clients use datastore.SessionToContext and
	// datastore.ContextToSession instead of using this key directly.
	mongoSessionKey ctxMongoKey = iota
)

// SessionToContext provides a way to push a Mgo datastore session into the
// current session, which can then be passed on to other routes or functions.
func SessionToContext(ctx context.Context, session mongo.Session) context.Context {
	return context.WithValue(ctx, mongoSessionKey, session)
}

// ContextToSession provides a way to obtain a mgo session, if contained
// within the presented context.
func ContextToSession(ctx context.Context) (sess mongo.Session, ok bool) {
	sess, ok = ctx.Value(mongoSessionKey).(mongo.Session)
	return
}
