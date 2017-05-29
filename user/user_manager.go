package user

// Manager provides a generic interface to users in order to build a DataStore
type Manager interface {
	Storer

	Authenticate(id string, secret []byte) (*User, error)
}

// Storer provides a definition of specific methods that are required to store a User in a data store.
type Storer interface {
	GetConcreteUser(id string) (*User, error)
	GetUser(id string) (User, error)
	GetUsers() (map[string]User, error)
	CreateUser(u *User) error
	UpdateUser(u *User) error
	DeleteUser(id string) error
	GrantScope(scope string) error
	RemoveScope(scope string) error
}
