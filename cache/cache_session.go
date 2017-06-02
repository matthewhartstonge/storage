package cache

// SessionCache allows storing a map between a session ID and a session signature
type SessionCache struct {
	ID        string `bson:"_id" json:"id" xml:"id"`
	Signature string `bson:"signature" json:"signature" xml:"signature"`
}

// GetKey returns the key of the cached session map
func (s SessionCache) GetKey() string {
	return s.ID
}

// GetValue returns session data as a string
func (s SessionCache) GetValue() string {
	return s.Signature
}
