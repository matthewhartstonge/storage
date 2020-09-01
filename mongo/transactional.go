package mongo

import (
	// Standard Library Imports
	"context"
	"fmt"

	// External Imports
	"github.com/ory/fosite"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func (s *Store) BeginTX(ctx context.Context) (context context.Context, err error) {
	log := logger.WithFields(logrus.Fields{
		"package": "mongo",
		"method":  "BeginTX",
	})

	if !s.DB.HasSessions {
		return ctx, nil
	}

	session, err := s.DB.Client().StartSession()
	if err != nil {
		log.WithError(err).Error("error starting mongo session")
		return ctx, err
	}

	if s.DB.HasTransactions {
		// Define mongo transaction options...
		opts := options.Transaction().
			SetReadPreference(readpref.Primary()).
			SetReadConcern(readconcern.Majority()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())).
			SetMaxCommitTime(&s.timeout)

		err = session.StartTransaction(opts)
		if err != nil {
			log.WithError(err).Error("error starting mongo transaction")

			txErr := session.AbortTransaction(ctx)
			if txErr != nil {
				log.WithError(txErr).Warn("error aborting mongo transaction")
			}

			return ctx, errors.Wrap(
				fosite.ErrSerializationFailure,
				fmt.Sprintf("starting transaction failed: %s\n", err.Error()),
			)
		}
	}

	return SessionToContext(ctx, session), nil
}

func (s *Store) Commit(ctx context.Context) error {
	log := logger.WithFields(logrus.Fields{
		"package": "mongo",
		"method":  "Commit",
	})

	if !s.DB.HasSessions {
		return nil
	}

	txn, ok := ContextToSession(ctx)
	if !ok {
		return errors.Wrap(
			fosite.ErrServerError,
			"commit failed: no session stored in context",
		)
	}
	defer txn.EndSession(ctx)

	if s.DB.HasTransactions {
		err := txn.CommitTransaction(ctx)
		if err != nil {
			log.WithError(err).Error("error committing mongo transaction")

			txErr := txn.AbortTransaction(ctx)
			if txErr != nil {
				log.WithError(txErr).Warn("error aborting mongo transaction")
			}

			return errors.Wrap(
				fosite.ErrSerializationFailure,
				fmt.Sprintf("commit failed: %s\n", err.Error()),
			)
		}
	}

	return nil
}

func (s *Store) Rollback(ctx context.Context) error {
	if !s.DB.HasSessions {
		return nil
	}

	txn, ok := ContextToSession(ctx)
	if !ok {
		return errors.Wrap(
			fosite.ErrServerError,
			"rollback failed: no session stored in context",
		)
	}
	defer txn.EndSession(ctx)

	if s.DB.HasTransactions {
		if err := txn.AbortTransaction(ctx); err != nil {
			return errors.Wrap(
				fosite.ErrSerializationFailure,
				fmt.Sprintf("rollback failed: %s\n", err.Error()),
			)
		}
	}

	return nil
}
