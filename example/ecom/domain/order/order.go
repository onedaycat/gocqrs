package order

import (
	"github.com/onedaycat/gocqrs"
)

type Order struct {
	gocqrs.AggregateBase
	ID         string
	CustomerID string
	Status     Status
	OrderItems []*OrderItem
}
