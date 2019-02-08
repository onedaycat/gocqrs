package query

type Repository interface {
	GetOrder(id string) (*Order, error)
	SaveOrder(order *Order) error
}
