package mongo

import (
	// Standard Library Imports
	"context"

	// External Imports
	"github.com/ory/fosite"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

	// TODO: reimplement when we can detect running mongodb version.
	// If attempting to use the transactions api mongo <4.0 mongo returns a
	// BSON serialization error.
	// For now, default to creating and using a unique session - sessions in
	// the mongo driver default to casual consistency, which should provide the
	// required atomicity.
	// Refer: https://jira.mongodb.org/projects/GODRIVER/issues/GODRIVER-1732

	// // Define mongo transaction options...
	// opts := options.Transaction().
	// 	SetReadPreference(readpref.Primary()).
	// 	SetReadConcern(readconcern.Majority()).
	// 	SetWriteConcern(writeconcern.New(writeconcern.WMajority())).
	// 	SetMaxCommitTime(&s.timeout)
	//
	// err = session.StartTransaction(opts)
	// if err != nil {
	// 	_ = session.AbortTransaction(ctx)
	// 	return ctx, err
	// }

	return SessionToContext(ctx, session), nil
}

func (s *Store) Commit(ctx context.Context) error {
	txn, ok := ContextToSession(ctx)
	if !ok {
		return errors.Wrap(
			fosite.ErrServerError,
			"commit failed: no transaction stored in context",
		)
	}
	defer txn.EndSession(ctx)

	// TODO: reimplement when we can detect running mongodb version.
	// err := txn.CommitTransaction(ctx)
	// if err != nil {
	// 	return errors.Wrap(fosite.ErrSerializationFailure, err.Error())
	// }

	return nil
}

func (s *Store) Rollback(ctx context.Context) error {
	txn, ok := ContextToSession(ctx)
	if !ok {
		return errors.Wrap(
			fosite.ErrServerError,
			"rollback failed: no transaction stored in context",
		)
	}
	defer txn.EndSession(ctx)

	// TODO: reimplement when we can detect running mongodb version.
	// err := txn.AbortTransaction(ctx)
	// if err != nil {
	// 	return errors.Wrap(fosite.ErrSerializationFailure, err.Error())
	// }

	return nil
}
