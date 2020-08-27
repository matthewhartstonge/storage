package mongo

import (
	"testing"

	"github.com/ory/fosite/storage"
)

// Ensure store implements storage.Transactional
func TestStore_ImplementsStorageTransactional(t *testing.T) {
	u := &Store{}

	var i interface{} = u
	if _, ok := i.(storage.Transactional); !ok {
		t.Error("Store does not implement interface storage.Transactional")
	}
}
