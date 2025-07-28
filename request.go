package storage

import (
	// Standard Library Imports
	"context"
	"encoding/json"
	"net/url"
	"time"

	// External Imports
	"github.com/google/uuid"
	"github.com/ory/fosite"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Request is a concrete implementation of a fosite.Requester, extended to
// support the required data for OAuth2 and OpenID.
type Request struct {
	// ID contains the request identifier, which is effectively a session id.
	// This leads to multiple requests being able to exist for the same session
	// due to graceful token rotation.
	ID string `bson:"id" json:"id" xml:"id"`
	// CreateTime is when the resource was created in seconds from the epoch.
	CreateTime int64 `bson:"createTime" json:"createTime" xml:"createTime"`
	// UpdateTime is the last time the resource was modified in seconds from
	// the epoch.
	UpdateTime int64 `bson:"updateTime" json:"updateTime" xml:"updateTime"`
	// RequestedAt is the time the request was made.
	RequestedAt time.Time `bson:"requestedAt" json:"requestedAt" xml:"requestedAt"`
	// FirstUsedAt sets when the refresh token was first used.
	FirstUsedAt time.Time `bson:"firstUsedAt" json:"firstUsedAt" xml:"firstUsedAt"`
	// ExpiresAt is the time the refresh token reuse expires.
	ExpiresAt time.Time `bson:"expiresAt" json:"expiresAt" xml:"expiresAt"`
	// UsageCount specifies how many times a refresh token has been used.
	UsageCount uint32 `bson:"usageCount" json:"usageCount" xml:"usageCount"`
	// Signature contains a unique session signature.
	// The Signature denotes the unique access token in use throughout the lifetime of a 'request' - think session.
	Signature string `bson:"signature" json:"signature" xml:"signature"`
	// AccessSignature specifies a refresh token's linked access token.
	AccessSignature string `bson:"accessSignature" json:"accessSignature" xml:"accessTokenSignature"`
	// ClientID contains a link to the Client that was used to authenticate
	// this session.
	ClientID string `bson:"clientId" json:"clientId" xml:"clientId"`
	// UserID contains the subject's unique ID which links back to a stored
	// user account.
	UserID string `bson:"userId" json:"userId" xml:"userId"`
	// Scopes contains the scopes that the user requested.
	RequestedScope fosite.Arguments `bson:"scopes" json:"scopes" xml:"scopes"`
	// GrantedScope contains the list of scopes that the user was actually
	// granted.
	GrantedScope fosite.Arguments `bson:"grantedScopes" json:"grantedScopes" xml:"grantedScopes"`
	// RequestedAudience contains the audience the user requested.
	RequestedAudience fosite.Arguments `bson:"requestedAudience" json:"requestedAudience" xml:"requestedAudience"`
	// GrantedAudience contains the list of audiences the user was actually
	// granted.
	GrantedAudience fosite.Arguments `bson:"grantedAudience" json:"grantedAudience" xml:"grantedAudience"`
	// Form contains the url values that were passed in to authenticate the
	// user's client session.
	Form url.Values `bson:"formData" json:"formData" xml:"formData"`
	// Active is specifically used for Authorize Code flow revocation.
	Active bool `bson:"active" json:"active" xml:"active"`
	// Session contains the session data. The underlying structure differs
	// based on OAuth strategy, so we need to store it as binary-encoded JSON.
	// Otherwise, it can be stored but not unmarshalled back into a
	// fosite.Session.
	Session []byte `bson:"sessionData" json:"sessionData" xml:"sessionData"`
}

func (r *Request) WithinGracePeriod(gracePeriod time.Duration) bool {
	return gracePeriod > 0 && r.FirstUsedAt.Add(gracePeriod).After(time.Now().UTC())
}

func (r *Request) WithinGraceUsage(graceUsage uint32) bool {
	return graceUsage == 0 || // no limit
		(r.UsageCount < graceUsage)
}

// NewRequest returns a new Mongo Store request object.
func NewRequest() Request {
	return Request{
		ID:                uuid.NewString(),
		CreateTime:        0,
		UpdateTime:        0,
		RequestedAt:       time.Now().UTC(),
		FirstUsedAt:       time.Time{},
		ExpiresAt:         time.Time{},
		UsageCount:        0,
		Signature:         "",
		AccessSignature:   "",
		ClientID:          "",
		UserID:            "",
		RequestedScope:    fosite.Arguments{},
		GrantedScope:      fosite.Arguments{},
		RequestedAudience: fosite.Arguments{},
		GrantedAudience:   fosite.Arguments{},
		Form:              make(url.Values),
		Active:            true,
		Session:           nil,
	}
}

// ToRequest transforms a mongo request to a fosite.Request
func (r *Request) ToRequest(ctx context.Context, session fosite.Session, cm ClientStorer) (*fosite.Request, error) {
	if session != nil {
		if err := json.Unmarshal(r.Session, session); err != nil {
			return nil, errors.WithStack(err)
		}
	} else {
		log.Debug("Got an empty session in toRequest")
	}

	client, err := cm.GetClient(ctx, r.ClientID)
	if err != nil {
		return nil, err
	}

	req := &fosite.Request{
		ID:                r.ID,
		RequestedAt:       r.RequestedAt,
		Client:            client,
		RequestedScope:    r.RequestedScope,
		GrantedScope:      r.GrantedScope,
		Form:              r.Form,
		Session:           session,
		RequestedAudience: r.RequestedAudience,
		GrantedAudience:   r.GrantedAudience,
	}
	return req, nil
}
