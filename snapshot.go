package gocqrs

type Snapshot struct {
	ID            string        `json:"id" bson:"_id"`
	AggregateType AggregateType `json:"b" bson:"b"`
	Version       int64         `json:"v" bson:"v"`
	Payload       *Payload      `json:"p" bson:"p"`
	LastUpdate    int64         `json:"t" bson:"t"`
	LastEvent     *EventMessage `json:"e" bson:"e"`
	IsRemoved     bool          `json:"r" bson:"r"`
}
