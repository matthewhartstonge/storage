package client

import "github.com/ory/fosite"

// Manager provides a generic interface to clients in order to build a DataStore
type Manager interface {
	Storer

	Authenticate(id string, secret []byte) (*Client, error)
}

// Storage conforms to fosite.Storage and provides methods
type Storer interface {
	fosite.Storage

	GetConcreteClient(id string) (*Client, error)
	GetClients() (map[string]Client, error)
	CreateClient(c *Client) error
	UpdateClient(c *Client) error
	DeleteClient(id string) error
}
