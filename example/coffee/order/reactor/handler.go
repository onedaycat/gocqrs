package reactor

import (
	"github.com/onedaycat/gocqrs"
	"github.com/onedaycat/gocqrs/example/coffee/order/order/command"
)

type Handler struct {
	cmdSrv command.Service
}

func NewHandler(cmdSrv command.Service) *Handler {
	return &Handler{
		cmdSrv: cmdSrv,
	}
}

func (h *Handler) updateStatusFromBarista(id, status string) error {
	switch status {
	case CoffeeBrewStartedStatus:
		_, err := h.cmdSrv.StartOrder(&command.StartOrder{id})
		return err
	case CoffeeBrewFinishedStatus:
		_, err := h.cmdSrv.FinishOrder(&command.FinishOrder{id})
		return err
	case CoffeeBrewDeliveredStatus:
		_, err := h.cmdSrv.DeliveryOrder(&command.DeliveryOrder{id})
		return err
	}

	return nil
}

func (h *Handler) Apply(msg *gocqrs.EventMessage) error {
	switch msg.Type {
	case CoffeBrewStatusUpdatedEvent:
		event := &CoffeBrewStatusUpdated{}
		if err := msg.Payload.UnmarshalPayload(event); err != nil {
			return err
		}

		return h.updateStatusFromBarista(msg.AggregateID, event.Status)
	case OrderBeanValidatedEvent:
		_, err := h.cmdSrv.AcceptOrder(&command.AcceptOrder{})
		return err
	}

	return nil
}
