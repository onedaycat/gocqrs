package gocqrs

type Snapshot struct {
	ID            string        `json:"id" bson:"_id"`
	AggregateType AggregateType `json:"b" bson:"b"`
	Version       int           `json:"v" bson:"v"`
	Payload       *Payload      `json:"p" bson:"p"`
	LastEvent     *EventMessage `json:"e" bson:"e"`
	IsRemoved     bool          `json:"r" bson:"r"`
}
