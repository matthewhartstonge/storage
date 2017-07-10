package storage_test

import (
	"github.com/MatthewHartstonge/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestDefaultConfig_IsDefault ensures the Default config is one that allows connection to an unauthenticated locally
// hosted mongo instance.
func TestDefaultConfig_IsDefault(t *testing.T) {
	expected := &storage.Config{
		Hostname:     "127.0.0.1",
		Port:         27017,
		DatabaseName: "OAuth2",
	}
	got := storage.DefaultConfig()
	assert.NotNil(t, got)
	assert.ObjectsAreEqualValues(expected, got)
}

// TestConnectionURI_SingleHostNoCredentials ensures a correctly formed mongo connection URI is generated for connecting
// to an unsecured mongo host.
func TestConnectionURI_SingleHostNoCredentials(t *testing.T) {
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

// TestConnectionURI_SingleHostCredentials ensures a correctly formed mongo connection URI is generated for connecting
// to a single mongo host with database access credentials.
func TestConnectionURI_SingleHostCredentials(t *testing.T) {
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

// TestConnectionURI_Hostnames ensures a correctly formed mongo connection URI is generated when a single hostname is
// passed as hostnames.
func TestConnectionURI_MultiHostnames(t *testing.T) {
	expected := "mongodb://testuser:testuserpass@127.0.0.1:27017/test"
	cfg := &storage.Config{
		Hostname:     "",
		Hostnames:    []string{"127.0.0.1"},
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

// TestConnectionURI_ReplSetNoCredentials ensures a correctly formed mongo connection URI is generated for connecting
// to an unsecured mongo replica set.
func TestConnectionURI_ReplSetNoCredentials(t *testing.T) {
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

// TestConnectionURI_ReplSetCredentials ensures a correctly formed mongo connection URI is generated for connecting
// to a mongo replica set with database access credentials.
func TestConnectionURI_ReplSetCredentials(t *testing.T) {
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

// TestNewMongoStore_ErrorsWithBadConfig ensures the underlying lib used for Mongo creates an error
func TestNewMongoStore_ErrorsWithBadConfig(t *testing.T) {
	cfg := &storage.Config{
		Hostname:     "notevenanaddress",
		Port:         27666,
		DatabaseName: "lulz",

		// Specify a low timeout as we know it should fail
		Timeout: 2,
	}
	conn, err := storage.NewMongoStore(cfg, nil)
	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.Nil(t, conn)
}

// TestNewMongoStore_ErrorsWithBadConfig ensures the underlying lib used for Mongo creates an error
func TestNewMongoStore_Succeeds(t *testing.T) {
	cfg := storage.DefaultConfig()
	conn, err := storage.NewMongoStore(cfg, nil)
	assert.Nil(t, err)
	assert.NotNil(t, conn)
	assert.IsType(t, conn, &storage.MongoStore{})
}
