package cache

type SessionObject struct {
	ID        string `bson:"_id" json:"key"`
	Signature string `bson:"signature" json:"signature"`
}

// GetKey returns the key of the cached session map
func (s SessionObject) GetKey() string {
	return s.ID
}

// GetValue returns session data as a string
func (s SessionObject) GetValue() (string, error) {
	return s.Signature, nil
}
