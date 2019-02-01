package gocqrs

type AggregateType string

type AggregateRoot interface {
	Apply(payload *EventMessage) error
	GetAggregateID() string
	SetAggregateID(id string)
	GetAggregateType() AggregateType
	GetVersion() int
	SetVersion(version int)
	GetCurrentVersion() int
	IncreaseVersion()
	GetEvents() []Event
	ClearEvents()
	Publish(event Event)
	MarkAsRemoved()
	IsRemoved() bool
}

type AggregateBase struct {
	id      string
	events  []Event
	removed bool
	version int
}

func InitAggregate(id string) *AggregateBase {
	return &AggregateBase{
		id:      id,
		events:  make([]Event, 0, 1),
		version: 0,
	}
}

func (a *AggregateBase) GetAggregateID() string {
	return a.id
}

func (a *AggregateBase) SetAggregateID(id string) {
	a.id = id
}

func (a *AggregateBase) Publish(event Event) {
	a.events = append(a.events, event)
}

func (a *AggregateBase) GetEvents() []Event {
	return a.events
}

func (a *AggregateBase) MarkAsRemoved() {
	a.removed = true
}

func (a *AggregateBase) IsRemoved() bool {
	return a.removed
}

func (a *AggregateBase) GetVersion() int {
	return a.version
}

func (a *AggregateBase) SetVersion(version int) {
	a.version = version
}

func (a *AggregateBase) GetCurrentVersion() int {
	return a.version + len(a.events)
}

func (a *AggregateBase) ClearEvents() {
	a.events = make([]Event, 0, 1)
}

func (a *AggregateBase) IncreaseVersion() {
	a.version++
}
