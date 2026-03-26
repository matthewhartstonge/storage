# fosite-storage-mongo
[![Build Status](https://github.com/matthewhartstonge/storage/actions/workflows/go.yaml/badge.svg?branch=development)](https://github.com/matthewhartstonge/storage/actions/workflows/go.yaml) [![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fmatthewhartstonge%2Fstorage.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fmatthewhartstonge%2Fstorage?ref=badge_shield) [![Go Report Card](https://goreportcard.com/badge/github.com/matthewhartstonge/storage)](https://goreportcard.com/report/github.com/matthewhartstonge/storage)

fosite-storage-mongo provides a native Go based [Mongo backed database storage][mongo-driver] 
that conforms to *all the interfaces!* required by [ory/fosite][fosite].
Interface implementations are inspired from the SQL implementations found in [ory/hydra][hydra].

**Table of contents**
- [Compatibility](#compatibility)
- [Development](#development)
    - [Testing](#testing)
- [Examples](#examples)

## Compatibility

| fosite version | storage version | 
|---------------:|----------------:|
|      `v0.49.X` |       `v0.40.X` |

## Development
To start hacking:
* Install [Go][Go] >=1.25
    * Use Go modules!
    * `go build` successfully!

### Testing
Use `go test ./...` to discover heinous crimes against coding!

## Examples
For a quick start check out the following examples based on the `fosite-example`
repo for reference:

- [MongoDB Example](./examples/mongo)

## Licensing
storage is under the Apache 2.0 License.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fmatthewhartstonge%2Fstorage.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fmatthewhartstonge%2Fstorage?ref=badge_large)

[//]: #
    [fosite]: <https://github.com/ory/fosite> 
    [go]: <https://golang.org/dl/>
    [hydra]: <https://github.com/ory/hydra>
    [mongo-driver]: <https://github.com/mongodb/mongo-go-driver>
