package gocqrs

import (
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
	SetSnapshotStrategy(strategies SnapshotStategyHandler)
	GetEvents(aggID string, seq int64, agg AggregateRoot) ([]*EventMessage, error)
	GetEventsByEventType(eventType EventType, seq int64) ([]*EventMessage, error)
	GetEventsByAggregateType(aggType AggregateType, seq int64) ([]*EventMessage, error)
	GetAggregate(aggID string, agg AggregateRoot) error
	GetSnapshot(aggID string, agg AggregateRoot) error
	Save(agg AggregateRoot) error
}

type eventStore struct {
	storage          Storage
	limit            int64
	snapshotStrategy SnapshotStategyHandler
}

func NewEventStore(storage Storage) EventStore {
	return &eventStore{
		storage:          storage,
		limit:            100,
		snapshotStrategy: LatestEventSanpshot(),
	}
}

func (es *eventStore) SetEventLimit(limit int64) {
	es.limit = limit
}

func (es *eventStore) SetSnapshotStrategy(strategies SnapshotStategyHandler) {
	es.snapshotStrategy = strategies
}

func (es *eventStore) GetAggregate(id string, agg AggregateRoot) error {
	err := es.GetSnapshot(id, agg)
	if err != nil && err != ErrNotFound {
		return err
	}

	return es.getAggregateFromEvent(id, agg, agg.GetSequence())
}

func (es *eventStore) getAggregateFromEvent(id string, agg AggregateRoot, seq int64) error {
	events, err := es.storage.GetEvents(id, seq, es.limit)
	if err != nil {
		return err
	}

	n := len(events)

	if n == 0 {
		return ErrNotFound
	}

	lastEvent := events[n-1]

	agg.SetAggregateID(lastEvent.AggregateID)
	agg.SetSequence(lastEvent.Seq)

	for _, event := range events {
		if err = agg.Apply(event); err != nil {
			return err
		}
	}

	for n >= int(es.limit) {
		if err = es.getAggregateFromEvent(id, agg, lastEvent.Seq); err != nil {
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

func (es *eventStore) GetEvents(id string, seq int64, agg AggregateRoot) ([]*EventMessage, error) {
	return es.storage.GetEvents(id, seq, es.limit)
}

func (es *eventStore) GetEventsByEventType(eventType EventType, seq int64) ([]*EventMessage, error) {
	return es.storage.GetEventsByEventType(eventType, seq, es.limit)
}

func (es *eventStore) GetEventsByAggregateType(aggType AggregateType, seq int64) ([]*EventMessage, error) {
	return es.storage.GetEventsByAggregateType(aggType, seq, es.limit)
}

func (es *eventStore) GetSnapshot(id string, agg AggregateRoot) error {
	snapshot, err := es.storage.GetSnapshot(id)
	if err != nil {
		return err
	}

	agg.SetAggregateID(snapshot.AggregateID)
	agg.SetSequence(snapshot.Seq)

	if err = snapshot.Payload.UnmarshalPayload(agg); err != nil {
		return err
	}

	return nil
}

func (es *eventStore) Save(agg AggregateRoot) error {
	payloads := agg.GetEventPayloads()
	n := len(payloads)
	if n == 0 {
		return nil
	}

	if n > 9 {
		return ErrEventLimitExceed
	}

	if agg.GetAggregateID() == emptyStr {
		return ErrNoAggregateID
	}

	events := make([]*EventMessage, n)
	now := clock.Now().Unix()
	aggType := agg.GetAggregateType()
	eventTypes := agg.GetEventTypes()

	var lastEvent *EventMessage

	for i := 0; i < n; i++ {
		agg.IncreaseSequence()
		aggid := agg.GetAggregateID()
		seq := agg.GetSequence()
		eid := eid.CreateEventID(aggid, seq)
		metadata := NewPayload(agg.GetMetadata())
		events[i] = &EventMessage{
			EventID:       eid,
			AggregateID:   aggid,
			AggregateType: aggType,
			Seq:           seq,
			EventType:     eventTypes[i],
			Payload:       NewPayload(payloads[i]),
			Time:          now,
			TimeSeq:       NewSeq(now, seq),
			Metadata:      metadata,
		}

		if len(payloads)-1 == i {
			lastEvent = events[i]
		}
	}

	if lastEvent.Seq == 0 {
		return ErrInvalidVersionNotAllowed
	}

	var snapshot *Snapshot
	if es.snapshotStrategy(agg, events) {
		snapshot = &Snapshot{
			AggregateID:   agg.GetAggregateID(),
			AggregateType: aggType,
			Payload:       NewPayload(agg),
			EventID:       lastEvent.EventID,
			Time:          lastEvent.Time,
			Seq:           lastEvent.Seq,
			TimeSeq:       lastEvent.TimeSeq,
			Metadata:      lastEvent.Metadata,
		}
	}

	if err := es.storage.Save(events, snapshot); err != nil {
		return err
	}

	agg.ClearEvents()

	return nil
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

func NewSeq(time int64, seq int64) int64 {
	if time < 0 {
		return 0
	}

	if seq < 0 {
		seq = 0
	}

	return (time * 100000) + seq
}
