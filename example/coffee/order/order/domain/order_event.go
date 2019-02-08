package domain

const (
	OrderPlacedEvent        = "coffee.order.order.OrderPlaced"
	OrderStatusUpdatedEvent = "coffee.order.order.OrderStatusUpdated"
)

type OrderPlaced struct {
	Price  int64  `json:"price"`
	Status string `json:"status"`
}

type OrderStatusUpdated struct {
	Status string `json:"status"`
}
