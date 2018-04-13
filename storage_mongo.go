package storage

import (
	// Standard Library Imports
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	// External Imports
	"github.com/globalsign/mgo"
	"github.com/ory/fosite"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"

	// Internal Imports
	"github.com/MatthewHartstonge/storage/cache"
	"github.com/MatthewHartstonge/storage/client"
	"github.com/MatthewHartstonge/storage/request"
	"github.com/MatthewHartstonge/storage/user"
)

const (
	defaultHost = "localhost"
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

	// Timeout specified in seconds.
	Timeout uint

	// Mongo Options
	Replset   string
	SSL       bool
	TLSConfig *tls.Config
}

// DefaultConfig returns a configuration for a locally hosted, unauthenticated mongo
func DefaultConfig() *Config {
	return &Config{
		Hostname:     "localhost",
		Port:         27017,
		DatabaseName: "OAuth2",
	}
}

// ConnectionInfo configures options for establishing a session with a MongoDB cluster.
func ConnectionInfo(cfg *Config) *mgo.DialInfo {
	if len(cfg.Hostnames) == 0 && cfg.Hostname == "" {
		cfg.Hostname = defaultHost
	}

	if len(cfg.Hostnames) == 0 {
		cfg.Hostnames = []string{cfg.Hostname}
	}

	if cfg.DatabaseName == "" {
		cfg.DatabaseName = "oauth2"
	}

	for i := range cfg.Hostnames {
		cfg.Hostnames[i] = fmt.Sprintf("%s:%s", cfg.Hostnames[i], strconv.Itoa(int(cfg.Port)))
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = 10
	}

	dialInfo := &mgo.DialInfo{
		Addrs:          cfg.Hostnames,
		Database:       cfg.DatabaseName,
		Username:       cfg.Username,
		Password:       cfg.Password,
		ReplicaSetName: cfg.Replset,
		Timeout:        time.Second * time.Duration(cfg.Timeout),
	}

	if cfg.SSL {
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), cfg.TLSConfig)
		}
	}

	return dialInfo
}

// ConnectToMongo returns a connection to mongo.
func ConnectToMongo(cfg *Config) (*mgo.Database, error) {
	dialInfo := ConnectionInfo(cfg)
	session, err := mgo.DialWithInfo(dialInfo)
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
	hasher := &fosite.BCrypt{
		WorkFactor: 10,
	}
	mongoClients := &client.MongoManager{
		DB:     session,
		Hasher: hasher,
	}
	mongoUsers := &user.MongoManager{
		DB:     session,
		Hasher: hasher,
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
		DB:       session,
		Hasher:   hasher,
		Cache:    mongoCache,
		Clients:  mongoClients,
		Requests: mongoRequester,
		Users:    mongoUsers,
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
		DB:     session,
		Hasher: hasher,
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
		DB:       session,
		Hasher:   hasher,
		Cache:    mongoCache,
		Clients:  mongoClients,
		Requests: mongoRequester,
		Users:    mongoUsers,
	}, nil
}

// NewExampleMongoStore returns an example mongo store that matches the fosite-example data. If a default
// unauthenticated mongo database can't be found at localhost:27017, it will panic as you've done it wrong.
func NewExampleMongoStore() *MongoStore {
	m, err := NewDefaultMongoStore()
	if err != nil {
		panic(err)
	}
	err = m.Clients.CreateClient(&client.Client{
		ID:            "my-client",
		Secret:        []byte(`$2a$10$IxMdI6d.LIRZPpSfEwNoeu4rY3FhDREsxFJXikcgdRRAStxUlsuEO`), // = "foobar"
		RedirectURIs:  []string{"http://localhost:3846/callback"},
		ResponseTypes: []string{"id_token", "code", "token"},
		GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
		Scopes:        []string{"fosite", "openid", "photos", "offline"},
	})
	log.Printf("Mongo Client Creation error: %s", err)

	err = m.Users.CreateUser(&user.User{
		ID:         uuid.New(),
		Username:   "peter",
		Password:   "secret",
		FirstName:  "Peter",
		LastName:   "Secret",
		ProfileURI: "https://gravatar.com/avatar/e305b2c62b732cde23dbdd6f5b6ed6a9.png?s=256", // md5( peter@example.com )
	})
	log.Printf("Mongo User Creation error: %s", err)
	return m
}

// MongoStore composes all stores into the one datastore to rule them all
type MongoStore struct {
	// DB is the Mongo connection that holds the base session that can be copied and closed.
	DB       *mgo.Database
	Hasher   fosite.Hasher
	Clients  *client.MongoManager
	Requests *request.MongoManager
	Users    *user.MongoManager
	// Cache Stores
	// - *cache.MemoryManager
	// - *cache.MongoManager
	// - *cache.RedisManager
	Cache *cache.MongoManager
	//
	//AccessTokenRequestIDs  *cache.MongoManager
	//RefreshTokenRequestIDs *cache.MongoManager
}

// Close ensures that each endpoint has it's connection closed properly.
func (m *MongoStore) Close() {
	// As people can customise how they build up their mongo connections, ensure to close all endpoint individually.
	m.Clients.DB.Session.Close()
	if m.Requests != nil {
		m.Requests.DB.Session.Close()
	}
	if m.Users != nil {
		m.Users.DB.Session.Close()
	}
	if m.Cache != nil {
		m.Cache.DB.Session.Close()
	}
	// Kill top level session.
	m.DB.Session.Close()
}

/* Hoist all the funcs! */

// GetClient returns a Client if found by an ID lookup.
func (m MongoStore) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	return m.Clients.GetClient(ctx, id)
}

// GetClients returns a map of clients mapped by client ID
func (m MongoStore) GetClients() (clients map[string]client.Client, err error) {
	return m.Clients.GetClients()
}

// CreateClient adds a new OAuth2.0 Client to the client store.
func (m *MongoStore) CreateClient(c *client.Client) error {
	return m.Clients.CreateClient(c)
}

// UpdateClient updates an OAuth 2.0 Client record. This is done using the equivalent of an object replace.
func (m *MongoStore) UpdateClient(c *client.Client) error {
	return m.Clients.UpdateClient(c)
}

// DeleteClient removes an OAuth 2.0 Client from the client store
func (m *MongoStore) DeleteClient(id string) error {
	return m.Clients.DeleteClient(id)
}

// RevokeRefreshToken finds a token stored in cache based on request ID and deletes the session by signature.
func (m *MongoStore) RevokeRefreshToken(ctx context.Context, requestID string) error {
	return m.Requests.RevokeRefreshToken(ctx, requestID)
}

// RevokeAccessToken finds a token stored in cache based on request ID and deletes the session by signature.
func (m *MongoStore) RevokeAccessToken(ctx context.Context, requestID string) error {
	return m.Requests.RevokeAccessToken(ctx, requestID)
}

// CreateAccessTokenSession creates a new session for an Access Token in mongo
func (m *MongoStore) CreateAccessTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	return m.Requests.CreateAccessTokenSession(ctx, signature, request)
}

// GetAccessTokenSession returns a session if it can be found by signature in mongo
func (m MongoStore) GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	return m.Requests.GetAccessTokenSession(ctx, signature, session)
}

// DeleteAccessTokenSession removes an Access Tokens current session from mongo
func (m *MongoStore) DeleteAccessTokenSession(ctx context.Context, signature string) (err error) {
	return m.Requests.DeleteAccessTokenSession(ctx, signature)
}

// PersistAuthorizeCodeGrantSession creates an Authorise Code Grant session in mongo
func (m *MongoStore) PersistAuthorizeCodeGrantSession(ctx context.Context, authorizeCode, accessSignature, refreshSignature string, request fosite.Requester) error {
	if ts, err := m.GetRefreshTokenSession(ctx, authorizeCode, nil); err != nil {
		return err
	} else if err := m.RevokeAccessToken(ctx, ts.GetID()); err != nil {
		return err
	} else if err := m.RevokeRefreshToken(ctx, ts.GetID()); err != nil {
		return err
	} else if err := m.CreateAccessTokenSession(ctx, accessSignature, request); err != nil {
		return err
	} else if err := m.CreateRefreshTokenSession(ctx, refreshSignature, request); err != nil {
		return err
	}

	return nil
}

// CreateAuthorizeCodeSession creates a new session for an authorize code grant in mongo
func (m *MongoStore) CreateAuthorizeCodeSession(ctx context.Context, code string, request fosite.Requester) (err error) {
	return m.Requests.CreateAuthorizeCodeSession(ctx, code, request)
}

// GetAuthorizeCodeSession finds an authorize code grant session in mongo
func (m MongoStore) GetAuthorizeCodeSession(ctx context.Context, code string, session fosite.Session) (request fosite.Requester, err error) {
	return m.Requests.GetAuthorizeCodeSession(ctx, code, session)
}

// DeleteAuthorizeCodeSession removes an authorize code session from mongo
func (m *MongoStore) DeleteAuthorizeCodeSession(ctx context.Context, code string) (err error) {
	return m.Requests.DeleteAuthorizeCodeSession(ctx, code)
}

// CreateImplicitAccessTokenSession stores an implicit access token based session in mongo
func (m *MongoStore) CreateImplicitAccessTokenSession(ctx context.Context, token string, request fosite.Requester) (err error) {
	return m.Requests.CreateImplicitAccessTokenSession(ctx, token, request)
}

// PersistRefreshTokenGrantSession stores a refresh token grant session in mongo
func (m *MongoStore) PersistRefreshTokenGrantSession(ctx context.Context, originalRefreshSignature, accessSignature, refreshSignature string, request fosite.Requester) (err error) {
	if ts, err := m.GetRefreshTokenSession(ctx, originalRefreshSignature, nil); err != nil {
		return err
	} else if err := m.RevokeAccessToken(ctx, ts.GetID()); err != nil {
		return err
	} else if err := m.RevokeRefreshToken(ctx, ts.GetID()); err != nil {
		return err
	} else if err := m.CreateAccessTokenSession(ctx, accessSignature, request); err != nil {
		return err
	} else if err := m.CreateRefreshTokenSession(ctx, refreshSignature, request); err != nil {
		return err
	}

	return nil
}

// CreateRefreshTokenSession stores a new Refresh Token Session in mongo
func (m *MongoStore) CreateRefreshTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	return m.Requests.CreateRefreshTokenSession(ctx, signature, request)
}

// GetRefreshTokenSession returns a Refresh Token Session that's been previously stored in mongo
func (m *MongoStore) GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	return m.Requests.GetRefreshTokenSession(ctx, signature, session)
}

// DeleteRefreshTokenSession removes a Refresh Token that has been previously stored in mongo
func (m *MongoStore) DeleteRefreshTokenSession(ctx context.Context, signature string) (err error) {
	return m.Requests.DeleteRefreshTokenSession(ctx, signature)
}

// Authenticate checks if supplied user credentials are valid for User Credentials Grant
func (m *MongoStore) Authenticate(ctx context.Context, username string, secret string) error {
	return m.Requests.Authenticate(ctx, username, secret)
}

// AuthenticateUserByUsername checks if supplied user credentials are valid
func (m *MongoStore) AuthenticateUserByUsername(ctx context.Context, username string, secret string) (*user.User, error) {
	return m.Users.Authenticate(username, []byte(secret))
}

// AuthenticateClient checks is supplied client credentials are valid
func (m MongoStore) AuthenticateClient(id string, secret []byte) (*client.Client, error) {
	return m.Clients.Authenticate(id, secret)
}

// CreateOpenIDConnectSession creates an open id connect session for a given authorize code in mongo. This is relevant
// for explicit open id connect flow.
func (m *MongoStore) CreateOpenIDConnectSession(ctx context.Context, authorizeCode string, requester fosite.Requester) (err error) {
	return m.Requests.CreateOpenIDConnectSession(ctx, authorizeCode, requester)
}

// GetOpenIDConnectSession gets a session based off the Authorize Code and returns a fosite.Requester which contains a
// session or an error.
func (m *MongoStore) GetOpenIDConnectSession(ctx context.Context, authorizeCode string, requester fosite.Requester) (req fosite.Requester, err error) {
	return m.Requests.GetOpenIDConnectSession(ctx, authorizeCode, requester)
}

// DeleteOpenIDConnectSession removes an open id connect session from mongo.
func (m *MongoStore) DeleteOpenIDConnectSession(ctx context.Context, authorizeCode string) (err error) {
	return m.Requests.DeleteOpenIDConnectSession(ctx, authorizeCode)
}
