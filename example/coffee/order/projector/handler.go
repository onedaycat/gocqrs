package projector

import (
	"github.com/onedaycat/gocqrs"
	"github.com/onedaycat/gocqrs/example/coffee/order/order/domain"
	"github.com/onedaycat/gocqrs/example/coffee/order/order/query"
)

type Handler struct {
	querySrv query.Service
}

func NewHandler(querySrv query.Service) *Handler {
	return &Handler{
		querySrv: querySrv,
	}
}

func (h *Handler) Apply(msg *gocqrs.EventMessage) error {
	switch msg.Type {
	case domain.OrderPlacedEvent:
		event := &domain.OrderPlaced{}
		if err := msg.Payload.UnmarshalPayload(event); err != nil {
			return err
		}

		order := &query.Order{
			ID:     msg.AggregateID,
			Price:  event.Price,
			Status: event.Status,
		}

		return h.querySrv.SaveOrder(order)

	case domain.OrderStatusUpdatedEvent:
		event := &domain.OrderStatusUpdated{}
		if err := msg.Payload.UnmarshalPayload(event); err != nil {
			return err
		}

		order, err := h.querySrv.GetOrder(msg.AggregateID)
		if err != nil {
			return err
		}

		order.Status = event.Status

		return h.querySrv.SaveOrder(order)
	}

	return nil
}
