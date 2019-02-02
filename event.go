package gocqrs

type EventType = string

type EventMessage struct {
	ID            string        `json:"id" bson:"_id"`
	AggregateID   string        `json:"a" bson:"a"`
	AggregateType AggregateType `json:"b" bson:"b"`
	Type          EventType     `json:"e" bson:"e"`
	Version       int64         `json:"v" bson:"v"`
	Payload       *Payload      `json:"p" bson:"p"`
	Time          int64         `json:"t" bson:"t"`
	Seq           int64         `json:"s" bson:"s"`
}

type Event interface {
	GetEventType() EventType
}
