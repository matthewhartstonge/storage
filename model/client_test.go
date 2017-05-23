package model

import (
	"testing"

	"github.com/ory/fosite"
	"github.com/stretchr/testify/assert"
)

// TestClient ensures that Client conforms to fosite interfaces and that inputs and outputs are formed correctly.
func TestClient(t *testing.T) {
	c := &Client{
		ID:           "foo",
		RedirectURIs: []string{"foo"},
		Scope:        "foo bar",
	}

	assert.EqualValues(t, "foo", c.GetID())
	assert.EqualValues(t, c.RedirectURIs, c.GetRedirectURIs())
	assert.EqualValues(t, []byte(c.Secret), c.GetHashedSecret())
	assert.EqualValues(t, fosite.Arguments{"authorization_code"}, c.GetGrantTypes())
	assert.EqualValues(t, fosite.Arguments{"code"}, c.GetResponseTypes())
	assert.EqualValues(t, (c.Owner), c.GetOwner())
	assert.EqualValues(t, (c.Public), c.IsPublic())
	assert.Len(t, c.GetScopes(), 2)
	assert.EqualValues(t, c.RedirectURIs, c.GetRedirectURIs())

	// fosite.Argument logic branches
	expectedGrantTypes := fosite.Arguments{"foo", "bar"}
	c.GrantTypes = []string{"foo", "bar"}
	assert.EqualValues(t, expectedGrantTypes, c.GetGrantTypes())

	expectedResponseTypes := fosite.Arguments{"bar", "foo"}
	c.ResponseTypes = []string{"bar", "foo"}
	assert.EqualValues(t, expectedResponseTypes, c.GetResponseTypes())
}
