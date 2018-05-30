package request_test

import (
	"testing"

	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
	"github.com/stretchr/testify/assert"

	"github.com/matthewhartstonge/storage/request"
)

func TestRequestMongoManager_ImplementsFositeCoreStorageInterface(t *testing.T) {
	r := &request.MongoManager{}

	var i interface{} = r
	_, ok := i.(oauth2.CoreStorage)
	assert.Equal(t, true, ok, "request.MongoManager does not implement interface oauth2.CoreStorage")
}

func TestRequestMongoManager_ImplementsFositeAccessTokenStorageInterface(t *testing.T) {
	r := &request.MongoManager{}

	var i interface{} = r
	_, ok := i.(oauth2.AccessTokenStorage)
	assert.Equal(t, true, ok, "request.MongoManager does not implement interface oauth2.AccessTokenStorage")
}

func TestRequestMongoManager_ImplementsFositeAuthorizeCodeStorageInterface(t *testing.T) {
	r := &request.MongoManager{}

	var i interface{} = r
	_, ok := i.(oauth2.AuthorizeCodeStorage)
	assert.Equal(t, true, ok, "request.MongoManager does not implement interface oauth2.AuthorizeCodeStorage")
}

func TestRequestMongoManager_ImplementsFositeClientCredentialsGrantStorageInterface(t *testing.T) {
	r := &request.MongoManager{}

	var i interface{} = r
	_, ok := i.(oauth2.ClientCredentialsGrantStorage)
	assert.Equal(t, true, ok, "request.MongoManager does not implement interface oauth2.ClientCredentialsGrantStorage")
}

func TestRequestMongoManager_ImplementsFositeRefreshTokenStorageInterface(t *testing.T) {
	r := &request.MongoManager{}

	var i interface{} = r
	_, ok := i.(oauth2.RefreshTokenStorage)
	assert.Equal(t, true, ok, "request.MongoManager does not implement interface oauth2.RefreshTokenStorage")
}

func TestRequestMongoManager_ImplementsFositeResourceOwnerPasswordCredentialsGrantStorageInterface(t *testing.T) {
	r := &request.MongoManager{}

	var i interface{} = r
	_, ok := i.(oauth2.ResourceOwnerPasswordCredentialsGrantStorage)
	assert.Equal(t, true, ok, "request.MongoManager does not implement interface oauth2.ResourceOwnerPasswordCredentialsGrantStorage")
}

func TestRequestMongoManager_ImplementsFositeOpenIDConnectRequestStorageInterface(t *testing.T) {
	r := &request.MongoManager{}

	var i interface{} = r
	_, ok := i.(openid.OpenIDConnectRequestStorage)
	assert.Equal(t, true, ok, "request.MongoManager does not implement interface openid.OpenIDConnectRequestStorage")
}
