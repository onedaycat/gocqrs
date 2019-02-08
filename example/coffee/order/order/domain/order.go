package domain

import (
	"github.com/onedaycat/gocqrs"
)

type Order struct {
	*gocqrs.AggregateBase
	Price  int64  `json:"price"`
	Status Status `json:"status"`
}

func NewOrder() *Order {
	return &Order{
		AggregateBase: gocqrs.InitAggregate(),
	}
}

func (o *Order) GetAggregateType() string {
	return "coffee.order.order"
}

func (o *Order) PlaceOrder() {
	o.Price = 100
	o.Status = OrderPlacedStatus

	o.Publish(OrderPlacedEvent, &OrderPlaced{
		Price:  o.Price,
		Status: o.Status.String(),
	})
}

func (o *Order) AcceptOrder() {
	o.Status = OrderAcceptedStatus

	o.Publish(OrderStatusUpdatedEvent, &OrderStatusUpdated{
		Status: OrderPlacedStatus.String(),
	})
}

func (o *Order) StartOrder() {
	o.Status = OrderStartedStatus

	o.Publish(OrderStatusUpdatedEvent, &OrderStatusUpdated{
		Status: OrderStartedStatus.String(),
	})
}

func (o *Order) FinishOrder() {
	o.Status = OrderFinishedStatus

	o.Publish(OrderStatusUpdatedEvent, &OrderStatusUpdated{
		Status: OrderFinishedStatus.String(),
	})
}

func (o *Order) DeliveryOrder() {
	o.Status = OrderDeliveredStatus

	o.Publish(OrderStatusUpdatedEvent, &OrderStatusUpdated{
		Status: OrderDeliveredStatus.String(),
	})
}

func (o *Order) Apply(msg *gocqrs.EventMessage) error {
	switch msg.Type {
	case OrderPlacedEvent:
		event := &OrderPlaced{}
		if err := msg.Payload.UnmarshalPayload(event); err != nil {
			return err
		}

		o.Price = 100
		o.Status = Status(event.Status)

	case OrderStatusUpdatedEvent:
		event := &OrderStatusUpdated{}
		if err := msg.Payload.UnmarshalPayload(event); err != nil {
			return err
		}

		o.Status = Status(event.Status)
	}

	return nil
}
