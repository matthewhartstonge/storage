package storage

import (
	"github.com/ory/fosite"
)

// Client provides the structure of an OAuth2.0 Client.
type Client struct {
	//// Client Meta
	// ID is the id for this client.
	ID string `bson:"id" json:"id" xml:"id"`

	// createTime is when the resource was created in seconds from the epoch.
	CreateTime int64 `bson:"createTime" json:"createTime" xml:"createTime"`

	// updateTime is the last time the resource was modified in seconds from
	// the epoch.
	UpdateTime int64 `bson:"updateTime" json:"updateTime" xml:"updateTime"`

	// AllowedAudiences contains a list of Audiences that the client has been
	// given rights to access.
	AllowedAudiences []string `bson:"allowedAudiences" json:"allowedAudiences,omitempty" xml:"allowedAudiences,omitempty"`

	// AllowedTenantAccess contains a list of Tenants that the client has been
	// given rights to access.
	AllowedTenantAccess []string `bson:"allowedTenantAccess" json:"allowedTenantAccess,omitempty" xml:"allowedTenantAccess,omitempty"`

	// GrantTypes contains a list of grant types the client is allowed to use.
	//
	// Pattern: client_credentials|authorize_code|implicit|refresh_token
	GrantTypes []string `bson:"grantTypes" json:"grantTypes" xml:"grantTypes"`

	// ResponseTypes contains a list of the OAuth 2.0 response type strings
	// that the client can use at the authorization endpoint.
	//
	// Pattern: id_token|code|token
	ResponseTypes []string `bson:"responseTypes" json:"responseTypes" xml:"responseTypes"`

	// Scopes contains a list of values the client is entitled to use when
	// requesting an access token (as described in Section 3.3 of OAuth 2.0
	// [RFC6749]).
	//
	// Pattern: ([a-zA-Z0-9\.]+\s)+
	Scopes []string `bson:"scopes" json:"scopes" xml:"scopes"`

	// Public is a boolean that identifies this client as public, meaning that
	// it does not have a secret. It will disable the client_credentials grant
	// type for this client if set.
	Public bool `bson:"public" json:"public" xml:"public"`

	// Disabled stops the client from being able to authenticate to the system.
	Disabled bool `bson:"disabled" json:"disabled" xml:"disabled"`

	//// Client Content
	// Name contains a human-readable string name of the client to be presented
	// to the end-user during authorization.
	Name string `bson:"name" json:"name" xml:"name"`

	// Secret is the client's secret. The secret will be included in the create
	// request as cleartext, and then never again. The secret is stored using
	// BCrypt so it is impossible to recover it.
	// Tell your users that they need to remember the client secret as it will
	// not be made available again.
	Secret string `bson:"secret,omitempty" json:"secret,omitempty" xml:"secret,omitempty"`

	// RedirectURIs contains a list of allowed redirect urls for the client, for
	// example: http://mydomain/oauth/callback.
	RedirectURIs []string `bson:"redirectUris" json:"redirectUris" xml:"redirectUris"`

	// Owner identifies the owner of the OAuth 2.0 Client.
	Owner string `bson:"owner" json:"owner" xml:"owner"`

	// PolicyURI allows the application developer to provide a URI string that
	// points to a human-readable privacy policy document that describes how the
	// deployment organization collects, uses, retains, and discloses personal
	// data.
	PolicyURI string `bson:"policyUri" json:"policyUri" xml:"policyUri"`

	// TermsOfServiceURI allows the application developer to provide a URI
	// string that points to a human-readable terms of service document that
	// describes and outlines the contractual relationship between the end-user
	// and the client application that the end-user accepts when authorizing
	// their use of the client.
	TermsOfServiceURI string `bson:"termsOfServiceUri" json:"termsOfServiceUri" xml:"termsOfServiceUri"`

	// ClientURI allows the application developer to provide a URI string that
	// points to a human-readable web page that provides information about the
	// client application.
	// If present, the server SHOULD display this URL to the end-user in a
	// click-able fashion.
	ClientURI string `bson:"clientUri" json:"clientUri" xml:"clientUri"`

	// LogoURI is an URL string that references a logo for the client.
	LogoURI string `bson:"logoUri" json:"logoUri" xml:"logoUri"`

	// Contacts contains a list ways to contact the developers responsible for
	// this OAuth 2.0 client, typically email addresses.
	Contacts []string `bson:"contacts" json:"contacts" xml:"contacts"`
}

// GetID returns the client's Client ID.
func (c *Client) GetID() string {
	return c.ID
}

// GetRedirectURIs returns the OAuth2.0 authorized Client redirect URIs.
func (c *Client) GetRedirectURIs() []string {
	return c.RedirectURIs
}

// GetHashedSecret returns the Client's Hashed Secret for authenticating with
// the Identity Provider.
func (c *Client) GetHashedSecret() []byte {
	return []byte(c.Secret)
}

// GetScopes returns an array of strings, wrapped as `fosite.Arguments` to
// provide functions that allow verifying the
// Client's scopes against incoming requests.
func (c *Client) GetScopes() fosite.Arguments {
	return c.Scopes
}

// GetGrantTypes returns an array of strings, wrapped as `fosite.Arguments` to
// provide functions that allow verifying
// the Client's Grant Types against incoming requests.
func (c *Client) GetGrantTypes() fosite.Arguments {
	// https://openid.net/specs/openid-connect-registration-1_0.html#ClientMetadata
	//
	// JSON array containing a list of the OAuth 2.0 Grant Types that the Client
	// is declaring that it will restrict itself to using.
	// If omitted, the default is that the Client will use only the
	// authorization_code Grant Type.
	if len(c.GrantTypes) == 0 {
		return fosite.Arguments{"authorization_code"}
	}
	return c.GrantTypes
}

// GetResponseTypes returns an array of strings, wrapped as `fosite.Arguments`
// to provide functions that allow verifying
// the Client's Response Types against incoming requests.
func (c *Client) GetResponseTypes() fosite.Arguments {
	// https://openid.net/specs/openid-connect-registration-1_0.html#ClientMetadata
	//
	// <JSON array containing a list of the OAuth 2.0 response_type values that
	// the Client is declaring that it will restrict itself to using. If
	// omitted, the default is that the Client will use only the code Response
	// Type.
	if len(c.ResponseTypes) == 0 {
		return fosite.Arguments{"code"}
	}
	return c.ResponseTypes
}

// GetOwner returns a string which contains the OAuth Client owner's name.
// Generally speaking, this will be a developer or an organisation.
func (c *Client) GetOwner() string {
	return c.Owner
}

// IsPublic returns a boolean as to whether the Client itself is either private
// or public. If public, only trusted OAuth grant types should be used as
// client secrets shouldn't be exposed to a public client.
func (c *Client) IsPublic() bool {
	return c.Public
}

// GetAudience returns the allowed audience(s) for this client.
func (c *Client) GetAudience() fosite.Arguments {
	return c.AllowedAudiences
}

// IsDisabled returns a boolean as to whether the Client itself has had it's
// access disabled.
func (c *Client) IsDisabled() bool {
	return c.Disabled
}

// EnableScopeAccess enables client scope access
func (c *Client) EnableScopeAccess(scopes ...string) {
	for i := range scopes {
		found := false
		for j := range c.Scopes {
			if scopes[i] == c.Scopes[j] {
				found = true
				break
			}
		}
		if !found {
			c.Scopes = append(c.Scopes, scopes[i])
		}
	}
}

// DisableScopeAccess disables client scope access.
func (c *Client) DisableScopeAccess(scopes ...string) {
	for i := range scopes {
		for j := range c.Scopes {
			if scopes[i] == c.Scopes[j] {
				copy(c.Scopes[j:], c.Scopes[j+1:])
				c.Scopes[len(c.Scopes)-1] = ""
				c.Scopes = c.Scopes[:len(c.Scopes)-1]
				break
			}
		}
	}
}

// EnableTenantAccess adds a single or multiple tenantIDs to the given client.
func (c *Client) EnableTenantAccess(tenantIDs ...string) {
	for i := range tenantIDs {
		found := false
		for j := range c.AllowedTenantAccess {
			if tenantIDs[i] == c.AllowedTenantAccess[j] {
				found = true
				break
			}
		}
		if !found {
			c.AllowedTenantAccess = append(c.AllowedTenantAccess, tenantIDs[i])
		}
	}
}

// DisableTenantAccess removes a single or multiple tenantIDs from the given
// client.
func (c *Client) DisableTenantAccess(tenantIDs ...string) {
	for i := range tenantIDs {
		for j := range c.AllowedTenantAccess {
			if tenantIDs[i] == c.AllowedTenantAccess[j] {
				copy(c.AllowedTenantAccess[j:], c.AllowedTenantAccess[j+1:])
				c.AllowedTenantAccess[len(c.AllowedTenantAccess)-1] = ""
				c.AllowedTenantAccess = c.AllowedTenantAccess[:len(c.AllowedTenantAccess)-1]
				break
			}
		}
	}
}

// Equal enables checking equality as having a byte array in a struct stops
// allowing equality checks.
func (c Client) Equal(x Client) bool {
	if c.ID != x.ID {
		return false
	}

	if c.CreateTime != x.CreateTime {
		return false
	}

	if c.UpdateTime != x.UpdateTime {
		return false
	}

	if !stringArrayEquals(c.AllowedAudiences, x.AllowedAudiences) {
		return false
	}

	if !stringArrayEquals(c.AllowedTenantAccess, x.AllowedTenantAccess) {
		return false
	}

	if !stringArrayEquals(c.GrantTypes, x.GrantTypes) {
		return false
	}

	if !stringArrayEquals(c.ResponseTypes, x.ResponseTypes) {
		return false
	}

	if !stringArrayEquals(c.Scopes, x.Scopes) {
		return false
	}

	if c.Public != x.Public {
		return false
	}

	if c.Disabled != x.Disabled {
		return false
	}

	if c.Name != x.Name {
		return false
	}

	if c.Secret != x.Secret {
		return false
	}

	if !stringArrayEquals(c.RedirectURIs, x.RedirectURIs) {
		return false
	}

	if c.Owner != x.Owner {
		return false
	}

	if c.PolicyURI != x.PolicyURI {
		return false
	}

	if c.TermsOfServiceURI != x.TermsOfServiceURI {
		return false
	}

	if c.ClientURI != x.ClientURI {
		return false
	}

	if c.LogoURI != x.LogoURI {
		return false
	}

	if !stringArrayEquals(c.Contacts, x.Contacts) {
		return false
	}

	return true
}

// IsEmpty returns whether or not the client resource is an empty record.
func (c Client) IsEmpty() bool {
	return c.Equal(Client{})
}
