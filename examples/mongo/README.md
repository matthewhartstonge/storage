# ORY Fosite Example Server (feat. Mongo)

[![Build Status](https://travis-ci.org/ory/fosite-example.svg?branch=master)](https://travis-ci.org/ory/fosite-example)

ORY Fosite is the security first OAuth2 & OpenID Connect framework for Go. Built simple, powerful and extensible. This repository contains an exemplary http server using ORY Fosite for serving OAuth2 requests.

## Prerequisites
The mongo demo expects that you have a local (default) instance of mongo running.
This means, a mongo listening at `localhost:27017` with no authentication enabled.

The easiest way to get this up and running is if you have docker pre-installed.

### Docker
In a terminal run the following:

```sh
docker run -d -p 27017:27017 mongo:8.0
```

### Local Installation
Install the community edition of MongoDB on your computer locally following the steps from [mongo's documentation site](https://docs.mongodb.com/manual/installation/#mongodb-community-edition-installation-tutorials)

## Install and run
The Fosite example server requires [`go@1.23` or higher installed](https://golang.org/dl/) as it uses go modules for dependency management. 
Once Go and mongo have been installed, run the demo:

```
$ git clone https://github.com/matthewhartstonge/storage.git
$ cd storage/example/mongo
$ go run main.go
```
