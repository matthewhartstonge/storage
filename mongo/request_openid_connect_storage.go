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

// CreateOpenIDConnectSession creates an open id connect session resource for a
// given authorize code. This is relevant for explicit open id connect flow.
func (r *RequestManager) CreateOpenIDConnectSession(ctx context.Context, authorizeCode string, request fosite.Requester) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityOpenIDSessions,
		"method":     "CreateOpenIDConnectSession",
	})

	// Copy a new DB session if none specified
	_, ok := ContextToSession(ctx)
	if !ok {
		var closer func()
		ctx, _, closer, err = newSession(ctx, r.DB)
		if err != nil {
			log.WithError(err).Debug("error starting session")
			return err
		}
		defer closer()
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "CreateOpenIDConnectSession",
	})
	defer span.Finish()

	// Store session request
	_, err = r.Create(ctx, storage.EntityOpenIDSessions, toMongo(authorizeCode, request))
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

// GetOpenIDConnectSession gets a session resource based off the Authorize Code
// and returns a fosite.Requester, or an error.
func (r *RequestManager) GetOpenIDConnectSession(ctx context.Context, authorizeCode string, requester fosite.Requester) (request fosite.Requester, err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityOpenIDSessions,
		"method":     "GetOpenIDConnectSession",
	})

	// Copy a new DB session if none specified
	_, ok := ContextToSession(ctx)
	if !ok {
		var closer func()
		ctx, _, closer, err = newSession(ctx, r.DB)
		if err != nil {
			log.WithError(err).Debug("error starting session")
			return nil, err
		}
		defer closer()
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "GetOpenIDConnectSession",
	})
	defer span.Finish()

	// Get the stored request
	req, err := r.GetBySignature(ctx, storage.EntityOpenIDSessions, authorizeCode)
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
	session := requester.GetSession()
	if session == nil {
		return nil, fosite.ErrNotFound
	}

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

// DeleteOpenIDConnectSession removes an open id connect session from mongo.
func (r *RequestManager) DeleteOpenIDConnectSession(ctx context.Context, authorizeCode string) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityOpenIDSessions,
		"method":     "DeleteOpenIDConnectSession",
	})

	// Copy a new DB session if none specified
	_, ok := ContextToSession(ctx)
	if !ok {
		var closer func()
		ctx, _, closer, err = newSession(ctx, r.DB)
		if err != nil {
			log.WithError(err).Debug("error starting session")
			return err
		}
		defer closer()
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "DeleteOpenIDConnectSession",
	})
	defer span.Finish()

	// Remove session request
	err = r.DeleteBySignature(ctx, storage.EntityOpenIDSessions, authorizeCode)
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
