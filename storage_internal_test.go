package storage

import (
	// Standard Library Imports
	"testing"

	// External Imports
	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/handler/pkce"
)

func TestStore_ImplementsFositeStorage(t *testing.T) {
	c := &Store{}

	var i interface{} = c
	if _, ok := i.(fosite.Storage); !ok {
		t.Error("Store does not implement interface fosite.Storage")
	}
}

func TestStore_ImplementsFositeClientManager(t *testing.T) {
	c := &Store{}

	var i interface{} = c
	if _, ok := i.(fosite.ClientManager); !ok {
		t.Error("Store does not implement interface fosite.ClientManager")
	}
}

func TestStore_ImplementsFositeAccessTokenStorageInterface(t *testing.T) {
	r := &Store{}

	var i interface{} = r
	if _, ok := i.(oauth2.AccessTokenStorage); !ok {
		t.Error("Store does not implement interface oauth2.AccessTokenStorage")
	}
}

func TestStore_ImplementsFositeAuthorizeCodeStorageInterface(t *testing.T) {
	r := &Store{}

	var i interface{} = r
	if _, ok := i.(oauth2.AuthorizeCodeStorage); !ok {
		t.Error("Store does not implement interface oauth2.AuthorizeCodeStorage")
	}
}

func TestStore_ImplementsFositeClientCredentialsGrantStorageInterface(t *testing.T) {
	r := &Store{}

	var i interface{} = r
	if _, ok := i.(oauth2.ClientCredentialsGrantStorage); !ok {
		t.Error("Store does not implement interface oauth2.ClientCredentialsGrantStorage")
	}
}

func TestStore_ImplementsFositeRefreshTokenStorageInterface(t *testing.T) {
	r := &Store{}

	var i interface{} = r
	if _, ok := i.(oauth2.RefreshTokenStorage); !ok {
		t.Error("Store does not implement interface oauth2.RefreshTokenStorage")
	}
}

func TestStore_ImplementsFositeResourceOwnerPasswordCredentialsGrantStorageInterface(t *testing.T) {
	r := &Store{}

	var i interface{} = r
	if _, ok := i.(oauth2.ResourceOwnerPasswordCredentialsGrantStorage); !ok {
		t.Error("Store does not implement interface oauth2.ResourceOwnerPasswordCredentialsGrantStorage")
	}
}

func TestStore_ImplementsFositeOpenidOpenIDConnectRequestStorageInterface(t *testing.T) {
	r := &Store{}

	var i interface{} = r
	if _, ok := i.(openid.OpenIDConnectRequestStorage); !ok {
		t.Error("Store does not implement interface openid.OpenIDConnectRequestStorage")
	}
}

func TestStore_ImplementsFositePkcePKCERequestStorageInterface(t *testing.T) {
	r := &Store{}

	var i interface{} = r
	if _, ok := i.(pkce.PKCERequestStorage); !ok {
		t.Error("Store does not implement interface pkce.PKCERequestStorage")
	}
}
