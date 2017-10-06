# fosite-storage-mongo
[![Build Status](https://travis-ci.org/MatthewHartstonge/storage.svg?branch=master)](https://travis-ci.org/MatthewHartstonge/storage) [![Coverage Status](https://coveralls.io/repos/github/MatthewHartstonge/storage/badge.svg?branch=master)](https://coveralls.io/github/MatthewHartstonge/storage?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/MatthewHartstonge/storage)](https://goreportcard.com/report/github.com/MatthewHartstonge/storage)

fosite-storage-mongo provides Mongo backed database storage that conforms to *all the interfaces!* required by fosite.

**Lastest Version:** v0.7.4

**Table of contents**
- [Documentation](#documentation)
  - [Development](#development)
    - [Testing](#testing)
  - [Example](#example)

## Documentation
We wanted a [Fosite][fosite]/[Hydra][hydra]* storage backend that supported MongoDB. 'Nuf said.

### Development
To start hacking:
* Install [glide][glide] - A golang package manager
* Run `glide install`
* `go build` successfully!

#### Testing
We use `go test $(glide novendor)` to discover our heinous crimes against coding. For those developing on windows, like 
ourselves, there is a slight problem with linux based commandline expansion for obvious reasons... 

In order to test correctly under windows, neither `go test` or `glide novendor` work happily together. For this reason, 
please run `./test.bat` which has been manually created to achieve what `glide novendor` does. 

### Example
Following the [fosite-example/authorizationserver][fosite-example-server] example, we can extend this to add support 
for Mongo storage via the compose configuration.

```go
package authorizationserver

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"time"

	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
	"github.com/MatthewHartstonge/storage"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

func RegisterHandlers() {
	// Set up oauth2 endpoints. You could also use gorilla/mux or any other router.
	http.HandleFunc("/oauth2/auth", authEndpoint)
	http.HandleFunc("/oauth2/token", tokenEndpoint)

	// revoke tokens
	http.HandleFunc("/oauth2/revoke", revokeEndpoint)
	http.HandleFunc("/oauth2/introspect", introspectionEndpoint)
}

// NewExampleMongoStore allows us to create an example Mongo Datastore and panics if you don't have an unauthenticated 
// mongo database that can be found at `localhost:27017`. NewExampleMongoStore has one Client and one User. Check out 
// storage.NewExampleMongoStore() for the implementation/specific client/user details.
var store = storage.NewExampleMongoStore()
var config = new(compose.Config)

// Because we are using oauth2 and open connect id, we use this little helper to combine the two in one
// variable.
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

	// be aware that open id connect factories need to be added after oauth2 factories to work properly.
	compose.OpenIDConnectExplicitFactory,
	compose.OpenIDConnectImplicitFactory,
	compose.OpenIDConnectHybridFactory,
)

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
			Issuer:    "https://fosite.my-application.com",
			Subject:   user,
			Audience:  "https://my-client.my-application.com",
			ExpiresAt: time.Now().Add(time.Hour * 6),
			IssuedAt:  time.Now(),
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

type stackTracer interface {
	StackTrace() errors.StackTrace
}
```


**Disclaimers**
* We haven't tested implementation with Hydra at all, but we implemented on top of fosite which Hydra uses 
  store it's data under the hood.

[//]: #
    [glide]: <https://glide.sh>
    [fosite]: <https://github.com/ory/fosite> 
    [hydra]: <https://github.com/ory/hydra>
    [fosite-example-server]: <https://github.com/ory/fosite-example/blob/master/authorizationserver/oauth2.go>