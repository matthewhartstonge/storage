package cache

type KeyValue struct {
	K string `bson:"_id" json:"key"`
	V string `bson:"value" json:"value"`
}
