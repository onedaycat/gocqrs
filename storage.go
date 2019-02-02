package gocqrs

import (
	"context"
)

// Get(id string, withSnapshot bool)
type Storage interface {
	Get(aggID string, time int64) ([]*EventMessage, error)
	GetByEventType(eventType EventType, time int64) ([]*EventMessage, error)
	GetByAggregateType(aggType AggregateType, time int64) ([]*EventMessage, error)
	GetSnapshot(aggID string) (*Snapshot, error)
	BeginTx(fn func(ctx context.Context) error) error
	Save(ctx context.Context, payloads []*EventMessage, snapshot *Snapshot) error
}
