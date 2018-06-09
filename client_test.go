package storage_test

import (
	// Standard Library Imports
	"testing"

	// External Imports
	"github.com/ory/fosite"

	// Internal Imports
	"github.com/matthewhartstonge/storage"
)

func TestClient_ImplementsFositeClientInterface(t *testing.T) {
	c := &storage.Client{}

	var i interface{} = c
	_, ok := i.(fosite.Client)
	if ok != true {
		t.Error("storage.Client does not implement interface fosite.Client")
	}
}
