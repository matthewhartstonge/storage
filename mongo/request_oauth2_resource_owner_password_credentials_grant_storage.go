package mongo

import (
	// Standard Library Imports
	"context"
	"errors"

	// External Imports
	"github.com/ory/fosite"
	"github.com/sirupsen/logrus"

	// Internal Imports
	"github.com/matthewhartstonge/storage"
)

// Provides a concrete implementation of oauth2.ResourceOwnerPasswordCredentialsGrantStorage
// oauth2.ResourceOwnerPasswordCredentialsGrantStorage also implements
// oauth2.AccessTokenStorage and oauth2.RefreshTokenStorage

// Authenticate confirms whether the specified password matches the stored
// hashed password within a User resource, found by username.
func (r *RequestManager) Authenticate(ctx context.Context, username string, secret string) (subject string, err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityUsers,
		"method":     "Authenticate",
	})

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "Authenticate",
	})
	defer span.Finish()

	u, err := r.Users.Authenticate(ctx, username, secret)
	if err != nil {
		if errors.Is(err, fosite.ErrNotFound) {
			log.WithError(err).Debug(logNotFound)
			return "", err
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		return "", err
	}

	return u.ID, nil
}
