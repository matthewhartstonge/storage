package storage

import "encoding/json"

// CustomData provides a way to binpack custom data into a mongo field. It
// utilizes [json.RawMessage] so data can be serialized or deserialized within
// application-level code to and from a Go type.
//
// Within mongo, the data is stored as base64 encoded binary data making it
// unindexable, un-queryable and unable to be understood at a glance from the
// record. But this does enable a user to extract custom data back out from
// mongo reliably instead of it being returned as a primitive.D on find
// operations.
type CustomData json.RawMessage

// Marshal serializes and stores v into itself as JSON.
func (c *CustomData) Marshal(v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	*c = b
	return nil
}

// Unmarshal deserializes the stored JSON into the value pointed to by v.
func (c *CustomData) Unmarshal(v any) error {
	err := json.Unmarshal(*c, v)
	if err != nil {
		return err
	}
	return nil
}
