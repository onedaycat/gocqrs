package query

type Service interface {
	GetOrder(id string) (*Order, error)
	SaveOrder(order *Order) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) GetOrder(id string) (*Order, error) {
	return s.repo.GetOrder(id)
}

func (s *service) SaveOrder(order *Order) error {
	return s.repo.SaveOrder(order)
}
