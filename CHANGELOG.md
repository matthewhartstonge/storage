# Storage Changelog
## v0.9.1
- Fixes AllowedPeopleAccess filtering.

## v0.9.0
v0.9.0 makes a few under the hood changes in order to conform method and attribute naming to make the API cleaner.

### Mongo Driver
First of all, big shout out to @niemeyer for his amazing effort and continued support through the years to the mgo 
Go driver!! It's no small feat, with the driver in use today in many production stacks. 

We have decided to move to the community supported fork of mgo as it has a couple of extra github issues tidied up and 
is moving to add support for Mongo 3.6 features which make some aggregation pipelines easier internally for us. As such,
this repo is also moving to use the community fork of [mgo][mgo].

Big shoutout to @domodwyer + contributors past and future!

### User
- `AllowedAccess` has been changed to `AllowedTenantAccess` to better represent the underlying data.
    - The `bson`/`json`/`xml` tags have also been updated from `tenantIDs` to `allowedTenantAccess`
- `AllowedPeopleAccess` has been added to the user model support enabling and disabling explicit access to people accounts.
- Added `EnablePeopleAccess` method to user
- Added `DisablePeopleAccess` method to user
- User `AddTenantIDs` method conformed to `EnableTenantAccess` 
- User `RemoveTenantIDs` method conformed to `DisableTenantAccess` 

### Client
- Client `TenantIDs` have been changed to conform to `AllowedTenantAccess`, same as user.
- Client `AddScopes` method has been changed to `EnableScopeAccess`
- Client `RemoveScopes` method has been changed to `DisableScopeAccess`
- Client `AddTenantIDs` method has been changed to `EnableTenantAccess`
- Client `RemoveTenantIDs` method has been changed to `DisableTenantAccess` 

## v0.8.0
- Makes users filterable with `user.Filter` via the `GetUsers(filters user.Filter)` function 

## v0.7.5
- Adds `PersonID` to the client record to enable foreign key lookups 

## v0.7.4
- Adds `TenantIDs` to the client record to enable `client_credentials` for multi-tenant applications

## v0.7.3
- Adds better error checking support for clients

## v0.7.2
- Adds support for disabling clients via the model

## v0.7.1
- Adds functions to enable sorting Clients by Name and Owner  
- Adds functions to enable sorting Users by Username, FirstName and LastName  

## v0.7.0
- Adds support for mongo connections over SSL
- `ConnectionURI` has been dropped in favour of `ConnectionInfo` to enable SSL connections

## v0.6.0
- Removes `request.PersistRefreshTokenGrantSession` from `request.Storer` interface as per required fosite v0.11.x breaking changes
- Removes `request.PersistAuthorizeCodeGrantSession` from `request.Storer` interface as per required fosite v0.11.x breaking changes
- Uses the new interfaces that were brought in to simplify storage with fosite v0.11.x

## v0.5.3
- Add omitempty for marshaling json tags

## v0.5.2
- Added returning `fosite.ErrNotFound` if unable to find a user record to delete

## v0.5.1
- Add omitempty for marshaling tags

## v0.5.0
- Opened the user model up to accept passwords via JSON/XML payloads. 

Ensure that on all API routes, if using the model directly, to either cast attributes to a response struct that does 
not contain a password attribute or clear out the password field before sending the response.

## v0.4.4
- Added error for conflicting user accounts on creation based on username

## v0.4.3
- Fixed a filtering case where organisation_id had not been changed to tenantIDs
- Fixes a couple of testcases

## v0.4.2
- Adds user account disabled boolean. 
- Adds user methods to check for equality and emptiness.

## v0.4.1
- Remove go 1.9 test helper function to enable testing on go 1.7 and go 1.8

## v0.4.0
- Removes user organisationID.
- Adds tenantIDs to the user model to enable multi-tenanted applications  

## v0.3.2
- Updates Storer interface to include the now existing concrete implementations of `RevokeRefreshToken` and `RevokeAccessToken` 
- Adds an edge case test for a single hostname in hostnames

## v0.3.1
- Users
    - Fixes an issue in GetUser() where error checking `err != mgo.ErrNotFound` should have been `err == mgo.ErrNotFound`
    - Fixes error handling being over generous with multi-returns of `errors.withstack(errors.withstack(...))`

## v0.3.0
- Adds support for fosite v0.9.0+

## v0.2.1
- Fixes bug related to findSessionBySignature where mgo requires a MongoRequest struct that has been malloc'd

## v0.2.0
- Make all marshalling conform to JS/JSON camelCase convention

## v0.1.0
- General pre-release!
