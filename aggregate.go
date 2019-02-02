package gocqrs

import (
	"github.com/onedaycat/gocqrs/internal/eid"
)

type AggregateType = string

type AggregateRoot interface {
	Apply(payload *EventMessage) error
	GetAggregateID() string
	SetAggregateID(id string)
	GetAggregateType() AggregateType
	GetVersion() int64
	SetVersion(version int64)
	GetCurrentVersion() int64
	GetLastUpdate() int64
	SetLastUpdate(t int64)
	IncreaseVersion()
	GetEvents() []Event
	ClearEvents()
	Publish(event Event)
	MarkAsRemoved()
	IsRemoved() bool
	IsNew() bool
}

type AggregateBase struct {
	id         string
	events     []Event
	removed    bool
	version    int64
	lastUpdate int64
}

// InitAggregate if id is empty, id will be generated
func InitAggregate() *AggregateBase {
	return &AggregateBase{
		id:      eid.GenerateAggregateID(),
		events:  make([]Event, 0, 1),
		version: 0,
	}
}

func InitAggregateWithID(id string) *AggregateBase {
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

func (a *AggregateBase) GetVersion() int64 {
	return a.version
}

func (a *AggregateBase) SetVersion(version int64) {
	a.version = version
}

func (a *AggregateBase) GetCurrentVersion() int64 {
	return a.version + int64(len(a.events))
}

func (a *AggregateBase) ClearEvents() {
	a.events = make([]Event, 0, 1)
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

func (a *AggregateBase) SetLastUpdate(lastUpdate int64) {
	a.lastUpdate = lastUpdate
}
