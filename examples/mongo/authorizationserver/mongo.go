package authorizationserver

import (
	"context"
	"time"

	"github.com/matthewhartstonge/storage"
	"github.com/matthewhartstonge/storage/mongo"
	log "github.com/sirupsen/logrus"
)

// init configures and starts an example mongo datastore, then
// returns a teardown function to clean up after itself.
func NewExampleMongoStore() *mongo.Store {
	ctx := context.Background()
	cfg := mongo.DefaultConfig()
	cfg.DatabaseName = "fositeStorageDemo"
	store, err := mongo.New(cfg, nil)
	if err != nil {
		// Make sure to check in on your mongo instance and drop the database
		// to ensure you can start this up again and not have conflicting data
		// attempted to be inserted.
		log.Warn("error configuring/starting up connection to mongo. please ensure you drop the oauth2 database locally if it exists..")
		log.WithError(err).Fatal("error creating new store")
	}

	// The general setup when working with the database is to create a session
	// which is a way to group a "logical" unit of work for mongo. Here, we
	// know we want to create a couple of clients and a user, therefore, we'll
	// group that into a server session, if we are using a mongo replica set.

	// If our mongo is running as a replica set we can use mongo sessions.
	// We luckily have `store.NewSession()` which does the hard work for us by
	// pushing the session into the context so all db handlers can use the same
	// connection/session and provides a function to be able to cleanly close
	// the session for us, which we can defer to later.
	ctx, closeSession, err := store.NewSession(ctx)
	if err != nil {
		// oh noes! creating a mongo session broke :/
		log.WithError(err).Fatal("error creating new session")
	}
	defer closeSession()

	// Inject our test clients
	clients := []storage.Client{
		{
			ID:               "my-client",
			Name:             "My Super Cool client for testing out Mongo storage",
			CreateTime:       time.Now().Unix(),
			Secret:           "foobar", // gets automagically hashed using fosite's hasher
			AllowedAudiences: []string{"https://my-client.my-application.com"},
			RedirectURIs:     []string{"http://localhost:3846/callback"},
			ResponseTypes:    []string{"id_token", "code", "token", "id_token token"},
			GrantTypes:       []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
			Scopes:           []string{"fosite", "openid", "photos", "offline"},
		},
		{
			ID:            "encoded:client",
			Name:          "Sup3r secret 3nc0d3d Client",
			CreateTime:    time.Now().Unix(),
			Secret:        "encoded&password", // gets automagically hashed using fosite's hasher
			RedirectURIs:  []string{"http://localhost:3846/callback"},
			ResponseTypes: []string{"id_token", "code", "token"},
			GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
			Scopes:        []string{"fosite", "openid", "photos", "offline"},
		},
	}
	createClients(ctx, store, clients)

	// Build and inject our test users
	users := []storage.User{
		{
			Username: "peter",
			Password: "secret",
		},
	}
	createUsers(ctx, store, users)

	return store
}

// TeardownMongo drops the database.
func TeardownMongo() {
	log.Info("dropping mongo database: oauth2")
	err := store.DB.Drop(nil)
	if err != nil {
		log.Error("error dropping oauth2 db:", err)
		return
	}
	log.Info("mongo database oauth2 dropped successfully!")
}

func createClients(ctx context.Context, store *mongo.Store, clients []storage.Client) {
	// Clean up after failed runs
	for _, client := range clients {
		logger := log.WithFields(log.Fields{
			"id":   client.ID,
			"name": client.Name,
		})

		// Attempt to remove any past remnant from bad builds/panics e.t.c.
		err := store.ClientManager.Delete(ctx, client.ID)
		if err == nil {
			logger.Debug("client found and deleted to enable clean start")
		}

		// Create the new client!
		_, err = store.ClientManager.Create(ctx, client)
		if err != nil {
			// err, it broke... ?
			panic(err)
		}
		logger.Info("new client created!")
	}
}

func createUsers(ctx context.Context, store *mongo.Store, users []storage.User) {
	for _, user := range users {
		logger := log.WithFields(log.Fields{
			"id":       user.ID,
			"username": user.Username,
		})

		// Attempt to remove any past remnant from bad builds/panics e.t.c.
		oldUser, err := store.UserManager.GetByUsername(ctx, user.Username)
		if err == nil {
			// yes, this could be done by setting an ID on the created user,
			// but here you can see how the storage handlers can work together
			err := store.UserManager.Delete(ctx, oldUser.ID)
			if err == nil {
				logger.Debug("client found and deleted to enable clean start")
			}
		}

		// Create the new user!
		user, err = store.UserManager.Create(ctx, user)
		if err != nil {
			// err, it broke... ?
			panic(err)
		}
		logger.WithField("id", user.ID).Info("new user created!")
	}
}
