package mongo

import (
	// Standard Library Imports
	"testing"

	// External Imports
	"github.com/ory/fosite/handler/oauth2"
)

func TestRequestMongoManager_ImplementsFositeRefreshTokenStorageInterface(t *testing.T) {
	r := &requestMongoManager{}

	var i interface{} = r
	if _, ok := i.(oauth2.RefreshTokenStorage); !ok {
		t.Error("requestMongoManager does not implement interface oauth2.RefreshTokenStorage")
	}
}
