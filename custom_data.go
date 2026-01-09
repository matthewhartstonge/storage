package storage

import "encoding/json"

type CustomData json.RawMessage

func (c *CustomData) Marshal(v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	*c = b
	return nil
}

func (c *CustomData) Unmarshal(v any) error {
	err := json.Unmarshal(*c, v)
	if err != nil {
		return err
	}
	return nil
}
