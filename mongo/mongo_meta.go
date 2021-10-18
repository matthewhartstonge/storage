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

	// IdxExpiry provides a mongo index for generating ttl based record
	// expiration indices.
	IdxExpiry = "idxExpiry"

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

// SessionToContext provides a way to push a mongo datastore session into the
// current context, which can then be passed on to other routes or functions.
func SessionToContext(ctx context.Context, session mongo.Session) context.Context {
	return mongo.NewSessionContext(ctx, session)
}

// ContextToSession provides a way to obtain a mongo session, if contained
// within the presented context.
func ContextToSession(ctx context.Context) (sess mongo.Session, ok bool) {
	if sess := mongo.SessionFromContext(ctx); sess != nil {
		return sess, true
	}

	return nil, false
}
