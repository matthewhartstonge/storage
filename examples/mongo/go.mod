module github.com/ory/fosite-example

go 1.13

// use the local code, rather than go'getting the module
replace github.com/matthewhartstonge/storage => ../../../storage

require (
	github.com/matthewhartstonge/storage v0.18.9
	github.com/ory/fosite v0.32.2
	github.com/sirupsen/logrus v1.4.2
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
)
