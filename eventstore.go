package gocqrs

import (
	"context"
	"time"

	"github.com/onedaycat/gocqrs/common/clock"
	"github.com/onedaycat/gocqrs/common/eid"
)

// RetryHandler if return bool is true is allow retry,
// if return bool is false no retry
type RetryHandler func() error

const emptyStr = ""

//go:generate mockery -name=EventStore
type EventStore interface {
	SetEventLimit(limit int64)
	Get(aggID string, seq int64, agg AggregateRoot) ([]*EventMessage, error)
	GetByEventType(eventType EventType, seq int64) ([]*EventMessage, error)
	GetByAggregateType(aggType AggregateType, seq int64) ([]*EventMessage, error)
	GetAggregate(aggID string, agg AggregateRoot, seq int64) error
	GetSnapshot(aggID string, agg AggregateRoot) error
	Save(agg AggregateRoot) error
}

type eventStore struct {
	storage  Storage
	eventBus EventBus
	limit    int64
}

func NewEventStore(storage Storage, eventBus EventBus) EventStore {
	return &eventStore{storage, eventBus, 100}
}

func (es *eventStore) SetEventLimit(limit int64) {
	es.limit = limit
}

func (es *eventStore) GetAggregate(id string, agg AggregateRoot, seq int64) error {
	events, err := es.storage.Get(id, seq, es.limit)
	if err != nil {
		return err
	}

	n := len(events)

	if n == 0 {
		return ErrNotFound
	}

	lastEvent := events[n-1]

	agg.SetAggregateID(lastEvent.AggregateID)
	agg.SetVersion(lastEvent.Version)
	agg.SetLastUpdate(lastEvent.Time)

	for _, event := range events {
		if err = agg.Apply(event); err != nil {
			return err
		}
	}

	for n >= int(es.limit) {
		if err = es.GetAggregate(id, agg, lastEvent.Seq); err != nil {
			if err == ErrNotFound {
				break
			}
			return err
		}
	}

	if agg.IsNew() {
		return ErrNotFound
	}

	return nil
}

func (es *eventStore) Get(id string, seq int64, agg AggregateRoot) ([]*EventMessage, error) {
	return es.storage.Get(id, seq, es.limit)
}

func (es *eventStore) GetByEventType(eventType EventType, seq int64) ([]*EventMessage, error) {
	return es.storage.GetByEventType(eventType, seq, es.limit)
}

func (es *eventStore) GetByAggregateType(aggType AggregateType, seq int64) ([]*EventMessage, error) {
	return es.storage.GetByAggregateType(aggType, seq, es.limit)
}

func (es *eventStore) GetSnapshot(id string, agg AggregateRoot) error {
	snapshot, err := es.storage.GetSnapshot(id)
	if err != nil {
		return err
	}

	agg.SetAggregateID(snapshot.ID)
	agg.SetVersion(snapshot.Version)
	agg.SetLastUpdate(snapshot.LastUpdate)
	if snapshot.IsRemoved {
		agg.MarkAsRemoved()
	}

	if err = snapshot.Payload.UnmarshalPayload(agg); err != nil {
		return err
	}

	return nil
}

func (es *eventStore) Save(agg AggregateRoot) error {
	events := agg.GetEvents()
	if len(events) == 0 {
		return nil
	}

	if len(events) > 9 {
		return ErrEventLimitExceed
	}

	if agg.GetAggregateID() == emptyStr {
		return ErrNoAggregateID
	}

	payloads := make([]*EventMessage, len(events))
	now := clock.Now().Unix()
	aggType := agg.GetAggregateType()
	eventTypes := agg.GetEventTypes()

	var lastEvent *EventMessage

	for i := 0; i < len(events); i++ {
		agg.IncreaseVersion()
		aggid := agg.GetAggregateID()
		version := agg.GetVersion()
		id := eid.CreateEventID(aggid, version)
		payloads[i] = &EventMessage{
			ID:            id,
			AggregateID:   aggid,
			AggregateType: aggType,
			Version:       version,
			Type:          eventTypes[i],
			Payload:       NewPayload(events[i]),
			Time:          now,
			Seq:           NewSeq(now, version),
		}

		if len(events)-1 == i {
			lastEvent = payloads[i]
		}
	}

	snapshot := &Snapshot{
		ID:            agg.GetAggregateID(),
		AggregateType: aggType,
		Version:       agg.GetVersion(),
		Payload:       NewPayload(agg),
		LastEvent:     lastEvent,
		LastUpdate:    lastEvent.Time,
		IsRemoved:     agg.IsRemoved(),
	}

	if snapshot.Version == 0 {
		return ErrInvalidVersionNotAllowed
	}

	return es.storage.BeginTx(func(ctx context.Context) error {

		if err := es.storage.Save(ctx, payloads, snapshot); err != nil {
			return err
		}

		if es.eventBus != nil {
			if err := es.eventBus.Publish(payloads); err != nil {
				return err
			}
		}

		agg.ClearEvents()
		agg.SetLastUpdate(lastEvent.Time)

		return nil
	})
}

func WithRetry(numberRetry int, delay time.Duration, fn RetryHandler) error {
	var err error
	currentRetry := 0
	for currentRetry < numberRetry {
		if err = fn(); err != nil {
			if err == ErrVersionInconsistency {
				if delay > 0 {
					time.Sleep(delay)
				}

				currentRetry++
				continue
			}

			return err
		}

		return nil
	}

	return nil
}

func NewSeq(time int64, version int64) int64 {
	if time < 0 {
		return 0
	}

	if version < 0 {
		version = 0
	}

	return (time * 100000) + version
}
