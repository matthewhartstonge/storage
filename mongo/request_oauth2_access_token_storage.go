package mongo

import (
	// Standard Library Imports
	"context"

	// External Imports
	"github.com/ory/fosite"
	"github.com/sirupsen/logrus"

	// Internal Imports
	"github.com/matthewhartstonge/storage"
)

// CreateAccessTokenSession creates a new session for an Access Token
func (r *RequestManager) CreateAccessTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityAccessTokens,
		"method":     "CreateAccessTokenSession",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = r.DB.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "CreateAccessTokenSession",
	})
	defer span.Finish()

	// Store cached session
	cacheObj := storage.SessionCache{
		ID:        request.GetID(),
		Signature: signature,
	}
	_, err = r.Cache.Create(ctx, storage.EntityCacheAccessTokens, cacheObj)
	if err != nil {
		if err == storage.ErrResourceExists {
			log.WithError(err).Debug(logConflict)
			return err
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		return err
	}

	// Store session request
	_, err = r.Create(ctx, storage.EntityAccessTokens, toMongo(signature, request))
	if err != nil {
		if err == storage.ErrResourceExists {
			log.WithError(err).Debug(logConflict)
			return err
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		return err
	}

	return err
}

// GetAccessTokenSession returns a session if it can be found by signature
func (r *RequestManager) GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityAccessTokens,
		"method":     "GetAccessTokenSession",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = r.DB.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "GetAccessTokenSession",
	})
	defer span.Finish()

	// Get the stored request
	req, err := r.GetBySignature(ctx, storage.EntityAccessTokens, signature)
	if err != nil {
		if err == fosite.ErrNotFound {
			log.WithError(err).Debug(logNotFound)
			return nil, err
		}
		// Log to StdOut
		log.WithError(err).Error(logError)
		return nil, err
	}

	// Transform to a fosite.Request
	request, err = req.ToRequest(ctx, session, r.Clients)
	if err != nil {
		if err == fosite.ErrNotFound {
			log.WithError(err).Debug(logNotFound)
			return nil, err
		}
		// Log to StdOut
		log.WithError(err).Error(logError)
		return nil, err
	}

	return request, err
}

// DeleteAccessTokenSession removes an Access Token's session
func (r *RequestManager) DeleteAccessTokenSession(ctx context.Context, signature string) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityAccessTokens,
		"method":     "DeleteAccessTokenSession",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = r.DB.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "DeleteAccessTokenSession",
	})
	defer span.Finish()

	// Remove cached session
	err = r.Cache.DeleteByValue(ctx, storage.EntityCacheAccessTokens, signature)
	if err != nil {
		if err == fosite.ErrNotFound {
			log.WithError(err).Debug(logNotFound)
			return err
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		return err
	}

	// Remove session request
	err = r.DeleteBySignature(ctx, storage.EntityAccessTokens, signature)
	if err != nil {
		if err == fosite.ErrNotFound {
			log.WithError(err).Debug(logNotFound)
			return err
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		return err
	}

	return nil
}
