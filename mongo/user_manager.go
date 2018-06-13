package mongo

import (
	// Standard Library Imports
	"context"
	"time"

	// External Imports
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/ory/fosite"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"

	// Internal Imports
	"github.com/matthewhartstonge/storage"
)

// userMongoManager provides a mongo backed implementation for user resources.
//
// Implements:
// - storage.Configurer
// - storage.AuthUserMigrator
// - storage.UserStorer
// - storage.UserManager
type userMongoManager struct {
	db     *mgo.Database
	hasher fosite.Hasher
}

func (u *userMongoManager) Configure(ctx context.Context) error {
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityUsers,
		"method":     "Configure",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = u.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	collection := u.db.C(storage.EntityUsers).With(mgoSession)

	// Ensure Indexes on collections
	index := mgo.Index{
		Name:       IdxUserId,
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := collection.EnsureIndex(index)
	if err != nil {
		log.WithError(err).Error(logError)
		return err
	}

	index = mgo.Index{
		Name:       IdxUsername,
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = collection.EnsureIndex(index)
	if err != nil {
		log.WithError(err).Error(logError)
		return err
	}

	return nil
}

// getConcrete returns an OAuth 2.0 User resource.
func (c *userMongoManager) getConcrete(ctx context.Context, userID string) (storage.User, error) {
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityUsers,
		"method":     "getConcrete",
		"userID":     userID,
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = c.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Build Query
	query := bson.M{
		"id": userID,
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "userMongoManager",
		Method:  "getConcrete",
		Query:   query,
	})
	defer span.Finish()

	result := storage.User{}
	collection := c.db.C(storage.EntityUsers).With(mgoSession)
	if err := collection.Find(query).One(&result); err != nil {
		if err == mgo.ErrNotFound {
			log.WithError(err).Debug(logNotFound)
			return result, fosite.ErrNotFound
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return result, err
	}
	return result, nil
}

func (u *userMongoManager) List(ctx context.Context, filter storage.ListUsersRequest) ([]storage.User, error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityUsers,
		"method":     "List",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = u.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Build Query
	query := bson.M{}
	if filter.AllowedTenantAccess != "" {
		query["allowedTenantAccess"] = filter.AllowedTenantAccess
	}
	if filter.AllowedPersonAccess != "" {
		query["allowedPersonAccess"] = filter.AllowedPersonAccess
	}
	if filter.PersonID != "" {
		query["personId"] = filter.AllowedPersonAccess
	}
	if filter.Username != "" {
		query["username"] = filter.Username
	}
	if len(filter.ScopesIntersection) > 0 {
		query["scopes"] = bson.M{"$all": filter.ScopesUnion}
	}
	if len(filter.ScopesUnion) > 0 {
		query["scopes"] = bson.M{"$in": filter.ScopesUnion}
	}
	if filter.FirstName != "" {
		query["firstName"] = filter.FirstName
	}
	if filter.LastName != "" {
		query["lastName"] = filter.LastName
	}
	if filter.Disabled {
		query["disabled"] = filter.Disabled
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "userMongoManager",
		Method:  "List",
		Query:   query,
	})
	defer span.Finish()

	var results []storage.User
	collection := u.db.C(storage.EntityUsers).With(mgoSession)
	err := collection.Find(query).All(&results)
	if err != nil {
		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return nil, err
	}
	return results, nil
}

func (u *userMongoManager) Create(ctx context.Context, user storage.User) (storage.User, error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityUsers,
		"method":     "Create",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = u.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Enable developers to provide their own IDs
	if user.ID == "" {
		user.ID = uuid.New()
	}
	if user.CreateTime == 0 {
		user.CreateTime = time.Now().Unix()
	}

	// Hash incoming secret
	hash, err := u.hasher.Hash([]byte(user.Password))
	if err != nil {
		log.WithError(err).Error(logNotHashable)
		return user, err
	}
	user.Password = string(hash)

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "userMongoManager",
		Method:  "Create",
	})
	defer span.Finish()

	// Create resource
	collection := u.db.C(storage.EntityUsers).With(mgoSession)
	err = collection.Insert(user)
	if err != nil {
		if mgo.IsDup(err) {
			// Log to StdOut
			log.WithError(err).Debug(logConflict)
			// Log to OpenTracing
			otLogErr(span, err)
			return user, storage.ErrResourceExists
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		user.Password = "REDACTED"
		otLogQuery(span, user)
		otLogErr(span, err)
		return user, err
	}
	return user, nil
}

func (u *userMongoManager) Get(ctx context.Context, userID string) (storage.User, error) {
	return u.getConcrete(ctx, userID)
}

// GetByUsername returns a user resource if found by username.
func (u *userMongoManager) GetByUsername(ctx context.Context, username string) (storage.User, error) {
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityUsers,
		"method":     "GetByUsername",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = u.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Build Query
	query := bson.M{
		"username": username,
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "userMongoManager",
		Method:  "getConcrete",
		Query:   query,
	})
	defer span.Finish()

	result := storage.User{}
	collection := u.db.C(storage.EntityUsers).With(mgoSession)
	if err := collection.Find(query).One(&result); err != nil {
		if err == mgo.ErrNotFound {
			log.WithError(err).Debug(logNotFound)
			return result, fosite.ErrNotFound
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return result, err
	}
	return result, nil
}

func (u *userMongoManager) Update(ctx context.Context, userID string, updatedUser storage.User) (storage.User, error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityUsers,
		"method":     "Update",
		"id":         userID,
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = u.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	currentResource, err := u.getConcrete(ctx, userID)
	if err != nil {
		if err == fosite.ErrNotFound {
			log.Debug(logNotFound)
			return currentResource, err
		}

		log.WithError(err).Error(logError)
		return currentResource, err
	}

	// Deny updating the entity Id
	updatedUser.ID = userID
	// Update modified time
	updatedUser.UpdateTime = time.Now().Unix()

	if currentResource.Password == updatedUser.Password || updatedUser.Password == "" {
		// If the password/hash is blank or hash matches, set using old hash.
		updatedUser.Password = currentResource.Password
	} else {
		newHash, err := u.hasher.Hash([]byte(updatedUser.Password))
		if err != nil {
			log.WithError(err).Error(logNotHashable)
			return currentResource, err
		}
		updatedUser.Password = string(newHash)
	}

	// Build Query
	selector := bson.M{
		"id": userID,
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager:  "userMongoManager",
		Method:   "Update",
		Selector: selector,
	})
	defer span.Finish()

	collection := u.db.C(storage.EntityUsers).With(mgoSession)
	if err := collection.Update(selector, updatedUser); err != nil {
		if err == mgo.ErrNotFound {
			// Log to StdOut
			log.WithError(err).Debug(logNotFound)
			// Log to OpenTracing
			otLogErr(span, err)
			return currentResource, fosite.ErrNotFound
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogQuery(span, updatedUser)
		otLogErr(span, err)
		return currentResource, err
	}
	return updatedUser, nil
}

// Migrate is provided solely for the case where you want to migrate users and
// upgrade their password using the AuthUserMigrator interface.
// This performs an upsert, either creating or overwriting the record with the
// newly provided full record. Use with caution, be secure, don't be dumb.
func (u *userMongoManager) Migrate(ctx context.Context, migratedUser storage.User) (storage.User, error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityUsers,
		"method":     "Migrate",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = u.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Generate a unique ID if not supplied
	if migratedUser.ID == "" {
		migratedUser.ID = uuid.New()
	}
	// Update create time
	if migratedUser.CreateTime == 0 {
		migratedUser.CreateTime = time.Now().Unix()
	}
	// Update modified time
	migratedUser.UpdateTime = time.Now().Unix()

	// Build Query
	selector := bson.M{
		"id": migratedUser.ID,
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager:  "userMongoManager",
		Method:   "Migrate",
		Selector: selector,
	})
	defer span.Finish()

	collection := u.db.C(storage.EntityUsers).With(mgoSession)
	if _, err := collection.Upsert(selector, migratedUser); err != nil {
		if err == mgo.ErrNotFound {
			// Log to StdOut
			log.WithError(err).Debug(logNotFound)
			// Log to OpenTracing
			otLogErr(span, err)
			return storage.User{}, fosite.ErrNotFound
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogQuery(span, migratedUser)
		otLogErr(span, err)
		return storage.User{}, err
	}
	return migratedUser, nil
}

func (u *userMongoManager) Delete(ctx context.Context, userID string) error {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityUsers,
		"method":     "Delete",
		"id":         userID,
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = u.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Build Query
	query := bson.M{
		"id": userID,
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "userMongoManager",
		Method:  "Delete",
		Query:   query,
	})
	defer span.Finish()

	collection := u.db.C(storage.EntityUsers).With(mgoSession)
	if err := collection.Remove(query); err != nil {
		if err == mgo.ErrNotFound {
			// Log to StdOut
			log.WithError(err).Debug(logNotFound)
			// Log to OpenTracing
			otLogErr(span, err)
			return fosite.ErrNotFound
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return err
	}
	return nil
}

func (u *userMongoManager) Authenticate(ctx context.Context, username string, password string) (storage.User, error) {
	return u.AuthenticateByUsername(ctx, username, password)
}

func (u *userMongoManager) AuthenticateByID(ctx context.Context, userID string, password string) (storage.User, error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityUsers,
		"method":     "AuthenticateByID",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = u.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "userMongoManager",
		Method:  "AuthenticateByID",
	})
	defer span.Finish()

	user, err := u.getConcrete(ctx, userID)
	if err != nil {
		log.WithError(err).Warn(logError)
		return user, err
	}

	if user.Disabled {
		log.Debug("disabled user denied access")
		return user, fosite.ErrAccessDenied
	}

	err = u.hasher.Compare([]byte(user.Password), []byte(password))
	if err != nil {
		log.WithError(err).Warn("failed to authenticate user password")
		return user, err
	}

	return user, nil
}

func (u *userMongoManager) AuthenticateByUsername(ctx context.Context, username string, password string) (storage.User, error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityUsers,
		"method":     "AuthenticateByUsername",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = u.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "userMongoManager",
		Method:  "AuthenticateByUsername",
	})
	defer span.Finish()

	user, err := u.GetByUsername(ctx, username)
	if err != nil {
		log.WithError(err).Warn(logError)
		return user, err
	}

	if user.Disabled {
		log.Debug("disabled user denied access")
		return user, fosite.ErrAccessDenied
	}

	err = u.hasher.Compare([]byte(user.Password), []byte(password))
	if err != nil {
		log.WithError(err).Warn("failed to authenticate user password")
		return user, err
	}

	return user, nil
}

// AuthenticateMigration enables developers to supply your own
// authentication function, which in turn, if true, will migrate the secret
// to the hasher implemented within fosite.
func (u *userMongoManager) AuthenticateMigration(ctx context.Context, currentAuth storage.AuthUserFunc, userID string, password string) (storage.User, error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityUsers,
		"method":     "AuthenticateMigration",
		"id":         userID,
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = u.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "userMongoManager",
		Method:  "AuthenticateMigration",
	})
	defer span.Finish()

	// Authenticate with old hasher
	user, authenticated := currentAuth()

	// Check for user not found
	if user.IsEmpty() && !authenticated {
		log.Debug(logNotFound)
		return user, fosite.ErrNotFound
	}

	if user.Disabled {
		log.Debug("disabled user denied access")
		return user, fosite.ErrAccessDenied
	}

	if !authenticated {
		// If user isn't authenticated, try authenticating with new hasher.
		err := u.hasher.Compare(user.GetHashedSecret(), []byte(password))
		if err != nil {
			log.WithError(err).Warn("failed to authenticate user password")
		}
		return user, err
	}

	// If the user is found and authenticated, create a new hash using the new
	// hasher, update the database record and return the record with no error.
	newHash, err := u.hasher.Hash([]byte(password))
	if err != nil {
		log.WithError(err).Error(logNotHashable)
		return user, err
	}

	// Save the new hash
	user.Password = string(newHash)
	return u.Update(ctx, userID, user)
}

func (u *userMongoManager) GrantScopes(ctx context.Context, userID string, scopes []string) (storage.User, error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityUsers,
		"method":     "GrantScopes",
		"id":         userID,
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = u.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "userMongoManager",
		Method:  "GrantScopes",
	})
	defer span.Finish()

	user, err := u.getConcrete(ctx, userID)
	if err != nil {
		if err == fosite.ErrNotFound {
			log.Debug(logNotFound)
			return user, err
		}

		log.WithError(err).Error(logError)
		return user, err
	}

	user.EnableScopeAccess(scopes...)
	return u.Update(ctx, user.ID, user)
}

func (u *userMongoManager) RemoveScopes(ctx context.Context, userID string, scopes []string) (storage.User, error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityUsers,
		"method":     "RemoveScopes",
		"id":         userID,
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = u.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "userMongoManager",
		Method:  "RemoveScopes",
	})
	defer span.Finish()

	user, err := u.getConcrete(ctx, userID)
	if err != nil {
		if err == fosite.ErrNotFound {
			log.Debug(logNotFound)
			return user, err
		}

		log.WithError(err).Error(logError)
		return user, err
	}

	user.DisableScopeAccess(scopes...)
	return u.Update(ctx, user.ID, user)
}
