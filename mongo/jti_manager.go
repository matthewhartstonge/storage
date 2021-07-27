package mongo

import (
	// Standard Library Imports
	"context"
	"time"

	// External Imports
	"github.com/ory/fosite"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// Internal Imports
	"github.com/matthewhartstonge/storage"
)

// DeniedJtiManager provides a mongo backed implementation for denying JSON Web
// Tokens (JWTs) by ID.
type DeniedJtiManager struct {
	DB *DB
}

// Configure implements storage.Configurer.
func (d *DeniedJtiManager) Configure(ctx context.Context) (err error) {
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityJtiDenylist,
		"method":     "Configure",
	})

	indices := []mongo.IndexModel{
		{
			Keys: bson.D{
				{
					Key:   "signature",
					Value: int32(1),
				},
			},
			Options: options.Index().
				SetName(IdxSignatureID).
				SetSparse(true).
				SetUnique(true),
		},
		{
			Keys: bson.D{
				{
					Key:   "exp",
					Value: int32(1),
				},
			},
			Options: options.Index().
				SetName(IdxExpires).
				SetSparse(true),
		},
	}

	collection := d.DB.Collection(storage.EntityJtiDenylist)
	_, err = collection.Indexes().CreateMany(ctx, indices)
	if err != nil {
		log.WithError(err).Error(logError)
		return err
	}

	return nil
}

// getConcrete returns a denied jti resource.
func (d *DeniedJtiManager) getConcrete(ctx context.Context, signature string) (result storage.DeniedJTI, err error) {
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityJtiDenylist,
		"method":     "getConcrete",
		"signature":  signature,
	})

	// Build Query
	query := bson.M{
		"signature": signature,
	}

	// Trace how long the Mongo operation takes to complete.
	span, _ := traceMongoCall(ctx, dbTrace{
		Manager: "DeniedJtiManager",
		Method:  "getConcrete",
		Query:   query,
	})
	defer span.Finish()

	var user storage.DeniedJTI
	collection := d.DB.Collection(storage.EntityJtiDenylist)
	err = collection.FindOne(ctx, query).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.WithError(err).Debug(logNotFound)
			return result, fosite.ErrNotFound
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return result, err
	}

	return user, nil
}

// Create creates a new User resource and returns the newly created User
// resource.
func (d *DeniedJtiManager) Create(ctx context.Context, deniedJTI storage.DeniedJTI) (result storage.DeniedJTI, err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityJtiDenylist,
		"method":     "Create",
	})

	// Trace how long the Mongo operation takes to complete.
	span, _ := traceMongoCall(ctx, dbTrace{
		Manager: "DeniedJtiManager",
		Method:  "Create",
	})
	defer span.Finish()

	// Create resource
	collection := d.DB.Collection(storage.EntityJtiDenylist)
	_, err = collection.InsertOne(ctx, deniedJTI)
	if err != nil {
		if isDup(err) {
			// Log to StdOut
			log.WithError(err).Debug(logConflict)
			// Log to OpenTracing
			otLogErr(span, err)
			return result, storage.ErrResourceExists
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogQuery(span, deniedJTI)
		otLogErr(span, err)
		return result, err
	}

	return deniedJTI, nil
}

// Get returns the specified User resource.
func (d *DeniedJtiManager) Get(ctx context.Context, signature string) (result storage.DeniedJTI, err error) {
	return d.getConcrete(ctx, signature)
}

func (d *DeniedJtiManager) Delete(ctx context.Context, jti string) (err error) {
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityJtiDenylist,
		"method":     "Delete",
		"jti":        jti,
	})

	// Build Query
	query := bson.M{
		"signature": storage.SignatureFromJTI(jti),
	}

	// Trace how long the Mongo operation takes to complete.
	span, _ := traceMongoCall(ctx, dbTrace{
		Manager: "UserManager",
		Method:  "Delete",
		Query:   query,
	})
	defer span.Finish()

	collection := d.DB.Collection(storage.EntityJtiDenylist)
	res, err := collection.DeleteOne(ctx, query)
	if err != nil {
		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return err
	}

	if res.DeletedCount == 0 {
		// Log to StdOut
		log.WithError(err).Debug(logNotFound)
		// Log to OpenTracing
		otLogErr(span, err)
		return fosite.ErrNotFound
	}

	return nil
}

// DeleteExpired removes all JTIs before the given time. Returns not found if
// no tokens were found before the given time.
func (d *DeniedJtiManager) DeleteBefore(ctx context.Context, expBefore int64) (err error) {
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityJtiDenylist,
		"method":     "DeleteExpired",
		"expBefore":  expBefore,
	})

	// Build Query
	query := bson.M{
		"exp": bson.M{
			"$lt": time.Now().Unix(),
		},
	}

	// Trace how long the Mongo operation takes to complete.
	span, _ := traceMongoCall(ctx, dbTrace{
		Manager: "UserManager",
		Method:  "Delete",
		Query:   query,
	})
	defer span.Finish()

	collection := d.DB.Collection(storage.EntityJtiDenylist)
	res, err := collection.DeleteMany(ctx, query)
	if err != nil {
		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return err
	}

	if res.DeletedCount == 0 {
		// Log to StdOut
		log.WithError(err).Debug(logNotFound)
		// Log to OpenTracing
		otLogErr(span, err)
		return fosite.ErrNotFound
	}

	return nil
}
