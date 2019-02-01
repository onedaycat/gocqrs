package gocqrs

type EventBus interface {
	Publish(events []*EventMessage) error
}
