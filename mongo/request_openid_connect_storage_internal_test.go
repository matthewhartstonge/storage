package mongo

import (
	// Standard Library Imports
	"testing"

	// External Imports
	"github.com/ory/fosite/handler/openid"
)

func TestRequestMongoManager_ImplementsFositeOpenidOpenIDConnectRequestStorageInterface(t *testing.T) {
	r := &RequestManager{}

	var i interface{} = r
	if _, ok := i.(openid.OpenIDConnectRequestStorage); !ok {
		t.Error("RequestManager does not implement interface openid.OpenIDConnectRequestStorage")
	}
}
