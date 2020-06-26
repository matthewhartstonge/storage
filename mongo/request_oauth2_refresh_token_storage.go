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

// CreateRefreshTokenSession implements fosite.RefreshTokenStorage.
func (r *RequestManager) CreateRefreshTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityRefreshTokens,
		"method":     "CreateRefreshTokenSession",
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
		Method:  "CreateRefreshTokenSession",
	})
	defer span.Finish()

	// Store session request
	_, err = r.Create(ctx, storage.EntityRefreshTokens, toMongo(signature, request))
	if err != nil {
		if err == storage.ErrResourceExists {
			log.WithError(err).Debug(logConflict)
			return err
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		return err
	}

	return nil
}

// GetRefreshTokenSession implements fosite.RefreshTokenStorage.
func (r *RequestManager) GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityRefreshTokens,
		"method":     "GetRefreshTokenSession",
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
		Method:  "GetRefreshTokenSession",
	})
	defer span.Finish()

	// Get the stored request
	req, err := r.GetBySignature(ctx, storage.EntityRefreshTokens, signature)
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

	return request, nil
}

// DeleteRefreshTokenSession implements fosite.RefreshTokenStorage.
func (r *RequestManager) DeleteRefreshTokenSession(ctx context.Context, signature string) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityRefreshTokens,
		"method":     "DeleteRefreshTokenSession",
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
		Method:  "DeleteRefreshTokenSession",
	})
	defer span.Finish()

	// Remove session request
	err = r.DeleteBySignature(ctx, storage.EntityRefreshTokens, signature)
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
