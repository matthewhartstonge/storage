package mongo

import (
	// Standard Library Imports
	"context"

	// External Imports
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// IdxClientID provides a mongo index based on clientId
	IdxClientID = "idxClientId"

	// IdxExpires provides a mongo index based on expires
	IdxExpires = "idxExpires"

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

// ctxMongoKey is an unexported type to enable passing mongo information via a
// context in this package. This prevents collisions with keys defined in other
// packages.
type ctxMongoKey int

const (
	// mongoSessionKey is the key for mongo.Session values in Contexts. It is
	// unexported; clients should use datastore.SessionToContext and
	// datastore.ContextToSession instead of attempting to use this key.
	mongoSessionKey ctxMongoKey = iota
)

// SessionToContext provides a way to push a mongo datastore session into the
// current session, which can then be passed on to other routes or functions.
func SessionToContext(ctx context.Context, session mongo.Session) context.Context {
	return context.WithValue(ctx, mongoSessionKey, session)
}

// ContextToSession provides a way to obtain a mongo session, if contained
// within the presented context.
func ContextToSession(ctx context.Context) (sess mongo.Session, ok bool) {
	sess, ok = ctx.Value(mongoSessionKey).(mongo.Session)
	return
}
