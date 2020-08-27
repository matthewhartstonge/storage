package mongo

import (
	// Standard Library Imports
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	// External Imports
	"github.com/ory/fosite"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

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

// Store provides a MongoDB storage driver compatible with fosite's required
// storage interfaces.
type Store struct {
	// Internals
	DB *mongo.Database
	// timeout provides a way to configure maximum time before killing an
	// in-flight request.
	timeout time.Duration

	// Public API
	Hasher fosite.Hasher
	storage.Store
}

// NewSession creates and returns a new mongo session.
// A deferrable session closer is returned in an attempt to enforce closing of
// sessions.
//
// NewSession boilerplate becomes:
// ```
// ctx, session, sessionCloser, err := store.NewSession(nil)
// if err != nil {
// 	panic(err)
// }
// defer sessionCloser()
// ```
func (s *Store) NewSession(ctx context.Context) (context.Context, mongo.Session, func(), error) {
	return newSession(ctx, s.DB)
}

// newSession creates a new mongo session.
func newSession(ctx context.Context, db *mongo.Database) (context.Context, mongo.Session, func(), error) {
	session, err := db.Client().StartSession()
	if err != nil {
		fields := logrus.Fields{
			"package": "mongo",
			"method":  "newSession",
		}
		logger.WithError(err).WithFields(fields).Error("error starting mongo session")
		return ctx, nil, nil, err
	}

	if ctx == nil {
		ctx = context.Background()
	}
	ctx = SessionToContext(ctx, session)

	return ctx, session, sessionCloser(ctx, session), nil
}

// sessionCloser encapsulates the logic required to close a mongo session.
func sessionCloser(ctx context.Context, session mongo.Session) func() {
	return func() {
		session.EndSession(ctx)
	}
}

// Close terminates the mongo connection.
func (s *Store) Close() {
	err := s.DB.Client().Disconnect(nil)
	if err != nil {
		fields := logrus.Fields{
			"package": "mongo",
			"method":  "Close",
		}
		logger.WithError(err).WithFields(fields).Error("error closing mongo connection")
	}
}

// Config defines the configuration parameters which are used by GetMongoSession.
type Config struct {
	Hostnames    []string    `default:"localhost" envconfig:"CONNECTIONS_MONGO_HOSTNAMES"`
	Port         uint16      `default:"27017"     envconfig:"CONNECTIONS_MONGO_PORT"`
	SSL          bool        `default:"false"     envconfig:"CONNECTIONS_MONGO_SSL"`
	AuthDB       string      `default:"admin"     envconfig:"CONNECTIONS_MONGO_AUTHDB"`
	Username     string      `default:""          envconfig:"CONNECTIONS_MONGO_USERNAME"`
	Password     string      `default:""          envconfig:"CONNECTIONS_MONGO_PASSWORD"`
	DatabaseName string      `default:""          envconfig:"CONNECTIONS_MONGO_NAME"`
	Replset      string      `default:""          envconfig:"CONNECTIONS_MONGO_REPLSET"`
	Timeout      uint        `default:"10"        envconfig:"CONNECTIONS_MONGO_TIMEOUT"`
	PoolMinSize  uint64      `default:"0"         envconfig:"CONNECTIONS_MONGO_POOL_MIN_SIZE"`
	PoolMaxSize  uint64      `default:"100"       envconfig:"CONNECTIONS_MONGO_POOL_MAX_SIZE"`
	TLSConfig    *tls.Config `ignored:"true"`
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
func ConnectionInfo(cfg *Config) *options.ClientOptions {
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

	dialInfo := options.Client().
		SetHosts(cfg.Hostnames).
		SetReplicaSet(cfg.Replset).
		SetConnectTimeout(time.Second * time.Duration(cfg.Timeout)).
		SetReadPreference(readpref.SecondaryPreferred()).
		SetMinPoolSize(cfg.PoolMinSize).
		SetMaxPoolSize(cfg.PoolMaxSize)

	if cfg.Username != "" || cfg.Password != "" {
		auth := options.Credential{
			AuthMechanism: "SCRAM-SHA-1",
			AuthSource:    cfg.AuthDB,
			Username:      cfg.Username,
			Password:      cfg.Password,
		}
		dialInfo.SetAuth(auth)
	}

	if cfg.SSL {
		dialInfo = dialInfo.SetTLSConfig(cfg.TLSConfig)
	}

	return dialInfo
}

// Connect returns a connection to a mongo database.
func Connect(cfg *Config) (*mongo.Database, error) {
	log := logger.WithFields(logrus.Fields{
		"package": "mongo",
		"method":  "Connect",
	})

	ctx := context.Background()
	dialInfo := ConnectionInfo(cfg)
	client, err := mongo.Connect(ctx, dialInfo)
	if err != nil {
		log.WithError(err).Error("Unable to build mongo connection!")
		return nil, err
	}

	// check connection works as mongo-go lazily connects.
	err = client.Ping(ctx, nil)
	if err != nil {
		log.WithError(err).Error("Unable to connect to mongo! Have you configured your connection properly?")
		return nil, err
	}

	return client.Database(cfg.DatabaseName), nil
}

// New allows for custom mongo configuration and custom hashers.
func New(cfg *Config, hashee fosite.Hasher) (*Store, error) {
	log := logger.WithFields(logrus.Fields{
		"package": "mongo",
		"method":  "NewFromConfig",
	})

	database, err := Connect(cfg)
	if err != nil {
		log.WithError(err).Error("Unable to connect to mongo! Are you sure mongo is running?")
		return nil, err
	}

	if hashee == nil {
		// Initialize default fosite Hasher.
		hashee = &fosite.BCrypt{
			WorkFactor: 10,
		}
	}

	// Build up the mongo endpoints
	mongoDeniedJtis := &DeniedJtiManager{
		DB: database,
	}
	mongoClients := &ClientManager{
		DB:     database,
		Hasher: hashee,

		DeniedJTIs: mongoDeniedJtis,
	}
	mongoUsers := &UserManager{
		DB:     database,
		Hasher: hashee,
	}
	mongoRequests := &RequestManager{
		DB: database,

		Clients: mongoClients,
		Users:   mongoUsers,
	}

	// Init DB collections, indices e.t.c.
	managers := []storage.Configurer{
		mongoClients,
		mongoDeniedJtis,
		mongoUsers,
		mongoRequests,
	}

	// Use the main DB database to configure the mongo collections on first up.
	ctx, _, closer, err := newSession(nil, database)
	if err != nil {
		log.WithError(err).Error("error starting session")
		return nil, err
	}
	defer closer()

	for _, manager := range managers {
		err := manager.Configure(ctx)
		if err != nil {
			log.WithError(err).Error("Unable to configure mongo collections!")
			return nil, err
		}
	}

	store := &Store{
		DB:      database,
		timeout: time.Second * time.Duration(cfg.Timeout),
		Hasher:  hashee,
		Store: storage.Store{
			ClientManager:    mongoClients,
			DeniedJTIManager: mongoDeniedJtis,
			RequestManager:   mongoRequests,
			UserManager:      mongoUsers,
		},
	}
	return store, nil
}

// NewDefaultStore returns a Store configured with the default mongo
// configuration and default Hasher.
func NewDefaultStore() (*Store, error) {
	cfg := DefaultConfig()
	return New(cfg, nil)
}

const (
	// errCodeDuplicate provides the mongo error code for duplicate key error.
	errCodeDuplicate = 11000
)

// isDup replicates mgo.IsDup functionality for the official driver in order
// to know when a conflict has occurred.
func isDup(err error) (isDup bool) {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == errCodeDuplicate {
				return true
			}
		}
	}

	return
}
