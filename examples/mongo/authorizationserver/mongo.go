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
	store, err := mongo.NewDefaultStore()
	if err != nil {
		if store != nil {

		}
		// Make sure to check in on your mongo instance and drop the database
		// to ensure you can start this up again and not have conflicting data
		// attempted to be inserted.
		log.Warn("error configuring/starting up connection to mongo. please ensure you drop the oauth2 database locally if it exists..")
		log.WithError(err).Fatal("error creating new store")
	}

	// The general setup when working with the database is to create a session
	// which is a way to group a "logical" unit of work for mongo. Here, we
	// know we want to create a couple of clients and a user, therefore, we'll
	// group that into a session.

	// We luckily have `store.NewSession()` which does the hard work for us by
	// psuhing the session into the context so all db handlers can use the same
	// connection/session and provides a function to be able to cleanly close
	// the session for us, which we can defer to later.
	ctx, _, closeSession, err := store.NewSession(context.TODO())
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
	for _, client := range clients {
		newClient, err := store.ClientManager.Create(ctx, client)
		if err != nil {
			// err, it broke... ?
			panic(err)
		}

		fields := log.Fields{
			"id":   newClient.ID,
			"name": newClient.Name,
		}
		log.WithFields(fields).Info("new client created!")
	}
}

func createUsers(ctx context.Context, store *mongo.Store, users []storage.User) {
	for _, user := range users {
		newUser, err := store.UserManager.Create(ctx, user)
		if err != nil {
			// err, it broke... ?
			panic(err)
		}

		fields := log.Fields{
			"id":       newUser.ID,
			"username": newUser.Username,
		}
		log.WithFields(fields).Info("new user created!")
	}
}

func cleanUp() {

}
