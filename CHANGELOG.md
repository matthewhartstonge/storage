# Storage Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Changed
- mongo: Changed `MongoStore` to `Store` to be more idiomatic.
- mongo: Changed `ConnectToMongo` to `Connect` to be more idiomatic.
- mongo: Changed `NewDefaultMongoStore` to `NewDefaultStore` to be more idiomatic.

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

[Unreleased]: https://github.com/matthewhartstonge/storage/tree/master
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
