package mongo

import (
	// Standard Library Imports
	"context"

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
	session, err := s.DB.Client().StartSession()
	if err != nil {
		fields := logrus.Fields{
			"package": "mongo",
			"method":  "BeginTX",
		}
		logger.WithError(err).WithFields(fields).Error("error starting mongo transaction")
		return ctx, err
	}

	// Define mongo transaction options...
	opts := options.Transaction().
		SetReadPreference(readpref.Primary()).
		SetReadConcern(readconcern.Majority()).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority())).
		SetMaxCommitTime(&s.timeout)

	err = session.StartTransaction(opts)
	if err != nil {
		_ = session.AbortTransaction(ctx)
		return ctx, err
	}

	return TransactionToContext(ctx, session), nil
}

func (s *Store) Commit(ctx context.Context) error {
	txn, ok := ContextToTransaction(ctx)
	if !ok {
		return errors.Wrap(
			fosite.ErrServerError,
			"commit failed: no transaction stored in context",
		)
	}
	defer txn.EndSession(ctx)

	err := txn.CommitTransaction(ctx)
	if err != nil {
		return errors.Wrap(fosite.ErrSerializationFailure, err.Error())
	}

	return nil
}

func (s *Store) Rollback(ctx context.Context) error {
	txn, ok := ContextToTransaction(ctx)
	if !ok {
		return errors.Wrap(
			fosite.ErrServerError,
			"rollback failed: no transaction stored in context",
		)
	}
	defer txn.EndSession(ctx)

	err := txn.AbortTransaction(ctx)
	if err != nil {
		return errors.Wrap(fosite.ErrSerializationFailure, err.Error())
	}

	return nil
}
