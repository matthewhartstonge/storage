package storage

import (
	"fmt"
	"github.com/MatthewHartstonge/storage/client"
	"github.com/ory/fosite"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"log"
	"strconv"
	"strings"
)

// Config provides a way to define the specific pieces that make up a mongo connection
type Config struct {
	// Default connection settings
	Hostname     string
	Hostnames    []string
	Port         uint16 // 0 to 65,535
	DatabaseName string

	// Credential Access
	Username string
	Password string

	// Replica Set
	Replset string
}

// DefaultConfig returns a configuration for a locally hosted, unauthenticated mongo
func DefaultConfig() *Config {
	return &Config{
		Hostname:     "127.0.0.1",
		Port:         27017,
		DatabaseName: "OAuth2",
	}
}

// ConnectionURI generates a formatted Mongo Connection URL
func ConnectionURI(cfg *Config) string {
	connectionString := "mongodb://"
	credentials := ""

	if cfg.Username != "" && cfg.Password != "" {
		credentials = fmt.Sprintf("%s:%s@", cfg.Username, cfg.Password)
	}

	hosts := ""
	if cfg.Hostnames != nil && cfg.Hostname == "" {
		hosts = strings.Join(cfg.Hostnames, fmt.Sprintf(":%s,", strconv.Itoa(int(cfg.Port))))
		cfg.Hostname = hosts
	}

	connectionString = fmt.Sprintf("%s%s%s:%s/%s",
		connectionString,
		credentials,
		cfg.Hostname,
		strconv.Itoa(int(cfg.Port)),
		cfg.DatabaseName,
	)

	if cfg.Replset != "" {
		connectionString += "?replicaSet=" + cfg.Replset
	}

	return connectionString
}

// MongoStore provides a fosite Datastore.
type MongoStore struct {
	Clients *client.MongoManager
	//AuthorizeCodes
	//IDSessions
	//AccessTokens
	//Implicit *session.Manager
	//RefreshTokens
	//AccessTokenRequestIDs
	//RefreshTokenRequestIDs
	//Users
}

// Close ensures that each endpoint has it's connection closed properly.
func (m *MongoStore) Close() {
	m.Clients.DB.Session.Close()
}

// ConnectToMongo returns a connection to mongo.
func ConnectToMongo(cfg *Config) (*mgo.Database, error) {
	uri := ConnectionURI(cfg)
	session, err := mgo.Dial(uri)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return session.DB(cfg.DatabaseName), nil
}

// NewDefaultMongoStore returns a MongoStore configured with the default mongo configuration and default hasher.
func NewDefaultMongoStore() *MongoStore {
	cfg := DefaultConfig()
	sess, err := ConnectToMongo(cfg)
	if err != nil {
		log.Fatalf("Could not connect to Mongo on %s:%s! Error:%s\n", cfg.Hostname, cfg.Port, err)
		return nil
	}
	h := &fosite.BCrypt{WorkFactor: 10}
	store := &MongoStore{
		Clients: &client.MongoManager{
			DB:     sess,
			Hasher: h,
		},
	}
	return store
}

// NewMongoStore allows for custom mongo configuration and custom hashers.
func NewMongoStore(c *Config, h fosite.Hasher) *MongoStore {
	sess, err := ConnectToMongo(c)
	if err != nil {
		log.Fatalf("Could not connect to Mongo on %s:%s! Error:%s\n", c.Hostname, c.Port, err)
		return nil
	}
	if h == nil {
		h = &fosite.BCrypt{WorkFactor: 10}
	}
	store := &MongoStore{
		Clients: &client.MongoManager{
			DB:     sess,
			Hasher: h,
		},
	}
	return store
}
