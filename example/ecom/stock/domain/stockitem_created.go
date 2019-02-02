package domain

import (
	"github.com/onedaycat/gocqrs"
)

const StockItemCreatedEvent = "domain.subdomain.aggregate.StockItemCreated"

type StockItemCreated struct {
	ProductID string `json:"productID"`
	Qty       Qty    `json:"qty"`
}

func (e *StockItemCreated) GetEventType() gocqrs.EventType {
	return StockItemCreatedEvent
}
