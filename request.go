package storage

import (
	// Standard Library Imports
	"context"
	"net/url"
	"time"

	// External Imports
	"github.com/ory/fosite"
	"github.com/pborman/uuid"
)

// MongoRequest is a concrete implementation of a fosite.Requester, extended to support the required data for
// OAuth2 and OpenID.
type Request struct {
	// ID contains the unique request identifier.
	ID string `bson:"id" json:"id" xml:"id"`
	// CreateTime is when the resource was created in seconds from the epoch.
	CreateTime int64 `bson:"createTime" json:"createTime" xml:"createTime"`
	// UpdateTime is the last time the resource was modified in seconds from
	// the epoch.
	UpdateTime int64 `bson:"updateTime" json:"updateTime" xml:"updateTime"`
	// RequestedAt is the time the request was made.
	RequestedAt time.Time `bson:"requestedAt" json:"requestedAt" xml:"requestedAt"`
	// Signature contains a unique session signature.
	Signature string `bson:"signature" json:"signature" xml:"signature"`
	// ClientID contains a link to the Client that was used to authenticate
	// this session.
	ClientID string `bson:"clientId" json:"clientId" xml:"clientId"`
	// UserID contains the subject's unique ID which links back to a stored
	// user account.
	UserID string `bson:"userId" json:"userId" xml:"userId"`
	// Scopes contains the scopes that the user requested.
	Scopes fosite.Arguments `bson:"scopes" json:"scopes" xml:"scopes"`
	// GrantedScopes contains the list of scopes that the user was actually
	// granted.
	GrantedScopes fosite.Arguments `bson:"grantedScopes" json:"grantedScopes" xml:"grantedScopes"`
	// Form contains the url values that were passed in to authenticate the
	// user's client session.
	Form url.Values `bson:"formData" json:"formData" xml:"formData"`
	// Active is specifically used for Authorize Code flow revocation.
	Active bool `bson:"active" json:"active" xml:"active"`
	// Session contains the session data. The underlying structure differs
	// based on OAuth strategy, but thanks to Mongo magic, we can magically
	// store an arbitary structure as long as it can marshal to json.
	Session fosite.Session `bson:"sessionData" json:"sessionData" xml:"sessionData"`
}

// NewRequest returns a new Mongo Store request object.
func NewRequest() Request {
	return Request{
		ID:            uuid.New(),
		RequestedAt:   time.Now(),
		Signature:     "",
		ClientID:      "",
		UserID:        "",
		Scopes:        []string{},
		GrantedScopes: []string{},
		Form:          make(url.Values),
		Active:        true,
		Session:       nil,
	}
}

// ToRequest transforms a mongo request to a fosite.Request
func (r *Request) ToRequest(ctx context.Context, session fosite.Session, cm ClientStorer) (*fosite.Request, error) {
	client, err := cm.GetClient(ctx, r.ClientID)
	if err != nil {
		return nil, err
	}

	req := &fosite.Request{
		ID:            r.ID,
		RequestedAt:   r.RequestedAt,
		Client:        client,
		Scopes:        r.Scopes,
		GrantedScopes: r.GrantedScopes,
		Form:          r.Form,
		Session:       r.Session,
	}
	return req, nil
}
