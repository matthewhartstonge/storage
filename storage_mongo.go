package storage

import (
	"fmt"
	"github.com/MatthewHartstonge/storage/cache"
	"github.com/MatthewHartstonge/storage/client"
	"github.com/MatthewHartstonge/storage/request"
	"github.com/MatthewHartstonge/storage/user"
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
		Hostname:     "localhost",
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
	// OAuth Stores
	Clients        *client.MongoManager
	AuthorizeCodes *request.MongoManager
	IDSessions     *request.MongoManager
	Implicit       *request.MongoManager
	AccessTokens   *request.MongoManager
	RefreshTokens  *request.MongoManager

	// User Store
	Users *user.MongoManager

	// Cache Stores
	// - *cache.MemoryManager
	// - *cache.MongoManager
	// - *cache.RedisManager
	AccessTokenRequestIDs  *cache.MongoManager
	RefreshTokenRequestIDs *cache.MongoManager
}

// Close ensures that each endpoint has it's connection closed properly.
func (m *MongoStore) Close() {
	// As people can customise how they build up their mongo connections, ensure to close all endpoint individually.
	m.Clients.DB.Session.Close()
	if m.AuthorizeCodes != nil {
		m.AuthorizeCodes.DB.Session.Close()
	}
	if m.IDSessions != nil {
		m.IDSessions.DB.Session.Close()
	}
	if m.Implicit != nil {
		m.Implicit.DB.Session.Close()
	}
	if m.AccessTokens != nil {
		m.AccessTokens.DB.Session.Close()
	}
	if m.RefreshTokens != nil {
		m.RefreshTokens.DB.Session.Close()
	}
	if m.Users != nil {
		m.Users.DB.Session.Close()
	}
	if m.AccessTokenRequestIDs != nil {
		m.AccessTokenRequestIDs.DB.Session.Close()
	}
	if m.RefreshTokenRequestIDs != nil {
		m.RefreshTokenRequestIDs.DB.Session.Close()
	}
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

	// Monotonic consistency will start reading from a slave if possible
	session.SetMode(mgo.Monotonic, true)
	return session.DB(cfg.DatabaseName), nil
}

// NewDefaultMongoStore returns a MongoStore configured with the default mongo configuration and default hasher.
func NewDefaultMongoStore() (*MongoStore, error) {
	cfg := DefaultConfig()
	session, err := ConnectToMongo(cfg)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	hasher := &fosite.BCrypt{WorkFactor: 10}
	mongoClients := &client.MongoManager{
		DB:     session,
		Hasher: hasher,
	}
	mongoUsers := &user.MongoManager{
		DB: session,
	}
	mongoCache := &cache.MongoManager{
		DB: session,
	}
	mongoRequester := &request.MongoManager{
		DB:      session,
		Cache:   mongoCache,
		Clients: mongoClients,
		Users:   mongoUsers,
	}
	return &MongoStore{
		Clients:                mongoClients,
		Users:                  mongoUsers,
		AuthorizeCodes:         mongoRequester,
		IDSessions:             mongoRequester,
		Implicit:               mongoRequester,
		AccessTokens:           mongoRequester,
		RefreshTokens:          mongoRequester,
		AccessTokenRequestIDs:  mongoCache,
		RefreshTokenRequestIDs: mongoCache,
	}, nil
}

// NewMongoStore allows for custom mongo configuration and custom hashers.
func NewMongoStore(cfg *Config, hasher fosite.Hasher) (*MongoStore, error) {
	session, err := ConnectToMongo(cfg)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if hasher == nil {
		hasher = &fosite.BCrypt{WorkFactor: 10}
	}
	mongoClients := &client.MongoManager{
		DB:     session,
		Hasher: hasher,
	}
	mongoUsers := &user.MongoManager{
		DB: session,
	}
	mongoCache := &cache.MongoManager{
		DB: session,
	}
	mongoRequester := &request.MongoManager{
		DB:      session,
		Cache:   mongoCache,
		Clients: mongoClients,
		Users:   mongoUsers,
	}
	return &MongoStore{
		Clients:                mongoClients,
		Users:                  mongoUsers,
		AuthorizeCodes:         mongoRequester,
		IDSessions:             mongoRequester,
		Implicit:               mongoRequester,
		AccessTokens:           mongoRequester,
		RefreshTokens:          mongoRequester,
		AccessTokenRequestIDs:  mongoCache,
		RefreshTokenRequestIDs: mongoCache,
	}, nil
}
