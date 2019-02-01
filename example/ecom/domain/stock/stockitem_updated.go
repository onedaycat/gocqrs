package stock

import "github.com/onedaycat/gocqrs"

const StockItemUpdatedEvent = "domain.subdomain.aggregate.StockItemUpdated"

type StockItemUpdated struct {
	ProductID string `json:"productID"`
	Qty       Qty    `json:"qty"`
}

func (e *StockItemUpdated) GetEventType() gocqrs.EventType {
	return StockItemUpdatedEvent
}
