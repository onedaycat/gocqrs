package gocqrs

type EventType = string

type EventMessage struct {
	EventID       string        `json:"i" bson:"_id"`
	PartitionKey  string        `json:"k" bson:"k"`
	AggregateID   string        `json:"a" bson:"a"`
	AggregateType AggregateType `json:"b" bson:"b"`
	EventType     EventType     `json:"e" bson:"e"`
	Payload       *Payload      `json:"p" bson:"p"`
	Time          int64         `json:"t" bson:"t"`
	Seq           int64         `json:"s" bson:"s"`
	TimeSeq       int64         `json:"x" bson:"x"`
	Metadata      *Payload      `json:"m" bson:"m"`
}
