package request_test

import (
	"testing"

	"github.com/ory/fosite/handler/openid"
	"github.com/stretchr/testify/assert"

	"github.com/matthewhartstonge/storage/request"
)

func TestRequestMongoManager_ImplementsFositeOpenIDConnectRequestStorageInterface(t *testing.T) {
	r := &request.MongoManager{}

	var i interface{} = r
	_, ok := i.(openid.OpenIDConnectRequestStorage)
	assert.Equal(t, true, ok, "request.MongoManager does not implement interface openid.OpenIDConnectRequestStorage")
}
