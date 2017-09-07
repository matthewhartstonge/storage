# Storage Changelog
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
