package mongo

import (
	// Standard Library Imports
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/handler/pkce"

	// External Imports
	"github.com/ory/fosite"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

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
	DB *DB

	// timeout provides a way to configure maximum time before killing an
	// in-flight request.
	timeout time.Duration

	// Public API
	Hasher fosite.Hasher
	storage.Store
}

// DB wraps the mongo database connection and the features that are enabled.
type DB struct {
	*mongo.Database
	Hasher             fosite.Hasher
	RefreshGraceUsage  uint32
	RefreshGracePeriod time.Duration
}

// NewSession creates and returns a new mongo session.
// A deferrable session closer is returned in an attempt to enforce proper
// session handling/closing of sessions to avoid session and memory leaks.
//
// NewSession boilerplate becomes:
// ```
// ctx := context.Background()
//
//	if store.DB.HasSessions {
//	    var closeSession func()
//	    ctx, closeSession, err = store.NewSession(nil)
//	    if err != nil {
//	        panic(err)
//	    }
//	    defer closeSession()
//	}
//
// ```
func (s *Store) NewSession(ctx context.Context) (context.Context, func(), error) {
	return newSession(ctx, s.DB)
}

// newSession creates a new mongo session.
func newSession(ctx context.Context, db *DB) (context.Context, func(), error) {
	session, err := db.Client().StartSession()
	if err != nil {
		fields := logrus.Fields{
			"package": "mongo",
			"method":  "newSession",
		}
		logger.WithError(err).WithFields(fields).Error("error starting mongo session")
		return ctx, nil, err
	}

	if ctx == nil {
		ctx = context.Background()
	}
	ctx = SessionToContext(ctx, session)

	return ctx, closeSession(ctx, session), nil
}

// closeSession encapsulates the logic required to close a mongo session.
func closeSession(ctx context.Context, session mongo.Session) func() {
	return func() {
		session.EndSession(ctx)
	}
}

// Close terminates the mongo connection.
func (s *Store) Close() {
	err := s.DB.Client().Disconnect(context.Background())
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
	Hostnames    []string `default:"localhost" envconfig:"CONNECTIONS_MONGO_HOSTNAMES"`
	Port         uint16   `default:"27017"     envconfig:"CONNECTIONS_MONGO_PORT"`
	SSL          bool     `default:"false"     envconfig:"CONNECTIONS_MONGO_SSL"`
	AuthDB       string   `default:"admin"     envconfig:"CONNECTIONS_MONGO_AUTHDB"`
	Username     string   `default:""          envconfig:"CONNECTIONS_MONGO_USERNAME"`
	Password     string   `default:""          envconfig:"CONNECTIONS_MONGO_PASSWORD"`
	DatabaseName string   `default:""          envconfig:"CONNECTIONS_MONGO_NAME"`
	Replset      string   `default:""          envconfig:"CONNECTIONS_MONGO_REPLSET"`
	Timeout      uint     `default:"10"        envconfig:"CONNECTIONS_MONGO_TIMEOUT"`
	PoolMinSize  uint64   `default:"0"         envconfig:"CONNECTIONS_MONGO_POOL_MIN_SIZE"`
	PoolMaxSize  uint64   `default:"100"       envconfig:"CONNECTIONS_MONGO_POOL_MAX_SIZE"`
	Compressors  []string `default:""          envconfig:"CONNECTIONS_MONGO_COMPRESSORS"`
	TokenTTL     uint32   `default:"0"         envconfig:"CONNECTIONS_MONGO_TOKEN_TTL"`
	// RefreshTokenGracePeriod defines the refresh token 'slop' in seconds where a token can be reused.
	RefreshTokenGracePeriod uint32 `default:"0" envconfig:"CONNECTIONS_MONGO_REFRESH_TOKEN_GRACE_PERIOD"`
	// RefreshTokenMaxUsage limits the number of times a refresh token can be reused.
	RefreshTokenMaxUsage uint32      `default:"0" envconfig:"CONNECTIONS_MONGO_REFRESH_TOKEN_MAX_USAGE"`
	TLSConfig            *tls.Config `ignored:"true"`
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

	clientOpts := options.Client()
	if len(cfg.Hostnames) == 1 && strings.HasPrefix(cfg.Hostnames[0], "mongodb+srv://") {
		// MongoDB SRV records can only be configured with ApplyURI,
		// but we can continue to mung with client options after it's set.
		clientOpts.ApplyURI(cfg.Hostnames[0])
	} else {
		for i := range cfg.Hostnames {
			cfg.Hostnames[i] = fmt.Sprintf("%s:%d", cfg.Hostnames[i], cfg.Port)
		}
		clientOpts.SetHosts(cfg.Hostnames)
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = 10
	}

	clientOpts.SetReplicaSet(cfg.Replset).
		SetConnectTimeout(time.Second * time.Duration(cfg.Timeout)).
		SetReadPreference(readpref.SecondaryPreferred()).
		SetWriteConcern(writeconcern.Majority()).
		SetMinPoolSize(cfg.PoolMinSize).
		SetMaxPoolSize(cfg.PoolMaxSize).
		SetCompressors(cfg.Compressors).
		SetAppName(cfg.DatabaseName)

	if cfg.Username != "" || cfg.Password != "" {
		auth := options.Credential{
			AuthMechanism: "SCRAM-SHA-1",
			AuthSource:    cfg.AuthDB,
			Username:      cfg.Username,
			Password:      cfg.Password,
		}
		clientOpts.SetAuth(auth)
	}

	if cfg.SSL {
		tlsConfig := cfg.TLSConfig
		if tlsConfig == nil {
			// Inject a default TLS config if the SSL switch is toggled, but a
			// TLS config has not been provided programmatically.
			tlsConfig = &tls.Config{
				InsecureSkipVerify: false,
				MinVersion:         tls.VersionTLS12,
			}
		}

		clientOpts.SetTLSConfig(tlsConfig)
	}

	return clientOpts
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

	// Wrap database with mongo feature detection.
	mongoDB := &DB{
		Database:           database,
		Hasher:             hashee,
		RefreshGraceUsage:  cfg.RefreshTokenMaxUsage,
		RefreshGracePeriod: time.Second * time.Duration(cfg.RefreshTokenGracePeriod),
	}

	if hashee == nil {
		// Initialize default fosite Hasher.
		hashee = &fosite.BCrypt{
			Config: &fosite.Config{
				HashCost: fosite.DefaultBCryptWorkFactor,
			},
		}
	}

	// Build up the mongo endpoints
	mongoDeniedJtis := &DeniedJtiManager{
		DB: mongoDB,
	}
	mongoClients := &ClientManager{
		DB:     mongoDB,
		Hasher: hashee,

		DeniedJTIs: mongoDeniedJtis,
	}
	mongoUsers := &UserManager{
		DB:     mongoDB,
		Hasher: hashee,
	}
	mongoRequests := &RequestManager{
		DB: mongoDB,

		Clients: mongoClients,
		Users:   mongoUsers,
	}

	// attempt to perform index updates in a session.
	var closeSession func()
	ctx, closeSession, err := newSession(context.Background(), mongoDB)
	if err != nil {
		log.WithError(err).Error("error starting session")
		return nil, err
	}
	defer closeSession()

	// Configure DB collections, indices, TTLs e.t.c.
	if err = configureDatabases(ctx, mongoClients, mongoDeniedJtis, mongoUsers, mongoRequests); err != nil {
		log.WithError(err).Error("Unable to configure mongo collections!")
		return nil, err
	}
	if err = mongoRequests.ConfigureExpiryWithTTL(ctx, int(cfg.TokenTTL)); err != nil {
		log.WithError(err).Error("Failed to configure mongo auto expiry!")
		return nil, err
	}

	store := &Store{
		DB:      mongoDB,
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

// configureDatabases calls the configuration handler for the provided
// configurers.
func configureDatabases(ctx context.Context, configurers ...storage.Configurer) error {
	for _, configurer := range configurers {
		if err := configurer.Configure(ctx); err != nil {
			return err
		}
	}

	return nil
}

// configureExpiry calls the configuration handler for the provided expirers.
// ttl should be a positive integer.
func configureExpiry(ctx context.Context, ttl int, expirers ...storage.Expirer) error {
	for _, expirer := range expirers {
		if err := expirer.ConfigureExpiryWithTTL(ctx, ttl); err != nil {
			return err
		}
	}

	return nil
}

// NewDefaultStore returns a Store configured with the default mongo
// configuration and default Hasher.
func NewDefaultStore() (*Store, error) {
	cfg := DefaultConfig()
	return New(cfg, nil)
}

// NewIndex generates a new index model, ready to be saved in mongo.
//
// Note:
//   - This function assumes you are entering valid index keys and relies on
//     mongo rejecting index operations if a bad index is created.
func NewIndex(name string, keys ...string) (model mongo.IndexModel) {
	return mongo.IndexModel{
		Keys:    generateIndexKeys(keys...),
		Options: generateIndexOptions(name, false),
	}
}

// NewUniqueIndex generates a new unique index model, ready to be saved in
// mongo.
func NewUniqueIndex(name string, keys ...string) mongo.IndexModel {
	return mongo.IndexModel{
		Keys:    generateIndexKeys(keys...),
		Options: generateIndexOptions(name, true),
	}
}

// NewExpiryIndex generates a new index with a time to live value before the
// record expires in mongodb.
func NewExpiryIndex(name string, key string, expireAfter int) (model mongo.IndexModel) {
	return mongo.IndexModel{
		Keys: bson.D{{Key: key, Value: int32(1)}},
		Options: generateIndexOptions(name, false).
			SetExpireAfterSeconds(int32(expireAfter)),
	}
}

// generateIndexKeys given a number of stringy keys will return a bson
// document containing keys in the structure required by mongo for defining
// index and sort order.
func generateIndexKeys(keys ...string) (indexKeys bson.D) {
	var indexKey bson.E
	for _, key := range keys {
		switch {
		case strings.HasPrefix(key, "-"):
			// Reverse Index
			indexKey.Key = strings.TrimLeft(key, "-")
			indexKey.Value = int32(-1)

		case strings.HasPrefix(key, "#"):
			// Hashed Index
			indexKey.Key = strings.TrimLeft(key, "#")
			indexKey.Value = "hashed"

		default:
			// Forward Index
			indexKey.Key = key
			indexKey.Value = int32(1)
		}

		indexKeys = append(indexKeys, indexKey)
	}

	return
}

// generateIndexOptions generates new index options.
func generateIndexOptions(name string, unique bool) *options.IndexOptions {
	opts := options.Index().
		SetUnique(unique)

	if name != "" {
		opts.SetName(name)
	}

	return opts
}

var (
	_ fosite.Storage                                      = (*Store)(nil)
	_ fosite.ClientManager                                = (*Store)(nil)
	_ oauth2.CoreStorage                                  = (*Store)(nil)
	_ oauth2.AuthorizeCodeStorage                         = (*Store)(nil)
	_ oauth2.ClientCredentialsGrantStorage                = (*Store)(nil)
	_ oauth2.ResourceOwnerPasswordCredentialsGrantStorage = (*Store)(nil)
	_ oauth2.AccessTokenStorage                           = (*Store)(nil)
	_ oauth2.RefreshTokenStorage                          = (*Store)(nil)
	_ oauth2.TokenRevocationStorage                       = (*Store)(nil)
	_ openid.OpenIDConnectRequestStorage                  = (*Store)(nil)
	_ pkce.PKCERequestStorage                             = (*Store)(nil)
)

// TODO: implement me
// var _ fosite.PARStorage = (*Store)(nil)
// var _ fosite_storage.Transactional = (*Store)(nil)
// var _ rfc7523.RFC7523KeyStorage = (*Store)(nil)
// var _ verifiable.NonceManager = (*Store)(nil)
