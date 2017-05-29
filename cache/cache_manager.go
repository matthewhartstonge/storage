package cache

// Manager provides a generic interface to key value cache objects in order to build a DataStore
type Manager interface {
	Storer
}

// Storer provides a way to create cache based objects in mongo
type Storer interface {
	Create(cacheObject KeyValue, collectionName string) error
	Get(key string, collectionName string) (*KeyValue, error)
	Update(cacheObject KeyValue, collectionName string) error
	Delete(key string, collectionName string) error
}
