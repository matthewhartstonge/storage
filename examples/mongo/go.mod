module github.com/ory/fosite-example

go 1.14

// use the local code, rather than go'getting the module
replace github.com/matthewhartstonge/storage => ../../../storage

require (
	github.com/matthewhartstonge/storage v0.0.0
	github.com/ory/fosite v0.33.0
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/net v0.12.0
	golang.org/x/oauth2 v0.10.0
)
