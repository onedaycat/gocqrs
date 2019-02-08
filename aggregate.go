package gocqrs

import (
	"github.com/onedaycat/gocqrs/common/eid"
)

type AggregateType = string

type AggregateRoot interface {
	Apply(payload *EventMessage) error
	GetAggregateID() string
	SetAggregateID(id string) *AggregateBase
	GetAggregateType() string
	GetVersion() int64
	SetVersion(version int64) *AggregateBase
	GetCurrentVersion() int64
	GetLastUpdate() int64
	SetLastUpdate(t int64) *AggregateBase
	IncreaseVersion()
	GetEvents() []Event
	GetEventTypes() []EventType
	ClearEvents()
	Publish(eventType EventType, event Event)
	MarkAsRemoved()
	IsRemoved() bool
	IsNew() bool
}

type AggregateBase struct {
	id         string
	events     []Event
	eventTypes []EventType
	removed    bool
	version    int64
	lastUpdate int64
}

// InitAggregate if id is empty, id will be generated
func InitAggregate() *AggregateBase {
	return &AggregateBase{
		id:         eid.GenerateID(),
		events:     make([]Event, 0, 1),
		eventTypes: make([]EventType, 0, 1),
		version:    0,
	}
}

func (a *AggregateBase) GetAggregateID() string {
	return a.id
}

func (a *AggregateBase) SetAggregateID(id string) *AggregateBase {
	a.id = id
	return a
}

func (a *AggregateBase) Publish(eventType EventType, event Event) {
	a.events = append(a.events, event)
	a.eventTypes = append(a.eventTypes, eventType)
}

func (a *AggregateBase) GetEvents() []Event {
	return a.events
}

func (a *AggregateBase) GetEventTypes() []EventType {
	return a.eventTypes
}

func (a *AggregateBase) MarkAsRemoved() {
	a.removed = true
}

func (a *AggregateBase) IsRemoved() bool {
	return a.removed
}

func (a *AggregateBase) GetVersion() int64 {
	return a.version
}

func (a *AggregateBase) SetVersion(version int64) *AggregateBase {
	a.version = version

	return a
}

func (a *AggregateBase) GetCurrentVersion() int64 {
	return a.version + int64(len(a.events))
}

func (a *AggregateBase) ClearEvents() {
	a.events = make([]Event, 0, 1)
	a.eventTypes = make([]EventType, 0, 1)
}

func (a *AggregateBase) IncreaseVersion() {
	a.version++
}

func (a *AggregateBase) IsNew() bool {
	return a.version == 0 && len(a.events) == 0
}

func (a *AggregateBase) GetLastUpdate() int64 {
	return a.lastUpdate
}

func (a *AggregateBase) SetLastUpdate(lastUpdate int64) *AggregateBase {
	a.lastUpdate = lastUpdate

	return a
}
