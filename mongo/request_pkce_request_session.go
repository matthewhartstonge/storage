package mongo

import (
	"context"

	"github.com/ory/fosite"
	"github.com/sirupsen/logrus"
)

func (r *requestMongoManager) GetPKCERequestSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": CollectionPKCESessions,
		"method":     "GetPKCERequestSession",
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
		Method:  "GetPKCERequestSession",
	})
	defer span.Finish()

	return context.TODO()
}

func (r *requestMongoManager) CreatePKCERequestSession(ctx context.Context, signature string, requester fosite.Requester) error {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": CollectionPKCESessions,
		"method":     "CreatePKCERequestSession",
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
		Method:  "CreatePKCERequestSession",
	})
	defer span.Finish()

	return context.TODO()
}

func (r *requestMongoManager) DeletePKCERequestSession(ctx context.Context, signature string) error {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": CollectionPKCESessions,
		"method":     "DeletePKCERequestSession",
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
		Method:  "DeletePKCERequestSession",
	})
	defer span.Finish()

	return context.TODO()
}
