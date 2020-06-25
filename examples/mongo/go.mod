module github.com/ory/fosite-example

go 1.13

// use the local code, rather than go'getting the module
replace github.com/matthewhartstonge/storage => ../../../storage

require (
	github.com/golang/mock v1.3.1 // indirect
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.1.0 // indirect
	github.com/matthewhartstonge/storage v0.18.9
	github.com/ory/fosite v0.30.2
	github.com/sirupsen/logrus v1.4.2
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4 // indirect
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/appengine v1.6.1 // indirect
	gopkg.in/square/go-jose.v2 v2.2.2 // indirect
)
