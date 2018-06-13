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

// CreateAuthorizeCodeSession stores the authorization request for a given
// authorization code.
func (r *requestMongoManager) CreateAuthorizeCodeSession(ctx context.Context, code string, request fosite.Requester) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityAuthorizationCodes,
		"method":     "CreateAuthorizeCodeSession",
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
		Method:  "CreateAuthorizeCodeSession",
	})
	defer span.Finish()

	// Store session request
	_, err = r.Create(ctx, storage.EntityAuthorizationCodes, toMongo(code, request))
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

// GetAuthorizeCodeSession hydrates the session based on the given code and
// returns the authorization request.
func (r *requestMongoManager) GetAuthorizeCodeSession(ctx context.Context, code string, session fosite.Session) (request fosite.Requester, err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityAuthorizationCodes,
		"method":     "GetAuthorizeCodeSession",
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
		Method:  "GetAuthorizeCodeSession",
	})
	defer span.Finish()

	// Get the stored request
	req, err := r.GetBySignature(ctx, storage.EntityAuthorizationCodes, code)
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

	if !req.Active {
		// If the authorization code has been invalidated with
		// `InvalidateAuthorizeCodeSession`, this method should return the
		// ErrInvalidatedAuthorizeCode error.
		// Make sure to also return the fosite.Requester value when returning
		// the ErrInvalidatedAuthorizeCode error!
		return request, fosite.ErrInvalidatedAuthorizeCode
	}

	return request, err
}

// InvalidateAuthorizeCodeSession is called when an authorize code is being
// used. The state of the authorization code should be set to invalid and
// consecutive requests to GetAuthorizeCodeSession should return the
// ErrInvalidatedAuthorizeCode error.
func (r *requestMongoManager) InvalidateAuthorizeCodeSession(ctx context.Context, code string) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityAuthorizationCodes,
		"method":     "InvalidateAuthorizeCodeSession",
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
		Method:  "InvalidateAuthorizeCodeSession",
	})
	defer span.Finish()

	// Get the stored request
	req, err := r.GetBySignature(ctx, storage.EntityAuthorizationCodes, code)
	if err != nil {
		if err == fosite.ErrNotFound {
			log.WithError(err).Debug(logNotFound)
			return err
		}
		// Log to StdOut
		log.WithError(err).Error(logError)
		return err
	}

	// InvalidateAuthorizeCodeSession
	req.Active = false

	// Push the update back
	req, err = r.Update(ctx, storage.EntityAuthorizationCodes, req.ID, req)
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
