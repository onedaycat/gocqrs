package gocqrs

import (
	"context"
	"fmt"
	"time"
)

type EventStore interface {
	Get(id string, agg AggregateRoot) error
	GetByTime(id string, time int64, agg AggregateRoot) error
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
	// payloads, err := es.storage.Get(id)
	// if err != nil {
	// 	return err
	// }

	// es.applyEvent(payloads, agg)

	return nil
}

func (es *eventStore) GetByTime(id string, time int64, agg AggregateRoot) error {
	// payloads, err := es.storage.GetByTime(id, time)
	// if err != nil {
	// 	return err
	// }

	// es.applyEvent(payloads, agg)

	return nil
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
	now := time.Now().UTC().Unix()
	aggType := agg.GetAggregateType()
	var lastEvent *EventMessage

	for i := 0; i < len(events); i++ {
		agg.IncreaseVersion()
		payloads[i] = &EventMessage{
			ID:            generateID(),
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
		curSnap, err := es.storage.GetSnapshot(snapshot.ID)
		if err != nil && err != ErrNotFound {
			return err
		}

		fmt.Println("wait 5s")
		time.Sleep(time.Second * 5)

		if err != ErrNotFound {
			if curSnap.Version+1 != payloads[0].Version {
				return ErrVersionInconsistency
			}
		}

		if err := es.storage.Save(ctx, payloads, snapshot); err != nil {
			return err
		}

		agg.ClearEvents()
		fmt.Println("done")

		return nil
	})

}
