package request

import (
	"encoding/json"
	"github.com/MatthewHartstonge/storage/client"
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
	ID            string    `bson:"_id" json:"id"`
	RequestedAt   time.Time `bson:"requested_at" json:"requested_at"`
	Signature     string    `bson:"signature"`
	ClientID      string    `bson:"client_id" json:"client_id"`
	Scopes        []string  `bson:"scopes" json:"scopes"`
	GrantedScopes []string  `bson:"granted_scopes" json:"granted_scopes"`
	Form          string    `bson:"form_data" json:"form_data"`
	Session       []byte    `bson:"session_data"`
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
