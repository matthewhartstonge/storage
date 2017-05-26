package user

type Manager interface {
	Storer

	Authenticate(id string, secret []byte) (*User, error)
}

// Storage conforms to fosite.Storage and provides methods
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
