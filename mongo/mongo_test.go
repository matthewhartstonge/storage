package mongo_test

import (
	// Standard Library Imports
	"context"
	"os"
	"testing"

	// External Imports
	"github.com/google/go-cmp/cmp"

	// Internal Imports
	"github.com/matthewhartstonge/storage"
	"github.com/matthewhartstonge/storage/mongo"
)

func TestMain(m *testing.M) {
	// If needed, enable logging when debugging for tests
	// mongo.SetLogger(logrus.New())
	// mongo.SetDebug(true)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func AssertError(t *testing.T, got interface{}, want interface{}, msg string) {
	t.Helper()
	t.Errorf("Error: %s\n	 got: %#+v\n	want: %#+v", msg, got, want)
}

func AssertFatal(t *testing.T, got interface{}, want interface{}, msg string) {
	t.Helper()
	t.Fatalf("Fatal: %s\n	 got: %#+v\n	want: %#+v", msg, got, want)
}

func setup(t *testing.T) (*mongo.Store, context.Context, func()) {
	// Build our default mongo storage layer
	cfg := mongo.DefaultConfig()
	cfg.DatabaseName = "fositeStorageTest"
	store, err := mongo.New(cfg, nil)
	if err != nil {
		AssertFatal(t, err, nil, "mongo connection error")
	}

	// Build a context with a mongo session ready to use for testing
	ctx := context.Background()
	var closeSession func()
	ctx, closeSession, err = store.NewSession(ctx)
	if err != nil {
		AssertFatal(t, err, nil, "error getting mongo session")
	}

	return store, ctx, func() {
		// Drop the database.
		err := store.DB.Drop(ctx)
		if err != nil {
			t.Errorf("error dropping database on cleanup: %s", err)
		}

		// Close the inner (test) session if it exists.
		closeSession()

		// Close the database connection.
		store.Close()
	}
}

type TestCustomData struct {
	ID      string
	Contact TestContactData
}

type TestContactData struct {
	Name  string
	Goals []string
}

func expectedCustomData() TestCustomData {
	return TestCustomData{
		ID: "001",
		Contact: TestContactData{
			Name: "Data, Custom Data",
			Goals: []string{
				"Store data successfully",
				"Make sure the data de/serializes as expected",
				"Update said data successfully",
			},
		},
	}
}

func AssertClientCustomData(t *testing.T, gotEntity, wantEntity storage.Client, gotData, wantData any) {
	t.Helper()
	if err := gotEntity.Data.Unmarshal(&gotData); err != nil {
		AssertError(t, err, nil, "failed to unmarshal custom data")
	}

	if diff := cmp.Diff(wantData, gotData); diff != "" {
		t.Errorf("Custom data mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(wantEntity, gotEntity); diff != "" {
		t.Errorf("client data mismatch (-want +got):\n%s", diff)
	}
}

func AssertUserCustomData(t *testing.T, gotEntity, wantEntity storage.User, gotData, wantData any) {
	t.Helper()
	if err := gotEntity.Data.Unmarshal(gotData); err != nil {
		AssertError(t, err, nil, "failed to unmarshal custom data")
	}

	if diff := cmp.Diff(wantData, gotData); diff != "" {
		t.Errorf("Custom data mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(wantEntity, gotEntity); diff != "" {
		t.Errorf("User data mismatch (-want +got):\n%s", diff)
	}
}
