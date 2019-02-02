package gocqrs

import (
	"context"
)

//go:generate mockery -name=Storage
// Get(id string, withSnapshot bool)
type Storage interface {
	Get(aggID string, seq, limit int64) ([]*EventMessage, error)
	GetByEventType(eventType EventType, seq, limit int64) ([]*EventMessage, error)
	GetByAggregateType(aggType AggregateType, seq, limit int64) ([]*EventMessage, error)
	GetSnapshot(aggID string) (*Snapshot, error)
	BeginTx(fn func(ctx context.Context) error) error
	Save(ctx context.Context, payloads []*EventMessage, snapshot *Snapshot) error
}
