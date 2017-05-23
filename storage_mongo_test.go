package storage_test

import (
	"github.com/MatthewHartstonge/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestConnectionURISingleHostNoCredentials ensures a correctly formed mongo connection URI is generated for connecting
// to an unsecured mongo host.
func TestConnectionURISingleHostNoCredentials(t *testing.T) {
	expected := "mongodb://127.0.0.1:27017/test"
	cfg := &storage.Config{
		Hostname:     "127.0.0.1",
		Hostnames:    nil,
		Port:         27017,
		DatabaseName: "test",

		// Credential Access
		Username: "",
		Password: "",

		// Replica Set
		Replset: "",
	}
	got := storage.ConnectionURI(cfg)
	assert.EqualValues(t, expected, got)
}

// TestConnectionURISingleHostCredentials ensures a correctly formed mongo connection URI is generated for connecting
// to a single mongo host with database access credentials.
func TestConnectionURISingleHostCredentials(t *testing.T) {
	expected := "mongodb://testuser:testuserpass@127.0.0.1:27017/test"
	cfg := &storage.Config{
		Hostname:     "127.0.0.1",
		Hostnames:    nil,
		Port:         27017,
		DatabaseName: "test",

		// Credential Access
		Username: "testuser",
		Password: "testuserpass",

		// Replica Set
		Replset: "",
	}
	got := storage.ConnectionURI(cfg)
	assert.EqualValues(t, expected, got)
}

// TestConnectionURIReplSetNoCredentials ensures a correctly formed mongo connection URI is generated for connecting
// to an unsecured mongo replica set.
func TestConnectionURIReplSetNoCredentials(t *testing.T) {
	expected := "mongodb://127.0.0.1:27017,127.0.1.1:27017,127.0.2.1:27017/test?replicaSet=sr0"
	cfg := &storage.Config{
		Hostname:     "",
		Hostnames:    []string{"127.0.0.1", "127.0.1.1", "127.0.2.1"},
		Port:         27017,
		DatabaseName: "test",

		// Credential Access
		Username: "",
		Password: "",

		// Replica Set
		Replset: "sr0",
	}
	got := storage.ConnectionURI(cfg)
	assert.EqualValues(t, expected, got)
}

// TestConnectionURIReplSetCredentials ensures a correctly formed mongo connection URI is generated for connecting
// to a mongo replica set with database access credentials.
func TestConnectionURIReplSetCredentials(t *testing.T) {
	expected := "mongodb://testuser:testuserpass@127.0.0.1:27017,127.0.1.1:27017,127.0.2.1:27017/test?replicaSet=sr0"
	cfg := &storage.Config{
		Hostname:     "",
		Hostnames:    []string{"127.0.0.1", "127.0.1.1", "127.0.2.1"},
		Port:         27017,
		DatabaseName: "test",

		// Credential Access
		Username: "testuser",
		Password: "testuserpass",

		// Replica Set
		Replset: "sr0",
	}
	got := storage.ConnectionURI(cfg)
	assert.EqualValues(t, expected, got)
}

// TestNewDatastoreErrorsWithBadConfig ensures the underlying lib used for Mongo creates an error
func TestNewDatastoreErrorsWithBadConfig(t *testing.T) {
	cfg := &storage.Config{
		Hostname:     "notevenanaddress",
		Port:         27666,
		DatabaseName: "lulz",
	}
	conn, err := storage.NewDatastore(cfg)
	assert.Nil(t, conn)
	assert.NotNil(t, err)
	assert.Error(t, err)
}
