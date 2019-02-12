package domain

const StockItemUpdatedEvent = "domain.subdomain.aggregate.StockItemUpdated"

type StockItemUpdated struct {
	ProductID string `json:"productID"`
	Qty       Qty    `json:"qty"`
}
