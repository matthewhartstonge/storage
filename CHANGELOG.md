# Storage Changelog
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
