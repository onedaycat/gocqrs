package domain

import (
	"github.com/onedaycat/gocqrs"
)

type StockItem struct {
	*gocqrs.AggregateBase
	ProductID string
	Qty       Qty
}

func NewStockItem() *StockItem {
	return &StockItem{
		AggregateBase: gocqrs.InitAggregate(),
	}
}

func (st *StockItem) Create(productID string, qty Qty) {
	st.Qty = qty
	st.Publish(&StockItemCreated{
		ProductID: st.GetAggregateID(),
		Qty:       st.Qty,
	})
}

func (st *StockItem) GetAggregateType() gocqrs.AggregateType {
	return "domain.subdomain.aggregate"
}

func (st *StockItem) Apply(msg *gocqrs.EventMessage) error {
	switch msg.Type {
	case StockItemCreatedEvent:
		event := &StockItemCreated{}
		if err := msg.Payload.UnmarshalPayload(event); err != nil {
			return err
		}

		st.ProductID = event.ProductID
		st.Qty = event.Qty
	case StockItemUpdatedEvent:
		event := &StockItemUpdated{}
		if err := msg.Payload.UnmarshalPayload(event); err != nil {
			return err
		}

		st.ProductID = event.ProductID
		st.Qty = event.Qty
	case StockItemRemovedEvent:
		st.MarkAsRemoved()
		return nil
	}

	return nil
}

func (st *StockItem) Add(amount Qty) {
	st.Qty = st.Qty.Add(amount)
	st.Publish(&StockItemUpdated{
		ProductID: st.ProductID,
		Qty:       st.Qty,
	})
}

func (st *StockItem) Sub(amount Qty) error {
	var err error
	st.Qty, err = st.Qty.Sub(amount)
	if err != nil {
		return err
	}

	st.Publish(&StockItemUpdated{
		ProductID: st.ProductID,
		Qty:       st.Qty,
	})

	return nil
}

func (st *StockItem) Remove() error {
	st.MarkAsRemoved()

	st.Publish(&StockItemRemoved{
		ProductID: st.ProductID,
	})

	return nil
}
