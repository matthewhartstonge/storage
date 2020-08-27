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
	// mongoTxnKey provides a context key for storing transaction based sessions
	// in contexts.
	mongoTxnKey
)

// SessionToContext provides a way to push a mongo datastore session into the
// current context, which can then be passed on to other routes or functions.
func SessionToContext(ctx context.Context, session mongo.Session) context.Context {
	return context.WithValue(ctx, mongoSessionKey, session)
}

// ContextToSession provides a way to obtain a mongo session, if contained
// within the presented context.
func ContextToSession(ctx context.Context) (sess mongo.Session, ok bool) {
	if tx, ok := ContextToTransaction(ctx); ok {
		// Always return the transaction session if available.
		return tx, ok
	}

	return ctxToSess(ctx, mongoSessionKey)
}

// TransactionToContext provides a way to push a mongo transaction-based
// session into the provided context, which can then be passed on to other
// routes or functions.
func TransactionToContext(ctx context.Context, session mongo.Session) context.Context {
	return context.WithValue(ctx, mongoTxnKey, session)
}

// ContextToTransaction provides a way to obtain a mongo transaction-based
// session, if contained within the presented context.
func ContextToTransaction(ctx context.Context) (sess mongo.Session, ok bool) {
	return ctxToSess(ctx, mongoTxnKey)
}

// ctxToSess provides a wrapper around extracting a mongo session from a
// context.
func ctxToSess(ctx context.Context, key ctxMongoKey) (sess mongo.Session, ok bool) {
	sess, ok = ctx.Value(key).(mongo.Session)
	return
}
