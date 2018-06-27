package mongo

import (
	// Standard Library Imports
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	// External Imports
	"github.com/globalsign/mgo"
	"github.com/ory/fosite"
	"github.com/sirupsen/logrus"

	// Local Imports
	"github.com/matthewhartstonge/storage"
)

func init() {
	// Bind a logger, but only to panic level. Leave it to the user to decide
	// whether they want datastore logging or not.
	SetLogger(logrus.New())
	logger.Level = logrus.PanicLevel
}

const (
	defaultHost         = "localhost"
	defaultPort         = 27017
	defaultDatabaseName = "oauth2"
)

type Store struct {
	// Internals
	DB *mgo.Database

	// Public API
	Hasher fosite.Hasher
	storage.Store
}

// NewSession returns a mongo session.
// Note: The session requires closing manually so no memory leaks occur.
// This is best achieved by calling `defer session.Close()` straight after
// obtaining the returned session object.
func (m *Store) NewSession() (session *mgo.Session) {
	return m.DB.Session.Copy()
}

// Close terminates the mongo session.
func (m *Store) Close() {
	m.DB.Session.Close()
}

// Config defines the configuration parameters which are used by GetMongoSession.
type Config struct {
	Hostnames    []string `default:"localhost" envconfig:"CONNECTIONS_MONGO_HOSTNAMES"`
	Port         uint16   `default:"27017"     envconfig:"CONNECTIONS_MONGO_PORT"`
	AuthDB       string   `default:"admin"     envconfig:"CONNECTIONS_MONGO_AUTHDB"`
	Username     string   `default:""          envconfig:"CONNECTIONS_MONGO_USERNAME"`
	Password     string   `default:""          envconfig:"CONNECTIONS_MONGO_PASSWORD"`
	DatabaseName string   `default:""          envconfig:"CONNECTIONS_MONGO_NAME"`
	Replset      string   `default:""          envconfig:"CONNECTIONS_MONGO_REPLSET"`
	Timeout      uint     `default:"10"        envconfig:"CONNECTIONS_MONGO_TIMEOUT"`
	SSL          bool     `default:"false"     envconfig:"CONNECTIONS_MONGO_SSL"`
	TLSConfig    *tls.Config
}

// DefaultConfig returns a configuration for a locally hosted, unauthenticated mongo
func DefaultConfig() *Config {
	return &Config{
		Hostnames:    []string{defaultHost},
		Port:         defaultPort,
		DatabaseName: defaultDatabaseName,
	}
}

// ConnectionInfo configures options for establishing a session with a MongoDB cluster.
func ConnectionInfo(cfg *Config) *mgo.DialInfo {
	if len(cfg.Hostnames) == 0 {
		cfg.Hostnames = []string{defaultHost}
	}

	if cfg.DatabaseName == "" {
		cfg.DatabaseName = defaultDatabaseName
	}

	if cfg.Port > 0 {
		for i := range cfg.Hostnames {
			cfg.Hostnames[i] = fmt.Sprintf("%s:%d", cfg.Hostnames[i], cfg.Port)
		}
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = 10
	}

	dialInfo := &mgo.DialInfo{
		Addrs:          cfg.Hostnames,
		Database:       cfg.DatabaseName,
		Username:       cfg.Username,
		Password:       cfg.Password,
		Source:         cfg.AuthDB,
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

// Connect returns a connection to mongo.
func Connect(cfg *Config) (*mgo.Database, error) {
	log := logger.WithFields(logrus.Fields{
		"package": "mongo",
		"method":  "Connect",
	})

	dialInfo := ConnectionInfo(cfg)
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.WithError(err).Error("Unable to connect to mongo! Have you configured your connection properly?")
		return nil, err
	}

	// Monotonic consistency will start reading from a slave if possible
	session.SetMode(mgo.Monotonic, true)
	return session.DB(cfg.DatabaseName), nil
}

// New allows for custom mongo configuration and custom hashers.
func New(cfg *Config, hashee fosite.Hasher) (*Store, error) {
	log := logger.WithFields(logrus.Fields{
		"package": "mongo",
		"method":  "NewFromConfig",
	})

	database, err := Connect(cfg)
	if err != nil {
		log.WithError(err).Error("Unable to connect to mongo! Are you sure mongo is running on localhost?")
		return nil, err
	}

	if hashee == nil {
		// Initialize default fosite Hasher.
		hashee = &fosite.BCrypt{
			WorkFactor: 10,
		}
	}

	// Build up the mongo endpoints
	mongoCache := &CacheManager{
		DB: database,
	}
	mongoClients := &ClientManager{
		DB:     database,
		Hasher: hashee,
	}
	mongoUsers := &UserManager{
		DB:     database,
		Hasher: hashee,
	}
	mongoRequests := &RequestManager{
		DB: database,

		Cache:   mongoCache,
		Clients: mongoClients,
		Users:   mongoUsers,
	}

	// Init DB collections, indices e.t.c.
	managers := []storage.Configurer{
		mongoCache,
		mongoClients,
		mongoUsers,
		mongoRequests,
	}

	// Use the main DB database to configure the mongo collections on first up.
	mgoSession := database.Session.Copy()
	defer mgoSession.Close()
	ctx := MgoSessionToContext(context.Background(), mgoSession)

	for _, manager := range managers {
		err := manager.Configure(ctx)
		if err != nil {
			log.WithError(err).Error("Unable to configure mongo collections!")
			return nil, err
		}
	}

	return &Store{
		DB:     database,
		Hasher: hashee,
		Store: storage.Store{
			CacheManager:   mongoCache,
			ClientManager:  mongoClients,
			RequestManager: mongoRequests,
			UserManager:    mongoUsers,
		},
	}, nil
}

// NewDefaultStore returns a Store configured with the default mongo
// configuration and default Hasher.
func NewDefaultStore() (*Store, error) {
	cfg := DefaultConfig()
	return New(cfg, nil)
}
