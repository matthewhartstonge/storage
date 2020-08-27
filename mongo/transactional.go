package mongo

import (
	// Standard Library Imports
	"context"

	// External Imports
	"github.com/ory/fosite"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func (s *Store) BeginTX(ctx context.Context) (context.Context, error) {
	session, ok := ContextToSession(ctx)
	if !ok {
		return ctx, errors.Wrap(
			fosite.ErrServerError,
			"transaction failed: no mongo session stored in context",
		)
	}

	// Define mongo transaction options...
	opts := options.Transaction().
		SetReadPreference(readpref.Primary()).
		SetReadConcern(readconcern.Majority()).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority())).
		SetMaxCommitTime(&s.timeout)

	err := session.StartTransaction(opts)
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

	err := txn.AbortTransaction(ctx)
	if err != nil {
		return errors.Wrap(fosite.ErrSerializationFailure, err.Error())
	}

	return nil
}
