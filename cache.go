package storage

// SessionCache allows storing a map between a session ID and a session signature
type SessionCache struct {
	// ID contains the unique identifier of the request.
	ID string `bson:"id" json:"id" xml:"id"`

	// createTime is when the resource was created in seconds from the epoch.
	CreateTime int64 `bson:"createTime" json:"createTime" xml:"createTime"`

	// updateTime is the last time the resource was modified in seconds from
	// the epoch.
	UpdateTime int64 `bson:"updateTime" json:"updateTime" xml:"updateTime"`

	// Signature contains the unique token signature.
	Signature string `bson:"signature" json:"signature" xml:"signature"`
}

// Key returns the key of the cached session map
func (s SessionCache) Key() string {
	return s.ID
}

// Value returns session data as a string
func (s SessionCache) Value() string {
	return s.Signature
}
