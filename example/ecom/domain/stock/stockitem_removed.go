package stock

import (
	"github.com/onedaycat/gocqrs"
)

const StockItemRemovedEvent = "ecom:StockItemRemoved"

type StockItemRemoved struct {
	ProductID string `json:"productID"`
}

func (e *StockItemRemoved) GetEventType() gocqrs.EventType {
	return StockItemRemovedEvent
}
