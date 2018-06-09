package mongo

import (
	// Standard Library Imports
	"testing"

	// External Imports
	"github.com/ory/fosite"

	// Internal Imports
	"github.com/matthewhartstonge/storage"
)

func TestClientMongoManager_ImplementsStorageConfigurer(t *testing.T) {
	c := &clientMongoManager{}

	var i interface{} = c
	_, ok := i.(storage.Configurer)
	if ok != true {
		t.Error("clientMongoManager does not implement interface storage.Configurer")
	}
}

func TestClientMongoManager_ImplementsStorageAuthClientMigrator(t *testing.T) {
	c := &clientMongoManager{}

	var i interface{} = c
	_, ok := i.(storage.AuthClientMigrator)
	if ok != true {
		t.Error("clientMongoManager does not implement interface storage.AuthClientMigrator")
	}
}

func TestClientMongoManager_ImplementsFositeClientManager(t *testing.T) {
	c := &clientMongoManager{}

	var i interface{} = c
	_, ok := i.(fosite.ClientManager)
	if ok != true {
		t.Error("clientMongoManager does not implement interface fosite.ClientManager")
	}
}

func TestClientMongoManager_ImplementsFositeStorage(t *testing.T) {
	c := &clientMongoManager{}

	var i interface{} = c
	_, ok := i.(fosite.Storage)
	if ok != true {
		t.Error("clientMongoManager does not implement interface fosite.Storage")
	}
}

func TestClientMongoManager_ImplementsStorageClientStorer(t *testing.T) {
	c := &clientMongoManager{}

	var i interface{} = c
	_, ok := i.(storage.ClientStorer)
	if ok != true {
		t.Error("clientMongoManager does not implement interface storage.ClientStorer")
	}
}

func TestClientMongoManager_ImplementsStorageClientManager(t *testing.T) {
	c := &clientMongoManager{}

	var i interface{} = c
	_, ok := i.(storage.ClientManager)
	if ok != true {
		t.Error("clientMongoManager does not implement interface storage.ClientManager")
	}
}
