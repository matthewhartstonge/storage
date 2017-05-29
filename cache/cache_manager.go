package cache

// Manager provides a generic interface to key value cache objects in order to build a DataStore
type Manager interface {
	Storer
}

// Cacher provides a generic interface for storing data in a mongo cache
type Cacher interface {
	GetKey() string
	GetValue() (interface{}, error)
}

// Storer provides a way to create cache based objects in mongo
type Storer interface {
	Create(cacheObject Cacher, collectionName string) error
	Get(key string, collectionName string) (*Cacher, error)
	Update(cacheObject Cacher, collectionName string) error
	Delete(key string, collectionName string) error
}
