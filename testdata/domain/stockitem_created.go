package domain

const StockItemCreatedEvent = "domain.subdomain.aggregate.StockItemCreated"

type StockItemCreated struct {
	ProductID string `json:"productID"`
	Qty       Qty    `json:"qty"`
}
