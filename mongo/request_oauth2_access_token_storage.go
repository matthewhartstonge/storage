package mongo

import (
	// Standard Library Imports
	"context"

	// External Imports
	"github.com/ory/fosite"
	"github.com/sirupsen/logrus"
)

// CreateAccessTokenSession creates a new session for an Access Token
func (r *requestMongoManager) CreateAccessTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": CollectionAccessTokens,
		"method":     "CreateAccessTokenSession",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = r.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "requestMongoManager",
		Method:  "CreateAccessTokenSession",
	})
	defer span.Finish()

	return context.TODO()
}

// GetAccessTokenSession returns a session if it can be found by signature
func (r *requestMongoManager) GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": CollectionAccessTokens,
		"method":     "GetAccessTokenSession",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = r.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "requestMongoManager",
		Method:  "GetAccessTokenSession",
	})
	defer span.Finish()

	return context.TODO()
}

// DeleteAccessTokenSession removes an Access Token's session
func (r *requestMongoManager) DeleteAccessTokenSession(ctx context.Context, signature string) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": CollectionAccessTokens,
		"method":     "DeleteAccessTokenSession",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = r.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "requestMongoManager",
		Method:  "DeleteAccessTokenSession",
	})
	defer span.Finish()

	return context.TODO()
}
