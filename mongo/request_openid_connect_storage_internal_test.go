package mongo

import (
	// Standard Library Imports
	"testing"

	// External Imports
	"github.com/ory/fosite/handler/openid"
)

func TestRequestMongoManager_ImplementsFositeOpenidOpenIDConnectRequestStorageInterface(t *testing.T) {
	r := &requestMongoManager{}

	var i interface{} = r
	if _, ok := i.(openid.OpenIDConnectRequestStorage); !ok {
		t.Error("requestMongoManager does not implement interface openid.OpenIDConnectRequestStorage")
	}
}
