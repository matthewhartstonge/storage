package storage_test

import (
	// Standard Library Imports
	"net"
	"testing"

	// External Imports
	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/assert"

	// Internal Imports
	"github.com/matthewhartstonge/storage"
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
	expected := &mgo.DialInfo{
		Addrs:    []string{"127.0.0.1:27017"},
		Direct:   false,
		Timeout:  10000000000,
		FailFast: false,
		Database: "test",
	}

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
	got := storage.ConnectionInfo(cfg)
	assert.Equal(t, expected, got)
}

// TestConnectionURI_SingleHostCredentials ensures a correctly formed mongo connection URI is generated for connecting
// to a single mongo host with database access credentials.
func TestConnectionURI_SingleHostCredentials(t *testing.T) {
	expected := &mgo.DialInfo{
		Addrs:    []string{"127.0.0.1:27017"},
		Direct:   false,
		Timeout:  10000000000,
		FailFast: false,
		Database: "test",
		Username: "testuser",
		Password: "testuserpass",
	}

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
	got := storage.ConnectionInfo(cfg)
	assert.Equal(t, expected, got)
}

// TestConnectionURI_Hostnames ensures a correctly formed mongo connection URI is generated when a single hostname is
// passed as hostnames.
func TestConnectionURI_MultiHostnames(t *testing.T) {
	expected := &mgo.DialInfo{
		Addrs:    []string{"127.0.0.1:27017"},
		Direct:   false,
		Timeout:  10000000000,
		FailFast: false,
		Database: "test",
		Username: "testuser",
		Password: "testuserpass",
	}

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
	got := storage.ConnectionInfo(cfg)
	assert.EqualValues(t, expected, got)
}

// TestConnectionURI_ReplSetNoCredentials ensures a correctly formed mongo connection URI is generated for connecting
// to an unsecured mongo replica set.
func TestConnectionURI_ReplSetNoCredentials(t *testing.T) {
	expected := &mgo.DialInfo{
		Addrs:          []string{"127.0.0.1:27017", "127.0.1.1:27017", "127.0.2.1:27017"},
		Timeout:        10000000000,
		Database:       "test",
		ReplicaSetName: "sr0",
	}

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
	got := storage.ConnectionInfo(cfg)
	assert.EqualValues(t, expected, got)
}

// TestConnectionURI_ReplSetCredentials ensures a correctly formed mongo connection URI is generated for connecting
// to a mongo replica set with database access credentials.
func TestConnectionURI_ReplSetCredentials(t *testing.T) {
	expected := &mgo.DialInfo{
		Addrs:          []string{"127.0.0.1:27017", "127.0.1.1:27017", "127.0.2.1:27017"},
		Timeout:        10000000000,
		Database:       "test",
		ReplicaSetName: "sr0",
		Username:       "testuser",
		Password:       "testuserpass",
	}

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
	got := storage.ConnectionInfo(cfg)
	assert.EqualValues(t, expected, got)
}

type SSLDialer interface {
	DialServer() func(addr *mgo.ServerAddr) (net.Conn, error)
}

// TestConnectionURI_SSLOption ensures a correctly formed mongo connection URI is generated for connecting
// to a SSL enabled mongo with database access credentials.
func TestConnectionURI_SSL(t *testing.T) {
	expected := &mgo.DialInfo{
		Addrs:    []string{"127.0.0.1:27017", "127.0.1.1:27017", "127.0.2.1:27017"},
		Timeout:  10000000000,
		Database: "test",
		Username: "testuser",
		Password: "testuserpass",
	}

	cfg := &storage.Config{
		Hostname:     "",
		Hostnames:    []string{"127.0.0.1", "127.0.1.1", "127.0.2.1"},
		Port:         27017,
		DatabaseName: "test",

		// Credential Access
		Username: "testuser",
		Password: "testuserpass",

		// Options
		SSL: true,
	}
	got := storage.ConnectionInfo(cfg)
	assert.NotNil(t, got.DialServer)
	got.DialServer = nil
	assert.Equal(t, expected, got)
}

// TestConnectionURI_SSLReplica ensures a correctly formed mongo connection URI is generated for connecting
// to a SSL enabled mongo replica set with database access credentials.
func TestConnectionURI_SSLReplica(t *testing.T) {
	expected := &mgo.DialInfo{
		Addrs:          []string{"127.0.0.1:27017", "127.0.1.1:27017", "127.0.2.1:27017"},
		Timeout:        10000000000,
		Database:       "test",
		Username:       "testuser",
		Password:       "testuserpass",
		ReplicaSetName: "sr0",
	}

	cfg := &storage.Config{
		Hostname:     "",
		Hostnames:    []string{"127.0.0.1", "127.0.1.1", "127.0.2.1"},
		Port:         27017,
		DatabaseName: "test",

		// Credential Access
		Username: "testuser",
		Password: "testuserpass",

		// Options
		Replset: "sr0",
		SSL:     true,
	}
	got := storage.ConnectionInfo(cfg)
	assert.NotNil(t, got.DialServer)
	got.DialServer = nil
	assert.Equal(t, expected, got)
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
