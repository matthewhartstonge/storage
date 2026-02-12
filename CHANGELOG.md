# Storage Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [0.39.2](https://github.com/matthewhartstonge/storage/compare/v0.39.1...v0.39.2) (2026-02-12)


### Bug Fixes

* **deps:** bump the default group across 2 directories with 3 updates ([#117](https://github.com/matthewhartstonge/storage/issues/117)) ([cdd33ed](https://github.com/matthewhartstonge/storage/commit/cdd33edad8650ede2ff9e1200917ea7b43e4ae86))

## [0.39.1](https://github.com/matthewhartstonge/storage/compare/v0.39.0...v0.39.1) (2026-01-26)


### Bug Fixes

* **deps:** bump go.mongodb.org/mongo-driver ([#114](https://github.com/matthewhartstonge/storage/issues/114)) ([0546fd1](https://github.com/matthewhartstonge/storage/commit/0546fd1a6a710a0560be379e1c1855d1ecf1fe55))

## [0.39.0](https://github.com/matthewhartstonge/storage/compare/v0.38.0...v0.39.0) (2026-01-12)


### Features

* Change user custom data field type from `any` to `json.RawMessage` to facilitate easier retrieval from db.  ([#109](https://github.com/matthewhartstonge/storage/issues/109)) ([62d3dd9](https://github.com/matthewhartstonge/storage/commit/62d3dd9cf91a463f3e28da4a489885107fa1cd40))

## [0.38.0](https://github.com/matthewhartstonge/storage/compare/v0.37.0...v0.38.0) (2025-11-20)


### Features

* **build:** updates minimum Go version to `go@1.24`. ([6a7c975](https://github.com/matthewhartstonge/storage/commit/6a7c97551201456761d67cd087735e3b66a2d446))

## [0.37.0](https://github.com/matthewhartstonge/storage/compare/v0.36.0...v0.37.0) (2025-08-07)


### Features

* **client:** adds support for `fosite.ClientWithSecretRotation`. ([6760715](https://github.com/matthewhartstonge/storage/commit/67607158ff4635b9d0c46e976dcc49529b90dde1))
* **client:** adds support for `fosite.ResponseModeClient`. ([e80fdb6](https://github.com/matthewhartstonge/storage/commit/e80fdb63666c92ed22a3d39a9d036b481b1e922f))


## [v0.36.0] - 2025-07-28

This jumps from `fosite@v0.35.1` => `fosite@v0.49.0` and with it comes a number of breaking changes.

### Breaking Changes

Also mentioned in the sections below, but highlighted here with relevant migration information:

- Requires `>=go@1.23`.
- `fosite.Hasher` has been removed from the individual entity managers (`ClientManager`, `UserManager`) in favour of using the hasher provided by the shared DB instance. You will need to reroute usage to the top level hasher (`store.Hasher`), or via the manager's DB shared instance:
    - `store.ClientManager.Hasher.*` => `store.ClientManager.DB.Hasher.*`
    - `store.UserManager.Hasher.*` => `store.UserManager.DB.Hasher.*`
- mongo: normalized `time.Now()` usage throughout to `UTC`. Traditionally, Go will use `time.Local()` which may not be useful working with external systems.
- mongo: The interface for `fosite.RefreshTokenStorage` has been updated and now requires the access token signature which **MUST** be hashed with `storage.SignatureHash(signature string) string`:

```diff
- func (r *RequestManager) CreateRefreshTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error)
+ func (r *RequestManager) CreateRefreshTokenSession(ctx context.Context, refreshSignature string, accessSignature string, request fosite.Requester) (err error)
```

- mongo: A number of changes have been made to all indices, therefore, all `storage` indices will require _manual_ removal, but on service startup the indices will be recreated as required across all collections.
  - `sparse` indexing has been removed.
  - The hashed `idxSignatureId` index has been removed in favour of performing internal hashing of access token signatures.
  - The unique requirement has been relaxed on the `idxSessionID` index and as such will need to be removed.

Use the following `mongosh` script for quick index removal:

```js
// Connect to your database if you haven't already
// For example:
// use myFositeDatabase;

const indexesToDrop = [
    // AccessTokens Collection
    { collection: "accessTokens", index: "idxSignatureId" },
    { collection: "accessTokens", index: "idxSessionId" },
    { collection: "accessTokens", index: "idxCompoundRequester" },
    { collection: "accessTokens", index: "idxExpiryRequestedAt" },
    // AuthorizationCodes Collection
    { collection: "authorizationCodes", index: "idxSignatureId" },
    { collection: "authorizationCodes", index: "idxSessionId" },
    { collection: "authorizationCodes", index: "idxCompoundRequester" },
    { collection: "authorizationCodes", index: "idxExpiryRequestedAt" },
    // Clients Collection
    { collection: "clients", index: "idxClientId" },
    // JtiDenylist Collection
    { collection: "jtiDenylist", index: "idxSignatureId" },
    { collection: "jtiDenylist", index: "idxExpires" },
    { collection: "jtiDenylist", index: "idxExpiryRequestedAt" },
    // OpenIDConnectSessions Collection
    { collection: "openIDConnectSessions", index: "idxSignatureId" },
    { collection: "openIDConnectSessions", index: "idxSessionId" },
    { collection: "openIDConnectSessions", index: "idxCompoundRequester" },
    { collection: "openIDConnectSessions", index: "idxExpiryRequestedAt" },
    // PkceSessions Collection
    { collection: "pkceSessions", index: "idxSignatureId" },
    { collection: "pkceSessions", index: "idxSessionId" },
    { collection: "pkceSessions", index: "idxCompoundRequester" },
    { collection: "pkceSessions", index: "idxExpiryRequestedAt" },
    // RefreshTokens Collection
    { collection: "refreshTokens", index: "idxSignatureId" },
    { collection: "refreshTokens", index: "idxSessionId" },
    { collection: "refreshTokens", index: "idxCompoundRequester" },
    { collection: "refreshTokens", index: "idxExpiryRequestedAt" },
    // Users Collection
    { collection: "users", index: "idxUserId" },
    { collection: "users", index: "idxUsername" }
];

function dropIndex(collectionName, indexName) {
    try {
        print(`Attempting to drop index '${indexName}' from collection '${collectionName}'...`);
        const result = db.getCollection(collectionName).dropIndex(indexName);
        if (result.ok === 1) {
            print(`Successfully dropped index '${indexName}' from collection '${collectionName}'.`);
        } else {
            print(`Failed to drop index '${indexName}' from collection '${collectionName}'. Result: ${JSON.stringify(result)}`);
        }
    } catch (e) {
        if (e.code === 27) { // 27 is the error code for IndexNotFound
            print(`Index '${indexName}' not found on collection '${collectionName}'. Skipping.`);
        } else {
            print(`Error dropping index '${indexName}' from collection '${collectionName}': ${e}`);
        }
    }
}

// Iterate through the array and drop each index
indexesToDrop.forEach(item => {
    if (item.collection && item.index) {
        dropIndex(item.collection, item.index);
    } else {
        print(`Skipping invalid entry in indexesToDrop array: ${JSON.stringify(item)}. Missing 'collection' or 'index' property.`);
    }
});

print("\nIndex removal script complete.");
```

### Added
- config: `CONNECTIONS_MONGO_REFRESH_TOKEN_GRACE_PERIOD` can be configured to set a multiple-use graceful token refresh window. Beneficial when working with web-based clients with multiple open tabs. Default: `0 == Not Enabled`.
- config: `CONNECTIONS_MONGO_REFRESH_TOKEN_MAX_USAGE` can be configured to enforce the maximum number of times a refresh token can be used. Default: `0 == unlimited`.
- `storage.SignatureHash(signature string) string` for hashing access token signatures to keep indexes small.
- `store.RequestManager.DeleteAll(ctx context.Context, entityName string, requestID string) (err error)`: to handle removing all records based on `requestID` (a given session) at once to cater for graceful token refreshing.
- `store.RequestManager.RotateRefreshToken(ctx context.Context, requestID string, refreshTokenSignature string) (err error)`: to support the latest `fosite.RefreshTokenStorage` interface definition.
- `user.Data` enables persisting custom data in a user record.
- `client.Data` enables persisting custom data in a client record.

### Changed
- deps!: upgrades to `fosite@v0.49.0`.
- mongo!: the `SessionID` index has been relaxed and is no longer unique to allow for graceful token refreshes.
- mongo!: routes `fosite.Hasher` through the shared singleton DB instance to simplify hasher plumbing.
- mongo!: access token signatures are now being directly hashed via `storage.SignatureHash` internally so we no longer need the hashed `#signature` index.
- mongo!: sparse indexes have been removed. The indexes built always had the specific properties required, so never required being sparse.
- mongo!: The interface for `fosite.RefreshTokenStorage` has been updated and now requires the access token signature which **MUST** be hashed with `storage.SignatureHash(signature string) string`:

```diff
- func (r *RequestManager) CreateRefreshTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error)
+ func (r *RequestManager) CreateRefreshTokenSession(ctx context.Context, refreshSignature string, accessSignature string, request fosite.Requester) (err error)
```

### Fixed
- mongo!: normalized `time.Now()` usage throughout to UTC.
- mongo: fixes `hashee` assignment ordering. There was a potential that the default hasher could have been `nil`.

### Removed
- deps: removed dependency on `github.com/pkg/errors`.
- mongo: as mentioned above, the `#signature` index has been removed in favour of internally hashing the signature before commiting to storage.
- mongo: internal function `configureExpiry` as no longer used.

## [v0.35.0] - 2025-07-21
### Changed
- deps: updates to `go@1.23`.
- deps: updates to `github.com/ory/fosite@v0.34.1`.
- deps: updates to `go.mongodb.org/mongo-driver@v1.17.4`.
- ci: added testing against `go@1.23, 1.24`.
- ci: added testing against `mongo@8.0`.
- ci: removed testing against `go@1.13 - go@1.22`.
- ci: removed testing against `mongo@4.4 - mongo@6.0`.

## [v0.34.0] - 2024-02-22
### Added
- user: adds support for regions.

### Changed
- deps: updates to `github.com/google/uuid@v1.6.0`.
- deps: updates to `github.com/stretchr/testify@v1.8.2`.
- deps: updates to `go.mongodb.org/mongo-driver@v1.13.2`.
- ci: added testing against `go@1.22`.
- ci: added testing against `mongo@7.0`.
- ci: removed testing against `go@1.13`.

## [v0.33.0] - 2023-08-25
### Fixed
- user: aligns `bson`, `json` and `xml` marshalling to the name of the `mfaFactors` property.  

## [v0.32.0] - 2023-07-18
### Added
- user: Adds support for MFA factors.

### Changed
- deps: updates to `github.com/opentracing/opentracing-go@v1.2.0`.
- deps: updates to `github.com/sirupsen/logrus@v1.9.3`.
- deps: updates to `go.mongodb.org/mongo-driver@v1.12.0`.
- deps: updates `examples/mongo` to `github.com/sirupsen/logrus@v1.9.3`.
- deps: updates `examples/mongo` to `golang.org/x/net@v0.12.0`.
- deps: updates `examples/mongo` to `golang.org/x/oauth2@v0.10.0`.

## [v0.31.0] - 2023-01-10
### Changed
- deps: updates to `github.com/google/uuid@v1.3.0`.
- deps: updates to `github.com/sirupsen/logrus@v1.8.1`.
- deps: updates to `github.com/stretchr/testify@v1.7.0`.
- deps: updates to `go.mongodb.org/mongo-driver@v1.11.1`.
- deps: updates `examples/mongo` to `github.com/sirupsen/logrus@v1.8.1`.
- deps: updates `examples/mongo` to `golang.org/x/net@v0.0.0-20220926192436-02166a98028e`.
- deps: updates `examples/mongo` to `golang.org/x/oauth2@v0.0.0-20220909003341-f21342109be1`.
- deps: updates `examples/mongo` to `go.mongodb.org/mongo-driver@v1.11.1`.
- examples/mongo/authorizationserver: migrates deprecated use of `Exact()` to `ExactOne()`.
- storage: gofmts the project with go@1.19.

### Fixed
- examples/mongo/authorizationserver: sets session subject and username. fixes [#65](https://github.com/matthewhartstonge/storage/issues/65).
- examples/mongo/authorizationserver: properly logs out the generated user id.
- mongo/mongo: reduces read errors occurring in a mongo replica set. fixes [#68](https://github.com/matthewhartstonge/storage/issues/68).

## [v0.30.1] - 2022-08-08
### Added
- user_manager: adds support for filtering users given a list of people ids.
- mongo/user_manager: adds support for filtering users given a list of people ids.

## [v0.30.0] - 2022-07-28
### Changed
- deps: upgrades to `fosite@v0.33.0`.

## [v0.29.0] - 2022-07-28
*Breaking Change:*
    If you are running on Mongo<4.0, please update as the indices will now
    build in the foreground. Mongo>4.0 has changed to a new indexing engine
    and this option is now deprecated.

### Removed
- mongo: deprecates `SetBackground` due to MongoDB 4.0 EOL.

## [v0.28.0] - 2021-10-18
### Added
- mongo: adds support for `mongodb+srv` connection strings.
- mongo: binds in a default TLS Config if `ssl=true` and a TLS config has not been provided.
- storage: adds `Expirer` interface to enable stores to add support for configuring record expiration.
- mongo: implements `storage.Expirer` interface to enable TTL based expiry on tokens.

### Changed
- mongo: migrated internal use of `isDup(err)` to `mongo.IsDuplicateKeyError(err)`.

### Removed
- mongo: removed internal `isDup(err)` function.

## [v0.27.0] - 2021-09-24
This release will add a new hashed index on `signature` for the `accessTokens`
collection. This makes the old `accessTokens.idxSignatureId` index redundant and
can be removed.

### Added
- mongo: migrates to using a hashed index for the signature index on access tokens.
    - The signature for an access token could grow quite large, leading to a
      large index.  By migrating to using a hashed index, the size can be
      reduced to 2% of the original indices size. In our testing we went from
      2.8GB -> 57MB.

### Fixed
- examples/mongo/authorizationserver: removes `mongo-features` example.

## [v0.26.0] - 2021-08-05
### Added
- utils: adds functions to help with adding and removing items from string sets.
- user: adds test cases for enabling and disabling person access.
- user: adds tests for `user.FullName()`.
- user: adds test cases to check create time and update time equality.
- user: adds test cases to check equality of allowed person ids, person id and extra fields in user record.
- user: adds support for storing user roles.
- storage: adds a benchmark for `user.Equal()`.

### Changed
- user: refactors enable and disable functions to use util append/remove functions.
- examples/mongo: updates `go.mod` to `go@v1.14` and tidies `go.sum`.

### Fixed
- mongo: `SetClientAssertionJWT` now logs unknown errors if deleting expired JTIs fails.
- mongo: fixes do not pass a nil Context (staticcheck)
- user: fixes whitespace issues when returning a user's full name.

### Removed
- deps: removed support for dep.

## [v0.25.1] - 2021-07-27
### Changed
- deps: updates to `mongo-driver@v1.5.4`.
    - This mongo driver release contains a fix to prevent clearing server 
      connection pools on operation-scoped timeouts.

## [v0.25.0] - 2021-06-01
### Added
- README: updates documentation.
    - Adds links to download Go.
    - Adds information for working with Go modules.
    - Changes build badge link to travis-ci.com.
    - Changes mgo link to the official MongoDB mongo-driver.

### Changed
- deps: migrates from `pborman/uuid` to `google/uuid`.
- deps: updates dependencies.
    - updates module to go1.14 (go@n-2).
    - updates to `mongo-driver@v1.5.2`.
    - updates to `testify@v1.6.1`.
    - migrates from `pborman/uuid` to `google/uuid@v1.2.0`.
    - removed `mongo-features@v0.4.0`.
- .travis: removes `go@1.13`, adds `go@1.16`.

### Fixed
- mongo: not found on token revocation should return nil.
- .travis: go install goveralls binary.

### Removed
- transactional: removes transactional interface implementation.
    - There isn't an easy way to tell via the mongo driver if the version of 
      mongo running is compatible with transactions (>mongo 4.4) without 
      enabling admin commands to be run for example, `db.adminCommand( { getParameter: 1, featureCompatibilityVersion: 1 } )`.
      Therefore, for now, it's easier to remove it until every current mongo
      version supports transactions.
- deps: removes use of `mongo-features` due to bugfix released via `mongo-driver`.
    - `mongo-driver` wasn't pulling or pushing sessions into the context correctly.
    - `mongo-features` also relied on admin commands/permissions to detect the
      running mongo version to ascertain if the mongo version connected to was
      transaction compatible, so no longer needed.

## [v0.24.0] - 2020-09-02
### Breaking changes
As mentioned under changed:
- `AuthClientFunc` and `AuthUserFunc` now take in a context.
- `store.DB` is now of type `*DB` not `*mongo.Database` but the API remains the
  same. If you explicitly require type `*mongo.Database`, you can obtain this by
  stepping into the `DB` wrapper `store.DB.Database`.

### Added
- deps: adds `mongo-features@v0.4.0` for mongoDB feature detection.
- mongo: adds `DB` a wrapper containing `*mongo.Database` and `*feat.Features`.
- mongo: implements mongo feature detection for correct session and transaction handling.

### Changed
- storage: `AuthClientFunc` and `AuthUserFunc` now accept a context.
    - `type AuthClientFunc func() (Client, bool)` => `type AuthClientFunc func(ctx context.Context) (Client, bool)` 
    - `type AuthUserFunc func() (User, bool)` => `type AuthUserFunc func(ctx context.Context) (User, bool)` 
- mongo: all handlers have moved from `DB *mongo.Database` to our wrapper 
  `DB *DB` in order to provide mongoDB feature detection for managing sessions 
  and transactions, if available.
- examples/mongo/authorizationserver: puts session creation behind a feature flag.

## [v0.23.0] - 2020-08-27
Deprecated - don't use.

The session and transaction implementation does not work for single node users 
(i.e. mongo not running as a replicaset), or those using mongo <v4.0.0.

### Added
- mongo: implements `storage.Transactional`

### Changed
- deps: upgrades to `mongo-driver@v1.3.7`

## [v0.22.2] - 2020-07-06
### Fixed
- mongo: fixes `UserManager.Migrate` returning not found on a successful insert.

## [v0.22.1] - 2020-07-06
### Fixed
- mongo: fixes `filter.ScopesIntersection` using `filter.ScopesUnion` instead
  of `filter.ScopesIntersection`.

## [v0.22.0] - 2020-07-02
### Changed
- deps: upgrades to `fosite@v0.32.2`

## [v0.21.0] - 2020-07-02
### Added
- storage: added support for managing and denying JTIs due to newly added 
  methods in `fosite@v0.31.X`'s interface `fosite.ClientManager`.
- mongo: added concrete implementation for `DeniedJTIManager` and 
  `DeniedJTIStorer` to comply to added methods in `fosite.ClientManager`.
- mongo: ensured update time is updated when updates are performed.
- mongo: added config options to adjust mongo connection min/max pool size.

### Changed
- deps: upgrades to `fosite@v0.31.3`
- readme: added version support information for `storage@v0.20.X` 
- readme: added version support information for `storage@v0.21.X`

### Removed
- storage: removed missed entity constants that helped define cache 
  table/schema/collection.
- mongo: removed dead-code index constants resulting from the removal of the 
  cache collection.

## [v0.20.0] - 2020-06-26
### Breaking changes
Removes 'Cache' implementation which actually added a level of indirection, 
doubling required database calls in some instances.

### Changed
- mongo: uses a defined database for testing.
- examples/mongo/authorizationserver: uses a defined database for the demo.

### Fixed
- travisci: fixes travis not running tests over the whole code base.

### Removed
- cache: removed cache structure, interfaces and db 
- `storage.SessionCache` (struct)
- `storage.Cacher` (interface)
- `storage.CacheManager` (interface)
- `storage.CacheStorer` (interface)
- `storage.RequestManager.Cache` (interface binding to a `storage.CacheStorer`)
- `mongo.CacheManager` (concrete implementation of `storage.CacheManager`)

## [v0.19.0] - 2020-06-26
### Breaking changes
This release migrates to the official Go MongoDB driver.

If you have any custom code using mgo that feeds into `storage`, you will need 
to migrate these to use [mongo-go-driver][mongo-go-driver] patterns.

### Added
- examples/mongo: added [fosite-example](./examples/mongo) featuring mongo 
  integration.

### Changed
- deps: updates to `fosite@v0.30.2`.
- deps: migrates from `globalsign/mgo` to `mongodb/mongo-go-driver`.
- readme: references `examples/mongo` instead of having a wad of example code 
  in the readme.

## [v0.18.9] - 2020-06-13
### Fixed
- mongo: `RevokeAccessToken` attempted to delete the access token twice from 
  the datastore leading to `fosite.ErrNotFound` always being returned.
- mongo: `RevokeRefreshToken` attempted to delete the refresh token twice from 
  the datastore leading to `fosite.ErrNotFound` always being returned.

## [v0.18.8] - 2020-06-11
### Fixed
- mongo: auth codes should be set to active by default on creation.

## [v0.18.7] - 2020-05-24
### Changed
- travisci: updated to test against `go@{1.14, tip}`

### Fixed
- mongo: fixed `ineffassign` and `staticcheck` issues.
- mongo: fixed `maligned` issues reducing config struct memory allocation from 
  138 bytes to 127 bytes.
- mongo: fixed missed error check.
- mongo: fixed `lint` issues where context was not the first parameter.
- mongo: fixed user delete test creating a client instead of a user for 
  deletion.
- mongo: fixed create client parameter ordering.

### Removed
- travisci: support for go < 1.13

## [v0.18.6] - 2019-09-25
### Added
- client: added `published` to enable filtering clients by published state.

## [v0.18.5] - 2019-09-24
### Changed
- deps:  updated to `fosite@v0.30.1`

### Fixed
- client: fixes `client.Equal` by doing a compare on allowed regions.

## [v0.18.4] - 2019-09-24
### Added
- client: added support for allowed regions. This enables filtering for clients 
  based on geographic region.
- mongo: added tests for `client.list`.

### Changed
- travis: updated CI testing to test against go versions `1.13.x`, `1.12.x`, 
  `1.11.x`.
- travis: migrated to go modules for dependency management.
- deps: updated to `fosite@v0.29.8` and `opentracing-go@v1.1.0`.

### Removed
- client: removed redundant type conversions in various return statements.

## [v0.18.3] - 2019-09-11
### Fixed
- mongo: fixes OpenTracing logging in the `cache` storage manager.

## [v0.18.2] - 2019-09-11
### Fixed
- Calls to `Cache.Get` and `Cache.Delete` in the `RevokeAccessToken` and
  `RevokeRefreshToken` handlers were specified in the wrong order.

## [v0.18.1] - 2019-02-07
### Added
- experimental support for go modules.

### Fixed
- Fixed the last ineffassign issue reported via goreportcard.
- Tested against upstream fosite@v0.28.x
- Tested against upstream fosite@v0.29.x
- RequestManager: `RequestManager.List` now uses `entityName` instead of 
  hardcoded `storage.EntityClients` [#24](https://github.com/matthewhartstonge/storage/issues/24)
- RequestManager: `RequestManager.Update` should use `entityName` instead of 
  hardcoded `storage.EntityClients` [#25](https://github.com/matthewhartstonge/storage/issues/25)

## [v0.18.0] - 2019-01-24
### Added
- Support for testing under Go 1.11

### Changed
- Adds support for Fosite `v0.27.x`
- `Client`: Now has an `AllowedAudiences` attribute to comply to the new 
  interface method required for `fosite.Client`.
- `Request`: Changed attribute `Scopes` to `RequestedScope`. bson, json and xml 
  tags remain the same.
- `Request`: Changed attribute `GrantedScopes` to `GrantedScope`. bson, json and
   xml tags remain the same.

### Fixed
- Fixes the last golint error which was not reported when run locally.
- Fixes ineffassign issues reported via goreportcard.

### Removed
- Support for testing under Go 1.8

## [v0.17.0] - 2018-11-07
### Changed
- Adds support for Fosite `v0.26.0`
- Exported Mongo index constants have been changed to align with idiomatic Go, 
  where the `Id` suffixes are now `ID`

### Fixed
- Fixed all golint errors

## [v0.16.0] - 2018-10-15
### Changed
- Adds support for Fosite `v0.25.0`

## [v0.15.0] - 2018-10-15
### Changed
- Adds support for Fosite `v0.23.0`

## [v0.14.0] - 2018-10-15
### Changed
- Adds support for Fosite `v0.22.0`
- Updated readme example to match upstream.

## [v0.13.0-beta] - 2018-09-04
We have been using this release in house for the past month with our own auth 
server. If you have any issues related to the mongo storage implementation, 
please report an issue. 

### Changed
- deps: updated `Gopkg.lock` to support dep `v0.5.0`

### Fixed
- mongo: Have added a struct tag to tell the `envconfig` package to ignore 
    processing `Config.TLSConfig`, as the instantiated config it creates breaks 
    TLS mongo connections.
- user manager: Fixes filtering not being performed on `PersonID`  

## [v0.13.0-alpha2] - 2018-07-12
### Added
- mongo: Added tests to CacheManager for Create, Get, Update, Delete and 
    DeleteByValue.

### Changed
- CacheManager: must support `Configurer` interface
- RequestManager: must support `Configurer` interface
- deps: updated to support fosite `v0.21.X`

### Fixed
- readme: version link for `v0.13.0-alpha1` 
- default config: Fixed a configuration bug, where repeat connections would 
    lead to the default port being appended multiple times to cfg.hostnames.
- cachemanager: DeleteByValue's query selector should have been querying by
    attribute `signature` not by the non-existant bson attribute `value`.  
- requestmanager: Reverted session data back to []byte due to not being able to 
    unmarshal into an interface. 

## [v0.13.0-alpha1] - 2018-06-18
### Added
- mongo: Added New to re-support custom mongo configuration and hashers.  
- Store: Added top level function override to enable `Store` to conform to 
	the required `fosite` interfaces.
- Store: Added interface tests to `Store` to ensure the functions are available 
	at the top level!

### Changed
- Storer: Changed `storage.Storer` to `Storage.Store` to be more idiomatic.
- Storer: Changed from named interfaces to a struct composition to enable the 
	required fosite interface functions to be raised to the top level.
- mongo: Changed `MongoStore` to `Store` to be more idiomatic.
- mongo: Changed `ConnectToMongo` to `Connect` to be more idiomatic.
- mongo: Changed `NewDefaultMongoStore` to `NewDefaultStore` to be more idiomatic.
- mongo: exported `cacheMongoManager` 
- mongo: exported `clientMongoManager` 
- mongo: exported `requestMongoManager` 
- mongo: exported `userMongoManager` 
- mongo: Changed `CacheMongoManager` to `CacheManager` to be more idiomatic.
- mongo: Changed `ClientMongoManager` to `ClientManager` to be more idiomatic.
- mongo: Changed `RequestMongoManager` to `RequestManager` to be more idiomatic.
- mongo: Changed `UserMongoManager` to `UserManager` to be more idiomatic.
- mongo: Changed unexported attributes `db` and `hasher` to be exported to 
	enable custom data store composition.

### Fixed
- documentation: typo in `user_manager` referring to clients instead of users.

## [v0.13.0-alpha] - 2018-06-14
### Added
- mongo: Added indices for quick look up.
- mongo: Added a way to pass the mongo session via `context.Context`
- OpenTracing: Added OpenTracing support. You can now gain distributed tracing 
    information on how your mongo queries are performing.   
- logging: Added logging support. Now provides a way to bind in your own logrus 
    logger to get information, or debug output from the storage driver.
- Client: Added to the domain model. Provides a data storage model for OAuth 2
    clients.
- AuthClientMigrator: Added to the domain model. Provides an interface to help 
    enable migration of hashes for legacy clients.
- AuthUserMigrator: Added to the domain model. Provides an interface to help 
    enable migration of hashes for legacy users.
- Configurer: Provides an interface to initialize datastore entities if 
    required.
- Cache: Added to the domain model. Provides caching functionality.
- Tests: Added fosite interface tests to easily test API compatibility with 
    newer version of fosite.
- Users: Added to the domain model. Provides a data storage model for OAuth 2 
    users.
- Requests: Added to the domain model. Provides a data storage model for OAuth 2
    auth session requests.
- Entity Names: Added to the domain model. Provides a way to use the entity 
    names consistently between multiple backend storage implementations.
- Storer: Added to the domain model. Provides a struct for composing backend 
    storage drivers. See `MongoStore` for an example of how to bind this in.
- AuthorizeCode: Added support for `InvalidateAuthorizeCodeSession(ctx context.Context, code string) (err error)`
    as per `fosite@v0.20.X`
- fosite: v0.20.X support.

### Changed
Pretty much everything.. 

Storage has been re-written in such a way that multiple datastore backends can 
be created, and bound together. This can become useful over time as you need to 
scale out and would like to switch components out to a different backend.

For example, hitting the cache. You could implement and compose in a Redis 
`CacheManager`, which you could bind into your mongo storage implementation.  

- OSS: Updated licenses and added attributions.
- Client: Secret is now stored as a string rather than bytes.
- Configurer: requires passing in `context.Context`

### Removed
- `DeleteAuthorizeCodeSession(ctx context.Context, code string) (err error)` has
    been removed from the interface and is no longer used by the upstream, 
    fosite, library.
- The old API

### Known Issues
- There are no mongo integration tests.
- Documentation needs to be updated to match current API.

## [v0.12.0] - 2018-05-31
### Added
- client: Tests to ensure storage implements fosite interfaces correctly
- request: Tests to ensure storage implements oauth2 interfaces correctly
- request: Tests to ensure storage implements openid interfaces correctly
- readme: compatibility table

### Changed:
- license: updated year, added github link.
- mongo: conformed collection names to match javascript naming conventions 
    (camelCase)
- deps: changed dependency manager to dep
- ci: changed TravisCI dependency manager to dep
- file naming: removed package name pre-pending to file names.

### Fixed
- Tested against upstream fosite@v0.12.0
- Tested against upstream fosite@v0.13.0
- Tested against upstream fosite@v0.14.0
- Tested against upstream fosite@v0.15.0
- Tested against upstream fosite@v0.16.0

### Removed
- request: Removed CreateImplicitAccessTokenSession function as per github 
    issue [removed implicit storage as its never used](https://github.com/ory/fosite/pull/171)
- storage_mongo: Removed CreateImplicitAccessTokenSession function as per 
    github issue [removed implicit storage as its never used](https://github.com/ory/fosite/pull/171)

## [v0.11.2] - 2018-05-30
### Changed
- git: updated repo links
- deps: updated glide lock

## [v0.11.1] - 2018-05-14
### Changed
- readme: updated latest version

### Fixed
- user: Equal() now supports comparisons including personID

### Removed
- legal: Removed mergo, now not in use

## [v0.11.0] - 2018-05-10
### Changed
- user: Removed use of lib mergo. Please move to passing through a full update, 
	instead of partials. This caused issues where fields were required to be 
	blanked out, for example, disabling a user account. 
- client: Removed use of lib mergo. Please move to passing through a full update, 
	instead of partials. This caused issues where fields were required to be 
	blanked out, for example, disabling a client.
- changelog: to be changelog compliant!
- glide: unpinned fosite version. Please ensure it works with your version of 
	fosite, please see readme disclaimers.

### Removed
- glide: mergo

## [v0.10.0] - 2018-04-13
### Changed
- Configuration now allows passing hostnames with included ports, for example: 
    `[]string{"mongo.example.com:123456", "mongo.example.com:234567"}`allowing 
    developers to bypass having to configure `config.Port` as well.
- Configuration now allows passing a custom tls.Config to the Config. This 
    requires manual initialization of a `tls.Config` struct, but enables users 
    to use their own TLS certs for connecting to mongo.
- Cleaned up the Readme

## [v0.9.1] - 2018-03-19
### Fixed
- Fixes AllowedPeopleAccess filtering.

## [v0.9.0] - 2018-03-19
v0.9.0 makes a few under the hood changes in order to conform method and 
attribute naming to make the API cleaner.

### Mongo Driver
First of all, big shout out to @niemeyer for his amazing effort and continued 
support through the years to the mgo Go driver!! It's no small feat, with the 
driver in use today in many production stacks. 

We have decided to move to the community supported fork of mgo as it has a 
couple of extra github issues tidied up and is moving to add support for Mongo 
3.6 features which make some aggregation pipelines easier internally for us. 
As such, this repo is also moving to use the community fork of [mgo][mgo].

Big shoutout to @domodwyer + contributors past and future!

### Added
- User: `AllowedPeopleAccess` has been added to the user model support enabling and disabling explicit access to people accounts.
- User: Added `EnablePeopleAccess` method to user
- User: Added `DisablePeopleAccess` method to user

### Changed
- User:`AllowedAccess` has been changed to `AllowedTenantAccess` to better represent the underlying data.
    - The `bson`/`json`/`xml` tags have also been updated from `tenantIDs` to `allowedTenantAccess`
- User: User `AddTenantIDs` method conformed to `EnableTenantAccess` 
- User: User `RemoveTenantIDs` method conformed to `DisableTenantAccess` 
- Client: `TenantIDs` have been changed to conform to `AllowedTenantAccess`, same as user.
- Client: `AddScopes` method has been changed to `EnableScopeAccess`
- Client: `RemoveScopes` method has been changed to `DisableScopeAccess`
- Client: `AddTenantIDs` method has been changed to `EnableTenantAccess`
- Client: `RemoveTenantIDs` method has been changed to `DisableTenantAccess` 

## [v0.8.0] - 2018-03-16
- Makes users filterable with `user.Filter` via the `GetUsers(filters user.Filter)` function 

## [v0.7.5] - 2017-10-12
### Added
- Adds `PersonID` to the client record to enable foreign key lookups 

## [v0.7.4] - 2017-10-06
### Added
- Adds `TenantIDs` to the client record to enable `client_credentials` for multi-tenant applications

## [v0.7.3] - 2017-10-03
### Added
- Adds better error checking support for clients

## [v0.7.2] - 2017-10-03
### Added
- Adds support for disabling clients via the model

## [v0.7.1] - 2017-10-03
### Added
- Adds functions to enable sorting Clients by Name and Owner  
- Adds functions to enable sorting Users by Username, FirstName and LastName  

## [v0.7.0] - 2017-10-02
### Added
- Adds support for mongo connections over SSL

### Removed
- `ConnectionURI` has been dropped in favour of `ConnectionInfo` to enable SSL connections

## [v0.6.0] - 2017-10-02
### Changed
- Uses the new interfaces that were brought in to simplify storage with fosite v0.11.x

### Removed
- Removes `request.PersistRefreshTokenGrantSession` from `request.Storer` interface as per required fosite v0.11.x breaking changes
- Removes `request.PersistAuthorizeCodeGrantSession` from `request.Storer` interface as per required fosite v0.11.x breaking changes

## [v0.5.3] - 2017-09-19
### Added
- Add omitempty for marshaling json tags

## [v0.5.2] - 2017-09-18
### Added
- Added returning `fosite.ErrNotFound` if unable to find a user record to delete

## [v0.5.1] - 2017-09-18
### Added
- Add omitempty for marshaling tags

## [v0.5.0] - 2017-09-18
### Added
- Opened the user model up to accept passwords via JSON/XML payloads. 

Ensure that on all API routes, if using the model directly, to either cast 
attributes to a response struct that does not contain a password attribute or 
clear out the password field before sending the response.

## [v0.4.4]  - 2017-09-18
### Added
- Added error for conflicting user accounts on creation based on username

## [v0.4.3] - 2017-09-15
### Fixed
- Fixed a filtering case where organisation_id had not been changed to tenantIDs
- Fixes a couple of testcases

## [v0.4.2] - 2017-09-11
### Added
- Adds user account disabled boolean. 
- Adds user methods to check for equality and emptiness.

## [v0.4.1] - 2017-09-08
### Removed
- Remove go 1.9 test helper function to enable testing on go 1.7 and go 1.8

## [v0.4.0] - 2017-09-07
### Added
- Adds tenantIDs to the user model to enable multi-tenanted applications  

### Removed
- Removes user organisationID.

## [v0.3.2] - 2017-07-10
### Added
- Adds an edge case test for a single hostname in hostnames

### Changed
- Updates Storer interface to include the now existing concrete implementations of `RevokeRefreshToken` and `RevokeAccessToken` 

## [v0.3.1] - 2017-06-08
### Fixed
- Users
    - Fixes an issue in GetUser() where error checking `err != mgo.ErrNotFound` should have been `err == mgo.ErrNotFound`
    - Fixes error handling being over generous with multi-returns of `errors.withstack(errors.withstack(...))`

## [v0.3.0] - 2017-06-07
### Changed
- Adds support for fosite v0.9.0+

## [v0.2.1] - 2017-06-02
### Fixed
- Fixes bug related to findSessionBySignature where mgo requires a MongoRequest struct that has been malloc'd

## [v0.2.0] - 2017-06-02
### Changed
- Make all marshalling conform to JS/JSON camelCase convention

## [v0.1.0] - 2017-05-31
### Added
- General pre-release!

[v0.36.0]: https://github.com/matthewhartstonge/storage/tree/v0.36.0
[v0.35.0]: https://github.com/matthewhartstonge/storage/tree/v0.35.0
[v0.34.0]: https://github.com/matthewhartstonge/storage/tree/v0.34.0
[v0.33.0]: https://github.com/matthewhartstonge/storage/tree/v0.33.0
[v0.32.0]: https://github.com/matthewhartstonge/storage/tree/v0.32.0
[v0.31.0]: https://github.com/matthewhartstonge/storage/tree/v0.31.0
[v0.30.1]: https://github.com/matthewhartstonge/storage/tree/v0.30.1
[v0.30.0]: https://github.com/matthewhartstonge/storage/tree/v0.30.0
[v0.29.0]: https://github.com/matthewhartstonge/storage/tree/v0.29.0
[v0.28.0]: https://github.com/matthewhartstonge/storage/tree/v0.28.0
[v0.27.0]: https://github.com/matthewhartstonge/storage/tree/v0.27.0
[v0.26.0]: https://github.com/matthewhartstonge/storage/tree/v0.26.0
[v0.25.1]: https://github.com/matthewhartstonge/storage/tree/v0.25.1
[v0.25.0]: https://github.com/matthewhartstonge/storage/tree/v0.25.0
[v0.24.0]: https://github.com/matthewhartstonge/storage/tree/v0.24.0
[v0.23.0]: https://github.com/matthewhartstonge/storage/tree/v0.23.0
[v0.22.2]: https://github.com/matthewhartstonge/storage/tree/v0.22.2
[v0.22.1]: https://github.com/matthewhartstonge/storage/tree/v0.22.1
[v0.22.0]: https://github.com/matthewhartstonge/storage/tree/v0.22.0
[v0.21.0]: https://github.com/matthewhartstonge/storage/tree/v0.21.0
[v0.20.0]: https://github.com/matthewhartstonge/storage/tree/v0.20.0
[v0.19.0]: https://github.com/matthewhartstonge/storage/tree/v0.19.0
[v0.18.9]: https://github.com/matthewhartstonge/storage/tree/v0.18.9
[v0.18.8]: https://github.com/matthewhartstonge/storage/tree/v0.18.8
[v0.18.7]: https://github.com/matthewhartstonge/storage/tree/v0.18.7
[v0.18.6]: https://github.com/matthewhartstonge/storage/tree/v0.18.6
[v0.18.5]: https://github.com/matthewhartstonge/storage/tree/v0.18.5
[v0.18.4]: https://github.com/matthewhartstonge/storage/tree/v0.18.4
[v0.18.3]: https://github.com/matthewhartstonge/storage/tree/v0.18.3
[v0.18.2]: https://github.com/matthewhartstonge/storage/tree/v0.18.2
[v0.18.1]: https://github.com/matthewhartstonge/storage/tree/v0.18.1
[v0.18.0]: https://github.com/matthewhartstonge/storage/tree/v0.18.0
[v0.17.0]: https://github.com/matthewhartstonge/storage/tree/v0.17.0
[v0.16.0]: https://github.com/matthewhartstonge/storage/tree/v0.16.0
[v0.15.0]: https://github.com/matthewhartstonge/storage/tree/v0.15.0
[v0.14.0]: https://github.com/matthewhartstonge/storage/tree/v0.14.0
[v0.13.0-beta]: https://github.com/matthewhartstonge/storage/tree/v0.13.0-beta
[v0.13.0-alpha2]: https://github.com/matthewhartstonge/storage/tree/v0.13.0-alpha2
[v0.13.0-alpha1]: https://github.com/matthewhartstonge/storage/tree/v0.13.0-alpha1
[v0.13.0-alpha]: https://github.com/matthewhartstonge/storage/tree/v0.13.0-alpha
[v0.12.0]: https://github.com/matthewhartstonge/storage/tree/v0.12.0
[v0.11.2]: https://github.com/matthewhartstonge/storage/tree/v0.11.2
[v0.11.1]: https://github.com/matthewhartstonge/storage/tree/v0.11.1
[v0.11.0]: https://github.com/matthewhartstonge/storage/tree/v0.11.0
[v0.10.0]: https://github.com/matthewhartstonge/storage/tree/v0.10.0
[v0.9.1]: https://github.com/matthewhartstonge/storage/tree/v0.9.1
[v0.9.0]: https://github.com/matthewhartstonge/storage/tree/v0.9.0
[v0.8.0]: https://github.com/matthewhartstonge/storage/tree/v0.8.0
[v0.7.5]: https://github.com/matthewhartstonge/storage/tree/v0.7.5
[v0.7.4]: https://github.com/matthewhartstonge/storage/tree/v0.7.4
[v0.7.3]: https://github.com/matthewhartstonge/storage/tree/v0.7.3
[v0.7.2]: https://github.com/matthewhartstonge/storage/tree/v0.7.2
[v0.7.1]: https://github.com/matthewhartstonge/storage/tree/v0.7.1
[v0.7.0]: https://github.com/matthewhartstonge/storage/tree/v0.7.0
[v0.6.0]: https://github.com/matthewhartstonge/storage/tree/v0.6.0
[v0.5.3]: https://github.com/matthewhartstonge/storage/tree/v0.5.3
[v0.5.2]: https://github.com/matthewhartstonge/storage/tree/v0.5.2
[v0.5.1]: https://github.com/matthewhartstonge/storage/tree/v0.5.1
[v0.5.0]: https://github.com/matthewhartstonge/storage/tree/v0.5.0
[v0.4.4]: https://github.com/matthewhartstonge/storage/tree/v0.4.4
[v0.4.3]: https://github.com/matthewhartstonge/storage/tree/v0.4.3
[v0.4.2]: https://github.com/matthewhartstonge/storage/tree/v0.4.2
[v0.4.1]: https://github.com/matthewhartstonge/storage/tree/v0.4.1
[v0.4.0]: https://github.com/matthewhartstonge/storage/tree/v0.4.0
[v0.3.2]: https://github.com/matthewhartstonge/storage/tree/v0.3.2
[v0.3.1]: https://github.com/matthewhartstonge/storage/tree/v0.3.1
[v0.3.0]: https://github.com/matthewhartstonge/storage/tree/v0.3.0
[v0.2.1]: https://github.com/matthewhartstonge/storage/tree/v0.2.1
[v0.2.0]: https://github.com/matthewhartstonge/storage/tree/v0.2.0
[v0.1.0]: https://github.com/matthewhartstonge/storage/tree/v0.1.0

[mongo-go-driver]: https://github.com/mongodb/mongo-go-driver
[mgo]: https://github.com/globalsign/mgo
