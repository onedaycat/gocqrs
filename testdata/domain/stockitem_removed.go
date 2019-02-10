package domain

import (
	"github.com/onedaycat/gocqrs"
)

const StockItemRemovedEvent = "ecom:StockItemRemoved"

type StockItemRemoved struct {
	ProductID string `json:"productID"`
	RemovedAt int64  `json:"removedAt"`
}

func (e *StockItemRemoved) GetEventType() gocqrs.EventType {
	return StockItemRemovedEvent
}
