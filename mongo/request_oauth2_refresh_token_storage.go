package mongo

import (
	// Standard Library Imports
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	// External Imports
	"github.com/ory/fosite"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

	// Internal Imports
	"github.com/matthewhartstonge/storage"
)

// CreateRefreshTokenSession implements fosite.RefreshTokenStorage.
func (r *RequestManager) CreateRefreshTokenSession(ctx context.Context, refreshSignature string, accessSignature string, request fosite.Requester) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityRefreshTokens,
		"method":     "CreateRefreshTokenSession",
	})

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "CreateRefreshTokenSession",
	})
	defer span.Finish()

	// Store session request
	exp := request.GetSession().GetExpiresAt(fosite.RefreshToken)
	entity, err := refreshToMongo(refreshSignature, storage.SignatureHash(accessSignature), request, exp)
	if err != nil {
		log.WithError(err).Error(logError)
		return err
	}

	_, err = r.Create(ctx, storage.EntityRefreshTokens, entity)
	if err != nil {
		if err == storage.ErrResourceExists {
			log.WithError(err).Debug(logConflict)
			return err
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		return err
	}

	return nil
}

// GetRefreshTokenSession implements fosite.RefreshTokenStorage.
func (r *RequestManager) GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityRefreshTokens,
		"method":     "GetRefreshTokenSession",
	})

	// Copy a new DB session if none specified
	_, ok := ContextToSession(ctx)
	if !ok {
		var closeSession func()
		ctx, closeSession, err = newSession(ctx, r.DB)
		if err != nil {
			log.WithError(err).Debug("error starting session")
			return nil, err
		}
		defer closeSession()
	}

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "GetRefreshTokenSession",
	})
	defer span.Finish()

	// Get the stored request
	req, err := r.GetBySignature(ctx, storage.EntityRefreshTokens, signature)
	if err != nil {
		switch {
		case errors.Is(err, fosite.ErrNotFound):
			log.WithError(err).Debug(logNotFound)
			return nil, err

		default:
			// Log to StdOut
			log.WithError(err).Error(logError)
			return nil, err
		}
	}

	if req.Active {
		// Token is active
		return req.ToRequest(ctx, session, r.Clients)
	}

	// If not active, perform graceful rotation checks to see if it's still within the specified 'grace' usage.
	if req.WithinGracePeriod(r.DB.RefreshGracePeriod) && req.WithinGraceUsage(r.DB.RefreshGraceUsage) {
		// We return the request as is, which indicates that the token is active (because we are in the grace period still).
		return req.ToRequest(ctx, session, r.Clients)
	}

	// Transform to a fosite.Request
	request, err = req.ToRequest(ctx, session, r.Clients)
	if err != nil {
		// Log to StdOut
		log.WithError(err).Error(logError)
		return nil, err
	}

	return request, fosite.ErrInactiveToken
}

// RotateRefreshToken performs a refresh token rotation.
// It handles both graceful and strict rotation modes.
// https://github.com/ory/hydra/blob/eda94bc9e984f68b36a65f823f528ae4f50e76af/persistence/sql/persister_oauth2.go#L770-L780
func (r *RequestManager) RotateRefreshToken(ctx context.Context, requestID string, refreshTokenSignature string) (err error) {
	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "RotateRefreshToken",
	})
	defer span.Finish()

	if r.DB.RefreshGracePeriod > 0 {
		// If we end up here, we have a valid refresh token and can proceed with a graceful rotation.
		return r.gracefulRefreshRotation(ctx, requestID, refreshTokenSignature)
	}

	// In strict rotation we only have one token chain for every request. Therefore, we remove all
	// access tokens associated with the request ID.
	return r.strictRefreshRotation(ctx, requestID)
}

// gracefulRefreshRotation updates the refresh token to reflect the expiration
// of its grace period.
// https://github.com/ory/hydra/blob/eda94bc9e984f68b36a65f823f528ae4f50e76af/persistence/sql/persister_oauth2.go#L683-L768
func (r *RequestManager) gracefulRefreshRotation(ctx context.Context, requestID string, refreshSignature string) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package": "mongo",
		"method":  "gracefulRefreshRotation",
	})

	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "gracefulRefreshRotation",
	})
	defer span.Finish()

	now := time.Now().UTC().Round(time.Millisecond)
	// The new expiry of the token starts now and ends at the end of the graceful token period.
	// After that, we can prune tokens from the store.
	expiresAt := newUsedExpiry().Add(r.DB.RefreshGracePeriod)

	// Signature is the primary key so no limit needed. We only update first_used_at if it is not set yet (otherwise
	// we would "refresh" the grace period again and again, and the refresh token would never "expire").
	filter := bson.M{
		"signature": refreshSignature,
	}
	if r.DB.RefreshGraceUsage > 0 {
		// check for maximum reuse count
		filter["$or"] = []bson.M{
			{"usageCount": bson.M{"$exists": false}},
			{"usageCount": bson.M{"$lt": r.DB.RefreshGraceUsage}},
		}
	}

	updatePipeline := mongo.Pipeline{
		{{
			Key: "$set", Value: bson.M{
				"active": false,
				"firstUsedAt": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							// Check if firstUsedAt does not exist or is null
							"$or": bson.A{
								bson.M{"$eq": bson.A{"$firstUsedAt", nil}},                  // If it's explicitly null
								bson.M{"$not": bson.A{bson.M{"$isNumber": "$firstUsedAt"}}}, // More robust check if it's not a valid date
							},
						},
						"then": now,            // If not set, set to current time
						"else": "$firstUsedAt", // Otherwise, keep its existing value
					},
				},
				"expiresAt":  expiresAt,
				"updateTime": now.Unix(),
				"usageCount": bson.M{
					"$add": bson.A{"$usageCount", 1},
				},
			},
		}},
	}

	var request storage.Request
	err = r.DB.Collection(storage.EntityRefreshTokens).FindOneAndUpdate(ctx, filter, updatePipeline).Decode(&request)
	if err != nil {
		switch {
		case errors.Is(err, fosite.ErrNotFound):
			// Tokens may have been pruned earlier, so we do not return an error here.
			log.WithError(err).Error(logNotFound)
			return nil

		default:
			log.WithError(err).Error(logError)
			return err
		}
	}

	if request.AccessSignature == "" {
		// If the access token is not found, we fall back to deleting all access tokens associated with the request ID.
		if err = r.DeleteAll(ctx, storage.EntityAccessTokens, requestID); err != nil {
			switch {
			case errors.Is(err, fosite.ErrNotFound):
				// Tokens may have been pruned earlier, so we do not return an error here.
				return nil

			default:
				log.WithError(err).Error(logError)
				return err
			}
		}
	}

	// We have the signature and we will only remove that specific access token as part of the rotation.
	if err = r.DeleteBySignature(ctx, storage.EntityAccessTokens, storage.SignatureHash(request.AccessSignature)); err != nil {
		switch {
		case errors.Is(err, fosite.ErrNotFound):
			// Tokens may have been pruned earlier, so we do not return an error here.
			return nil

		default:
			log.WithError(err).Error(logError)
			return err
		}
	}

	return nil
}

// strictRefreshRotation implements the strict refresh token rotation strategy. In strict rotation, we disable all
// refresh and access tokens associated with a request ID and subsequently create the only valid, new token pair.
// https://github.com/ory/hydra/blob/eda94bc9e984f68b36a65f823f528ae4f50e76af/persistence/sql/persister_oauth2.go#L646-L681
func (r *RequestManager) strictRefreshRotation(ctx context.Context, requestID string) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package": "mongo",
		"method":  "strictRefreshRotation",
	})

	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "strictRefreshRotation",
	})
	defer span.Finish()

	// In strict rotation we only have one token chain for every request. Therefore, we remove all
	// access tokens associated with the request ID.
	if err = r.DeleteAll(ctx, storage.EntityAccessTokens, requestID); err != nil {
		switch {
		case errors.Is(err, fosite.ErrNotFound):
			// Tokens may have been pruned earlier, so we do not return an error here.
			log.WithError(err).Debug(logNotFound)
			return nil

		default:
			log.WithError(err).Error(logError)
			return err
		}
	}

	// The same applies to refresh tokens in strict mode. We disable all old refresh tokens when rotating.
	filter := bson.M{
		"id":     requestID,
		"active": true,
	}
	update := bson.M{
		"$set": bson.M{
			"active": false,
			// We don't expire immediately, but in 30 minutes to avoid prematurely removing
			// rows while they may still be needed (e.g. for reuse detection).
			"expiresAt":  newUsedExpiry(),
			"updateTime": time.Now().UTC().Round(time.Millisecond).Unix(),
		},
	}
	res, err := r.DB.Collection(storage.EntityRefreshTokens).UpdateMany(ctx, filter, update)
	if err != nil {
		// Log to StdOut
		log.WithError(err).Error(logError)
		// Log to OpenTracing
		otLogErr(span, err)
		return err
	}

	if res.MatchedCount == 0 {
		// Log to StdOut
		log.WithError(err).Debug(logNotFound)
		// Log to OpenTracing
		otLogErr(span, err)
		return fosite.ErrNotFound
	}

	return nil
}

// DeleteRefreshTokenSession implements fosite.RefreshTokenStorage.
func (r *RequestManager) DeleteRefreshTokenSession(ctx context.Context, signature string) (err error) {
	// Initialize contextual method logger
	log := logger.WithFields(logrus.Fields{
		"package":    "mongo",
		"collection": storage.EntityRefreshTokens,
		"method":     "DeleteRefreshTokenSession",
	})

	// Trace how long the Mongo operation takes to complete.
	span, ctx := traceMongoCall(ctx, dbTrace{
		Manager: "RequestManager",
		Method:  "DeleteRefreshTokenSession",
	})
	defer span.Finish()

	// Remove session request
	err = r.DeleteBySignature(ctx, storage.EntityRefreshTokens, signature)
	if err != nil {
		if err == fosite.ErrNotFound {
			log.WithError(err).Debug(logNotFound)
			return err
		}

		// Log to StdOut
		log.WithError(err).Error(logError)
		return err
	}

	return nil
}

func newUsedExpiry() time.Time {
	// Reuse detection is racy and would generally happen within seconds. Using 30 minutes here is a paranoid
	// setting but ensures that we do not prematurely remove rows while they may still be needed (e.g. for reuse detection).
	return time.Now().UTC().Round(time.Millisecond).Add(time.Minute * 30)
}
