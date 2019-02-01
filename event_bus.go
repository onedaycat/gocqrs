package gocqrs

type EventBus interface {
	Publish()
	Subscribe()
}
