package request

import (
	"github.com/MatthewHartstonge/storage/client"
	"net/url"
	"time"
)

// Request is a concrete implementation of a fosite.Requester, extended to support the required data for OAuth2 and
// OpenID.
type Request struct {
	ID            string        `bson:"_id" json:"id"`
	RequestedAt   time.Time     `bson:"requested_at" json:"requested_at"`
	Client        client.Client `bson:"client" json:"client"`
	Scopes        []string      `bson:"scopes" json:"scopes"`
	GrantedScopes []string      `bson:"granted_scopes" json:"granted_scopes"`
	Form          url.Values    `bson:"form" json:"form"`

	// Potenially required extra fields
	// Signature
}
