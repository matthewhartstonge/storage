package mongo

import (
	// Standard Library Imports
	"testing"

	// External Imports
	"github.com/ory/fosite/handler/pkce"
)

func TestRequestMongoManager_ImplementsFositePkcePKCERequestStorageInterface(t *testing.T) {
	r := &RequestManager{}

	var i interface{} = r
	if _, ok := i.(pkce.PKCERequestStorage); !ok {
		t.Error("RequestManager does not implement interface pkce.PKCERequestStorage")
	}
}
