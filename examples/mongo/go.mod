module github.com/ory/fosite-example

go 1.14

// use the local code, rather than go'getting the module
replace github.com/matthewhartstonge/storage => ../../../storage

require (
	github.com/matthewhartstonge/storage v0.0.0
	github.com/ory/fosite v0.33.0
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/net v0.0.0-20220926192436-02166a98028e
	golang.org/x/oauth2 v0.0.0-20220909003341-f21342109be1
)
