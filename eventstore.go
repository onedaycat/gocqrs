package gocqrs

import (
	"context"
	"time"

	"github.com/onedaycat/gocqrs/internal/clock"
	"github.com/onedaycat/gocqrs/internal/eid"
)

// RetryHandler if return bool is true is allow retry,
// if return bool is false no retry
type RetryHandler func() (error, bool)

type EventStore interface {
	Get(id string, agg AggregateRoot) error
	GetByTime(id string, time int64, agg AggregateRoot) ([]*EventMessage, error)
	GetByEventType(eventType EventType, time int64) ([]*EventMessage, error)
	GetByAggregateType(aggType AggregateType, time int64) ([]*EventMessage, error)
	GetSnapshot(id string, agg AggregateRoot) error
	Save(agg AggregateRoot) error
}

type eventStore struct {
	storage  Storage
	eventBus EventBus
}

func NewEventStore(storage Storage, eventBus EventBus) EventStore {
	return &eventStore{storage, eventBus}
}

func (es *eventStore) Get(id string, agg AggregateRoot) error {
	events, err := es.storage.Get(id, 0)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err = agg.Apply(event); err != nil {
			return err
		}
	}

	return nil
}

func (es *eventStore) GetByTime(id string, time int64, agg AggregateRoot) ([]*EventMessage, error) {
	return es.storage.Get(id, time)
}

func (es *eventStore) GetByEventType(eventType EventType, time int64) ([]*EventMessage, error) {
	return es.storage.GetByEventType(eventType, time)
}

func (es *eventStore) GetByAggregateType(aggType AggregateType, time int64) ([]*EventMessage, error) {
	return es.storage.GetByAggregateType(aggType, time)
}

func (es *eventStore) GetSnapshot(id string, agg AggregateRoot) error {
	snapshot, err := es.storage.GetSnapshot(id)
	if err != nil {
		return err
	}

	agg.SetAggregateID(snapshot.ID)
	agg.SetVersion(snapshot.Version)

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

	payloads := make([]*EventMessage, len(events))
	now := clock.Now().Unix()
	aggType := agg.GetAggregateType()
	var lastEvent *EventMessage

	for i := 0; i < len(events); i++ {
		agg.IncreaseVersion()
		payloads[i] = &EventMessage{
			ID:            eid.New(agg.GetAggregateID(), agg.GetVersion()).String(),
			AggregateID:   agg.GetAggregateID(),
			AggregateType: aggType,
			Version:       agg.GetVersion(),
			Type:          events[i].GetEventType(),
			Payload:       NewPayload(events[i]),
			Time:          now,
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
		IsRemoved:     agg.IsRemoved(),
	}

	return es.storage.BeginTx(func(ctx context.Context) error {
		// curSnap, err := es.storage.GetSnapshot(snapshot.ID)
		// if err != nil && err != ErrNotFound {
		// 	return err
		// }

		// if err != ErrNotFound {
		// 	if curSnap.Version+1 != payloads[0].Version {
		// 		return ErrVersionInconsistency
		// 	}
		// }

		if err := es.storage.Save(ctx, payloads, snapshot); err != nil {
			return err
		}

		agg.ClearEvents()

		return nil
	})
}

func WithRetry(numberRetry int, delay time.Duration, fn RetryHandler) {
	var err error
	var isRetry bool
	currentRetry := 0
	for currentRetry < numberRetry {
		if err, isRetry = fn(); err == nil || !isRetry {
			return
		}

		if delay > 0 {
			time.Sleep(delay)
		}

		currentRetry++
	}
}
