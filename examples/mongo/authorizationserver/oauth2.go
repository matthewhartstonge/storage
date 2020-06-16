package authorizationserver

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"time"

	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
)

func RegisterHandlers() {
	// Set up oauth2 endpoints. You could also use gorilla/mux or any other router.
	http.HandleFunc("/oauth2/auth", authEndpoint)
	http.HandleFunc("/oauth2/token", tokenEndpoint)

	// revoke tokens
	http.HandleFunc("/oauth2/revoke", revokeEndpoint)
	http.HandleFunc("/oauth2/introspect", introspectionEndpoint)
}

// This secret is used to sign access and refresh tokens as well as authorize codes.
// It has to be 32-bytes long for HMAC signing.
// In order to generate secure keys, the best thing to do is use crypto/rand:
//
// ```
// package main
//
// import (
//	"crypto/rand"
//	"encoding/hex"
//	"fmt"
// )
//
// func main() {
//	var secret = make([]byte, 32)
//	_, err := rand.Read(secret)
//	if err != nil {
//		panic(err)
//	}
// }
// ```
//
// If you require this to key to be stable, for example, when running multiple fosite servers, you can generate the
// 32byte random key as above and push it out to a base64 encoded string.
// This can then be injected and decoded as the `var secret []byte` on server start.
var secret = []byte("some-cool-secret-that-is-32bytes")

// check the api docs of compose.Config for further configuration options
var config = &compose.Config{
	AccessTokenLifespan: time.Minute * 30,
	// ...
}

// NewExampleMongoStore allows us to create an example Mongo datastore that will
// panic if you don't have an unauthenticated mongo database that can be found
// at `localhost:27017`.
var store = NewExampleMongoStore()

var oauth2 = compose.ComposeAllEnabled(config, store, secret, mustRSAKey())

// A session is passed from the `/auth` to the `/token` endpoint. You probably want to store data like: "Who made the request",
// "What organization does that person belong to" and so on.
// For our use case, the session will meet the requirements imposed by JWT access tokens, HMAC access tokens and OpenID Connect
// ID Tokens plus a custom field

// newSession is a helper function for creating a new session. This may look like a lot of code but since we are
// setting up multiple strategies it is a bit longer.
// Usually, you could do:
//
//  session = new(fosite.DefaultSession)
func newSession(user string) *openid.DefaultSession {
	return &openid.DefaultSession{
		Claims: &jwt.IDTokenClaims{
			Issuer:      "https://fosite.my-application.com",
			Subject:     user,
			Audience:    []string{"https://my-client.my-application.com"},
			ExpiresAt:   time.Now().Add(time.Hour * 6),
			IssuedAt:    time.Now(),
			RequestedAt: time.Now(),
			AuthTime:    time.Now(),
		},
		Headers: &jwt.Headers{
			Extra: make(map[string]interface{}),
		},
	}
}

func mustRSAKey() *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	return key
}
