package storage

import (
	"fmt"
	"github.com/MatthewHartstonge/storage/client"
	"github.com/MatthewHartstonge/storage/request"
	"github.com/ory/fosite"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"strconv"
	"strings"
	"time"
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

	// Timeout specified in seconds.
	Timeout uint
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

	// OAuth Stores
	AuthorizeCodes *request.MongoManager
	IDSessions     *request.MongoManager
	Implicit       *request.MongoManager
	AccessTokens   *request.MongoManager
	RefreshTokens  *request.MongoManager

	// TODO: Create User Manager
	//Users *user.MongoManager

	// TODO: Create different cache storage backends?
	// - *cache.MemoryManager
	// - *cache.MongoManager
	// - *cache.RedisManager
	//AccessTokenRequestIDs *cache.MemoryManager
	//RefreshTokenRequestIDs *cache.MemoryManager
}

// Close ensures that each endpoint has it's connection closed properly.
func (m *MongoStore) Close() {
	m.Clients.DB.Session.Close()
}

// ConnectToMongo returns a connection to mongo.
func ConnectToMongo(cfg *Config) (*mgo.Database, error) {
	uri := ConnectionURI(cfg)
	if cfg.Timeout == 0 {
		cfg.Timeout = 10
	}
	session, err := mgo.DialWithTimeout(uri, time.Second*time.Duration(cfg.Timeout))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return session.DB(cfg.DatabaseName), nil
}

// NewDefaultMongoStore returns a MongoStore configured with the default mongo configuration and default hasher.
func NewDefaultMongoStore() (*MongoStore, error) {
	cfg := DefaultConfig()
	sess, err := ConnectToMongo(cfg)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	h := &fosite.BCrypt{WorkFactor: 10}
	c := &client.MongoManager{
		DB:     sess,
		Hasher: h,
	}
	r := &request.MongoManager{
		MongoManager: *c,
		DB:           sess,
	}
	return &MongoStore{
		Clients:        c,
		AuthorizeCodes: r,
		IDSessions:     r,
		Implicit:       r,
		AccessTokens:   r,
		RefreshTokens:  r,
	}, nil
}

// NewMongoStore allows for custom mongo configuration and custom hashers.
func NewMongoStore(cfg *Config, h fosite.Hasher) (*MongoStore, error) {
	sess, err := ConnectToMongo(cfg)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if h == nil {
		h = &fosite.BCrypt{WorkFactor: 10}
	}
	c := &client.MongoManager{
		DB:     sess,
		Hasher: h,
	}
	r := &request.MongoManager{
		MongoManager: *c,
		DB:           sess,
	}
	return &MongoStore{
		Clients:        c,
		AuthorizeCodes: r,
		IDSessions:     r,
		Implicit:       r,
		AccessTokens:   r,
		RefreshTokens:  r,
	}, nil
}
