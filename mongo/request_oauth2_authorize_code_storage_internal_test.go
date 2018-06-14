package mongo

import (
	// Standard Library Imports
	"testing"

	// External Imports
	"github.com/ory/fosite/handler/oauth2"
)

func TestRequestMongoManager_ImplementsFositeAuthorizeCodeStorageInterface(t *testing.T) {
	r := &RequestManager{}

	var i interface{} = r
	if _, ok := i.(oauth2.AuthorizeCodeStorage); !ok {
		t.Error("RequestManager does not implement interface oauth2.AuthorizeCodeStorage")
	}
}
