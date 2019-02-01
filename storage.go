package gocqrs

import (
	"context"
)

// Get(id string, withSnapshot bool)
type Storage interface {
	Get(id string, time int64) ([]*EventMessage, error)
	GetByEventType(eventType EventType, time int64) ([]*EventMessage, error)
	GetByAggregateType(aggType AggregateType, time int64) ([]*EventMessage, error)
	GetSnapshot(id string) (*Snapshot, error)
	BeginTx(fn func(ctx context.Context) error) error
	Save(ctx context.Context, payloads []*EventMessage, snapshot *Snapshot) error
}
