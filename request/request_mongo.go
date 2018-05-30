package request

import (
	"encoding/json"
	"github.com/matthewhartstonge/storage/client"
	"github.com/ory/fosite"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"log"
	"net/url"
	"time"
)

// MongoRequest is a concrete implementation of a fosite.Requester, extended to support the required data for
// OAuth2 and OpenID.
type MongoRequest struct {
	ID            string    `bson:"_id" json:"id" xml:"id"`
	RequestedAt   time.Time `bson:"requestedAt" json:"requestedAt" xml:"requestedAt"`
	Signature     string    `bson:"signature" json:"signature" xml:"signature"`
	ClientID      string    `bson:"clientId" json:"clientId" xml:"clientId"`
	Scopes        []string  `bson:"scopes" json:"scopes" xml:"scopes"`
	GrantedScopes []string  `bson:"grantedScopes" json:"grantedScopes" xml:"grantedScopes"`
	Form          string    `bson:"formData" json:"formData" xml:"formData"`
	Session       []byte    `bson:"sessionData" json:"sessionData" xml:"sessionData"`
}

func NewRequest() *MongoRequest {
	return &MongoRequest{
		ID:            uuid.New(),
		RequestedAt:   time.Now(),
		Scopes:        []string{},
		GrantedScopes: []string{},
		Form:          "",
		Session:       []byte(""),
	}
}

// toRequest transforms a mongo database reference to a fosite request
func (m *MongoRequest) toRequest(session fosite.Session, cm client.Manager) (*fosite.Request, error) {
	if session != nil {
		if err := json.Unmarshal(m.Session, session); err != nil {
			return nil, errors.WithStack(err)
		}
	} else {
		log.Println("Got an empty session in toRequest")
	}

	c, err := cm.GetClient(nil, m.ClientID)
	if err != nil {
		return nil, err
	}

	val, err := url.ParseQuery(m.Form)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	r := &fosite.Request{
		ID:            m.ID,
		RequestedAt:   m.RequestedAt,
		Client:        c,
		Scopes:        m.Scopes,
		GrantedScopes: m.GrantedScopes,
		Form:          val,
		Session:       session,
	}

	return r, nil
}
