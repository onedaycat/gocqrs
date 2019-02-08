package command

import (
	"github.com/onedaycat/gocqrs"
	"github.com/onedaycat/gocqrs/example/coffee/order/order/domain"
)

type Service interface {
	PlaceOrder(cmd *PlaceOrder) (*domain.Order, error)
	AcceptOrder(cmd *AcceptOrder) (*domain.Order, error)
	StartOrder(cmd *StartOrder) (*domain.Order, error)
	FinishOrder(cmd *FinishOrder) (*domain.Order, error)
	DeliveryOrder(cmd *DeliveryOrder) (*domain.Order, error)
}

type service struct {
	es gocqrs.EventStore
}

func NewService(es gocqrs.EventStore) Service {
	return &service{
		es: es,
	}
}

func (s *service) PlaceOrder(cmd *PlaceOrder) (*domain.Order, error) {
	order := domain.NewOrder()
	order.PlaceOrder()

	if err := s.es.Save(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *service) AcceptOrder(cmd *AcceptOrder) (*domain.Order, error) {
	order := domain.NewOrder()
	if err := s.es.GetSnapshot(cmd.ID, order); err != nil {
		return nil, err
	}

	order.AcceptOrder()

	if err := s.es.Save(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *service) StartOrder(cmd *StartOrder) (*domain.Order, error) {
	order := domain.NewOrder()
	if err := s.es.GetSnapshot(cmd.ID, order); err != nil {
		return nil, err
	}

	order.StartOrder()

	if err := s.es.Save(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *service) FinishOrder(cmd *FinishOrder) (*domain.Order, error) {
	order := domain.NewOrder()
	if err := s.es.GetSnapshot(cmd.ID, order); err != nil {
		return nil, err
	}

	order.FinishOrder()

	if err := s.es.Save(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *service) DeliveryOrder(cmd *DeliveryOrder) (*domain.Order, error) {
	order := domain.NewOrder()
	if err := s.es.GetSnapshot(cmd.ID, order); err != nil {
		return nil, err
	}

	order.DeliveryOrder()

	if err := s.es.Save(order); err != nil {
		return nil, err
	}

	return order, nil
}
