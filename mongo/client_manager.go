package mongo

import (
	// Standard Library imports
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

// clientMongoManager provides a fosite storage implementation for Clients.
//
// Implements:
// - fosite.Storage
// - fosite.ClientManager
// - storage.AuthClientMigrator
// - storage.ClientManager
// - storage.ClientStorer
type clientMongoManager struct {
	db     *mgo.Database
	hasher fosite.Hasher
}

// Configure sets up the Mongo collection for OAuth 2.0 client resources.
func (c *clientMongoManager) Configure() error {
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": CollectionClients,
		"method":     "Configure",
	})

	collection := c.db.C(CollectionClients).With(c.db.Session.Clone())
	defer collection.Database.Session.Close()

	// Ensure Indexes on collections
	index := mgo.Index{
		Name:       IdxClientId,
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

	return nil
}

// getConcrete returns an OAuth 2.0 Client resource.
func (c *clientMongoManager) getConcrete(ctx context.Context, clientID string) (storage.Client, error) {
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": CollectionClients,
		"method":     "getConcrete",
		"clientID":   clientID,
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
		"id": clientID,
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "clientMongoManager",
		Method:  "getConcrete",
		Query:   query,
	})
	defer span.Finish()

	result := storage.Client{}
	collection := c.db.C(CollectionClients).With(mgoSession)
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

// List filters resources to return a list of OAuth 2.0 client resources.
func (c *clientMongoManager) List(ctx context.Context, filter storage.ListClientsRequest) ([]storage.Client, error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": CollectionClients,
		"method":     "List",
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = c.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Build Query
	query := bson.D{}
	if filter.TenantID != "" {
		query = append(query, bson.DocElem{Name: "tenantIds", Value: filter.TenantID})
	}
	if filter.RedirectURI != "" {
		query = append(query, bson.DocElem{Name: "redirectUris", Value: filter.RedirectURI})
	}
	if filter.GrantType != "" {
		query = append(query, bson.DocElem{Name: "grantTypes", Value: filter.GrantType})
	}
	if filter.ResponseType != "" {
		query = append(query, bson.DocElem{Name: "responseTypes", Value: filter.ResponseType})
	}
	if filter.Scope != "" {
		query = append(query, bson.DocElem{Name: "scopes", Value: filter.Scope})
	}
	if filter.Contact != "" {
		query = append(query, bson.DocElem{Name: "contacts", Value: filter.Contact})
	}
	if filter.Public {
		query = append(query, bson.DocElem{Name: "public", Value: filter.Public})
	}
	if filter.Disabled {
		query = append(query, bson.DocElem{Name: "disabled", Value: filter.Disabled})
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "clientMongoManager",
		Method:  "List",
		Query:   query,
	})
	defer span.Finish()

	var results []storage.Client
	collection := c.db.C(CollectionClients).With(mgoSession)
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

// Create stores a new OAuth2.0 Client resource.
func (c *clientMongoManager) Create(ctx context.Context, client storage.Client) (storage.Client, error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": CollectionClients,
		"method":     "Create",
		"id":         client.ID,
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = c.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Enable developers to provide their own IDs
	if client.ID == "" {
		client.ID = uuid.New()
	}
	if client.CreateTime == 0 {
		client.CreateTime = time.Now().Unix()
	}

	// Hash incoming secret
	hash, err := c.hasher.Hash([]byte(client.Secret))
	if err != nil {
		log.WithError(err).Error(logNotHashable)
		return client, err
	}
	client.Secret = string(hash)

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "clientMongoManager",
		Method:  "Create",
	})
	defer span.Finish()

	// Create resource
	collection := c.db.C(CollectionClients).With(mgoSession)
	err = collection.Insert(c)
	if err != nil {
		if mgo.IsDup(err) {
			// Log to StdOut
			log.WithError(err).Debug(logConflict)
			// Log to OpenTracing
			otLogErr(span, err)
			return client, storage.ErrClientExists
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		client.Secret = "REDACTED"
		otLogQuery(span, client)
		otLogErr(span, err)
		return client, err
	}
	return client, nil
}

// Get finds and returns an OAuth 2.0 client resource.
func (c *clientMongoManager) Get(ctx context.Context, clientID string) (storage.Client, error) {
	return c.getConcrete(ctx, clientID)
}

// GetClient finds and returns an OAuth 2.0 client resource.
//
// GetClient implements:
// - fosite.Storage
// - fosite.ClientManager
func (c *clientMongoManager) GetClient(ctx context.Context, clientID string) (fosite.Client, error) {
	client, err := c.getConcrete(ctx, clientID)
	return &client, err
}

// Update updates an OAuth 2.0 client resource.
func (c *clientMongoManager) Update(ctx context.Context, clientID string, updatedClient storage.Client) (storage.Client, error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": CollectionClients,
		"method":     "Update",
		"id":         clientID,
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = c.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	currentResource, err := c.getConcrete(ctx, clientID)
	if err != nil {
		if err == fosite.ErrNotFound {
			log.Debug(logNotFound)
			return currentResource, err
		}

		log.Error(logError)
		return currentResource, err
	}

	// Deny updating the entity Id
	updatedClient.ID = clientID
	// Update modified time
	updatedClient.UpdateTime = time.Now().Unix()

	if string(updatedClient.Secret) == "" {
		// If the password isn't updated, grab it from the stored object
		updatedClient.Secret = currentResource.Secret
	} else {
		newHash, err := c.hasher.Hash([]byte(updatedClient.Secret))
		if err != nil {
			log.WithError(err).Error(logNotHashable)
			return currentResource, err
		}
		updatedClient.Secret = string(newHash)
	}

	// Build Query
	selector := bson.M{
		"id": clientID,
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager:  "clientMongoManager",
		Method:   "Update",
		Selector: selector,
	})
	defer span.Finish()

	collection := c.db.C(CollectionClients).With(mgoSession)
	if err := collection.Update(selector, updatedClient); err != nil {
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
		otLogQuery(span, updatedClient)
		otLogErr(span, err)
		return currentResource, err
	}
	return updatedClient, nil
}

// Delete removes an OAuth 2.0 Client resource.
func (c *clientMongoManager) Delete(ctx context.Context, clientID string) error {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": CollectionClients,
		"method":     "Delete",
		"client":     clientID,
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
		"id": clientID,
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "clientMongoManager",
		Method:  "Delete",
		Query:   query,
	})
	defer span.Finish()

	collection := c.db.C(CollectionClients).With(mgoSession)
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

// Authenticate verifies the identity of a client resource.
func (c *clientMongoManager) Authenticate(ctx context.Context, clientID string, secret []byte) (storage.Client, error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": CollectionClients,
		"method":     "Authenticate",
		"client":     clientID,
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = c.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	client, err := c.getConcrete(ctx, clientID)
	if err != nil {
		log.WithError(err).Error(logError)
		return client, err
	}

	if client.Public {
		// The client doesn't have a secret, therefore is authenticated
		// implicitly.
		log.Debug("public client allowed access")
		return client, nil
	}

	if client.Disabled {
		log.Debug("disabled client denied access")
		return client, fosite.ErrAccessDenied
	}

	err = c.hasher.Compare(client.GetHashedSecret(), secret)
	if err != nil {
		log.WithError(err).Warn("failed to authenticate client secret")
		return client, err
	}

	return client, nil
}

func (c *clientMongoManager) AuthenticateMigration(ctx context.Context, currentAuth storage.AuthFunc, clientID string, secret []byte) (storage.Client, error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": CollectionClients,
		"method":     "AuthenticateMigration",
		"client":     clientID,
	})

	// Copy a new DB session if none specified
	mgoSession, ok := ContextToMgoSession(ctx)
	if !ok {
		mgoSession = c.db.Session.Copy()
		ctx = MgoSessionToContext(ctx, mgoSession)
		defer mgoSession.Close()
	}

	// Authenticate with old hasher
	client, authenticated := currentAuth()

	// Check for client not found
	if client.IsEmpty() && !authenticated {
		log.Debug(logNotFound)
		return client, fosite.ErrNotFound
	}

	if client.Public {
		// The client doesn't have a secret, therefore is authenticated
		// implicitly.
		log.Debug("public client allowed access")
		return client, nil
	}

	if client.Disabled {
		log.Debug("disabled client denied access")
		return client, fosite.ErrAccessDenied
	}

	if !authenticated {
		// If client isn't authenticated, try authenticating with new hasher.
		err := c.hasher.Compare(client.GetHashedSecret(), secret)
		if err != nil {
			log.WithError(err).Warn("failed to authenticate client secret")
		}
		return client, err
	}

	// If the client is found and authenticated, create a new hash using the new
	// hasher, update the database record and return the record with no error.
	newHash, err := c.hasher.Hash(secret)
	if err != nil {
		log.WithError(err).Error(logNotHashable)
		return client, err
	}

	// Save the new hash
	client.Secret = string(newHash)
	return c.Update(ctx, clientID, client)
}
