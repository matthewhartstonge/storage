# fosite-storage-mongo
[![Coverage Status](https://coveralls.io/repos/github/matthewhartstonge/storage/badge.svg?branch=master)](https://coveralls.io/github/matthewhartstonge/storage?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/matthewhartstonge/storage)](https://goreportcard.com/report/github.com/matthewhartstonge/storage) [![Build Status](https://travis-ci.org/matthewhartstonge/storage.svg?branch=master)](https://travis-ci.org/matthewhartstonge/storage) [![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fmatthewhartstonge%2Fstorage.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fmatthewhartstonge%2Fstorage?ref=badge_shield)

fosite-storage-mongo provides a native Go based [Mongo backed database storage][mgo] 
that conforms to *all the interfaces!* required by [fosite][fosite].

**Table of contents**
- [Compatibility](#compatibility)
- [Development](#development)
    - [Testing](#testing)
- [Example](#example)
- [Disclaimer](#disclaimer)

## Compatibility
The following table lists the compatible versions of fosite-storage-mongo with
fosite. If you are currently using this in production, it would be awesome to 
know what versions you are successfully paired with.

| storage version | minimum fosite version | maximum fosite version | 
|----------------:|-----------------------:|-----------------------:|
|       `v0.18.X` |              `v0.27.X` |              `v0.30.X` |
|       `v0.17.X` |              `v0.26.X` |              `v0.26.X` |
|       `v0.16.X` |              `v0.25.X` |              `v0.25.X` |
|       `v0.15.X` |              `v0.23.X` |              `v0.24.X` |
|       `v0.14.X` |              `v0.22.X` |              `v0.22.X` |
|       `v0.13.X` |              `v0.20.X` |              `v0.21.X` |
|       `v0.12.X` |              `v0.11.0` |              `v0.16.X` |
|       `v0.11.X` |              `v0.11.0` |              `v0.16.X` |

## Development
To start hacking:
* Install [dep][dep] - A golang package manager
* Run `dep ensure`
* `go build` successfully!

### Testing
Since Go 1.9, we use `go test ./...` to discover our heinous crimes against 
coding.

## Example
Following the [fosite-example/authorizationserver][fosite-example-server] 
example, we can extend this to add support for Mongo storage via the compose 
configuration.

```go
package authorizationserver

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"time"

	"github.com/matthewhartstonge/storage"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
)

func RegisterHandlers() {
	// Set up oauth2 endpoints. 
	// You could also use gorilla/mux or any other router.
	http.HandleFunc("/oauth2/auth", authEndpoint)
	http.HandleFunc("/oauth2/token", tokenEndpoint)

	// revoke tokens
	http.HandleFunc("/oauth2/revoke", revokeEndpoint)
	http.HandleFunc("/oauth2/introspect", introspectionEndpoint)
}

// NewExampleMongoStore allows us to create an example Mongo Datastore and 
// panics if you don't have an unauthenticated mongo database that can be found 
// at `localhost:27017`. NewExampleMongoStore has one Client and one User. 
// Check out storage.NewExampleMongoStore() for the implementation/specific 
// client/user details.
var store = storage.NewExampleMongoStore()
var config = new(compose.Config)

// Because we are using oauth2 and open connect id, we use this little helper 
// to combine the two in one variable.
var strat = compose.CommonStrategy{
	// alternatively you could use:
	//  OAuth2Strategy: compose.NewOAuth2JWTStrategy(mustRSAKey())
	CoreStrategy: compose.NewOAuth2HMACStrategy(config, []byte("some-super-cool-secret-that-nobody-knows")),

	// open id connect strategy
	OpenIDConnectTokenStrategy: compose.NewOpenIDConnectStrategy(mustRSAKey()),
}

var oauth2 = compose.Compose(
	config,
	store,
	strat,

	// enabled handlers
	compose.OAuth2AuthorizeExplicitFactory,
	compose.OAuth2AuthorizeImplicitFactory,
	compose.OAuth2ClientCredentialsGrantFactory,
	compose.OAuth2RefreshTokenGrantFactory,
	compose.OAuth2ResourceOwnerPasswordCredentialsFactory,

	compose.OAuth2TokenRevocationFactory,
	compose.OAuth2TokenIntrospectionFactory,

	// Be aware that open id connect factories need to be added after oauth2 
	// factories to work properly.
	compose.OpenIDConnectExplicitFactory,
	compose.OpenIDConnectImplicitFactory,
	compose.OpenIDConnectHybridFactory,
	compose.OpenIDConnectRefreshFactory,
)

// A session is passed from the `/auth` to the `/token` endpoint. You probably 
// want to store data like: "Who made the request", "What organization does 
// that person belong to" and so on.
// For our use case, the session will meet the requirements imposed by JWT 
// access tokens, HMAC access tokens and OpenID Connect ID Tokens plus a custom 
// field.

// newSession is a helper function for creating a new session. This may look 
// like a lot of code but since we are setting up multiple strategies it is a 
// bit longer.
//
// Usually, you could do:
// session = new(fosite.DefaultSession)
//
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
```


## Disclaimer
* We are currently using this project in house with Storage `v0.16.0` and Fosite
  `v0.25.x` with good success.
* My aim is to keep storage to date with Fosite releases, as always though, my 
    time is limited due to my human frame. 
* If you are able to provide help in keeping storage up to date, feel free to 
    raise a github issue and discuss where you are able/willing to help. I'm 
    always happy to review PRs and merge code in :ok_hand:
* We haven't tested implementation with Hydra at all but theoretically this 
    should be compatible as Hydra uses Fosite to store it's data under the hood.

## Licensing
storage is under the Apache 2.0 License.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fmatthewhartstonge%2Fstorage.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fmatthewhartstonge%2Fstorage?ref=badge_large)

[//]: #
    [mgo]: <https://github.com/globalsign/mgo>
    [dep]: <https://github.com/golang/dep>
    [fosite]: <https://github.com/ory/fosite> 
    [hydra]: <https://github.com/ory/hydra>
    [fosite-example-server]: <https://github.com/ory/fosite-example/blob/master/authorizationserver/oauth2.go>
