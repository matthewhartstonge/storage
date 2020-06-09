package mongo_test

import (
	// Standard Library Imports
	"context"
	"fmt"
	"os"
	"testing"

	// Public Imports
	"github.com/matthewhartstonge/storage/mongo"
)

func TestMain(m *testing.M) {
	// If needed, enable logging when debugging for tests
	//mongo.SetLogger(logrus.New())
	//mongo.SetDebug(true)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func AssertError(t *testing.T, got interface{}, want interface{}, msg string) {
	t.Errorf(fmt.Sprintf("Error: %s\n	 got: %#+v\n	want: %#+v", msg, got, want))
}

func AssertFatal(t *testing.T, got interface{}, want interface{}, msg string) {
	t.Fatalf(fmt.Sprintf("Fatal: %s\n	 got: %#+v\n	want: %#+v", msg, got, want))
}

func setup(t *testing.T) (*mongo.Store, context.Context, func()) {
	// Build our default mongo storage layer
	store, err := mongo.NewDefaultStore()
	if err != nil {
		AssertFatal(t, err, nil, "mongo connection error")
	}

	// Build a context with a mongo session ready to use for testing
	ctx, _, closeSession, err := store.NewSession(nil)
	if err != nil {
		AssertFatal(t, err, nil, "error getting mongo session")
	}

	return store, ctx, func() {
		// Close the inner (test) session.
		closeSession()

		// Drop the database.
		err := store.DB.Drop(ctx)
		if err != nil {
			t.Errorf("error dropping database on cleanup: %s", err)
		}

		// Close the database connection.
		store.Close()
	}
}
