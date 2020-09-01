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

// CreatePKCERequestSession implements fosite.PKCERequestStorage.
func (r *RequestManager) CreatePKCERequestSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityPKCESessions,
		"method":     "CreatePKCERequestSession",
	})

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "CreatePKCERequestSession",
	})
	defer span.Finish()

	// Store session request
	_, err = r.Create(ctx, storage.EntityPKCESessions, toMongo(signature, request))
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

// GetPKCERequestSession implements fosite.PKCERequestStorage.
func (r *RequestManager) GetPKCERequestSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityPKCESessions,
		"method":     "GetPKCERequestSession",
	})

	if r.DB.HasSessions {
		// Copy a new DB session if none specified
		_, ok := ContextToSession(ctx)
		if !ok {
			var closeSession func()
			ctx, closeSession, err = newSession(ctx, r.DB)
			if err != nil {
				log.WithError(err).Debug("error starting session")
				return nil, err
			}
			defer closeSession()
		}
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "GetPKCERequestSession",
	})
	defer span.Finish()

	// Get the stored request
	req, err := r.GetBySignature(ctx, storage.EntityPKCESessions, signature)
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

// DeletePKCERequestSession implements fosite.PKCERequestStorage.
func (r *RequestManager) DeletePKCERequestSession(ctx context.Context, signature string) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityPKCESessions,
		"method":     "DeletePKCERequestSession",
	})

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "DeletePKCERequestSession",
	})
	defer span.Finish()

	// Remove session request
	err = r.DeleteBySignature(ctx, storage.EntityPKCESessions, signature)
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
