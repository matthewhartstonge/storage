module github.com/ory/fosite-example

go 1.14

// use the local code, rather than go'getting the module
replace github.com/matthewhartstonge/storage => ../../../storage

require (
	github.com/matthewhartstonge/storage v0.0.0
	github.com/ory/fosite v0.33.0
	github.com/sirupsen/logrus v1.4.2
	golang.org/x/net v0.0.0-20200625001655-4c5254603344
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
)
