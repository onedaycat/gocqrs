package command

import (
	"github.com/onedaycat/gocqrs"
	"github.com/onedaycat/gocqrs/example/ecom/stock/domain"
)

type StockService interface{}

type stockService struct {
	es gocqrs.EventStore
}

func NewHandler() StockService {
	return &stockService{}
}

func (h *stockService) CreateStockProduct(cmd *CreateStockProduct) (string, error) {
	st := domain.NewStockItem()
	st.Create(cmd.ProductID, domain.Qty(cmd.Qty))

	if err := h.es.Save(st); err != nil {
		return "", err
	}

	return st.GetAggregateID(), nil
}

func (h *stockService) UpdateStockQuantity(cmd *UpdateStockQuantity) (string, error) {
	st := domain.NewStockItem()
	if err := h.es.GetSnapshot(cmd.ID, st); err != nil {
		return "", err
	}

	if cmd.Qty < 0 {
		st.Sub(domain.Qty(cmd.Qty))
	} else if cmd.Qty > 0 {
		st.Add(domain.Qty(cmd.Qty))
	}

	if err := h.es.Save(st); err != nil {
		return "", err
	}

	return st.GetAggregateID(), nil
}

func (h *stockService) RemoveStockProduct(cmd *UpdateStockQuantity) (string, error) {
	st := domain.NewStockItem()
	if err := h.es.GetSnapshot(cmd.ID, st); err != nil {
		return "", err
	}

	st.Remove()

	if err := h.es.Save(st); err != nil {
		return "", err
	}

	return st.GetAggregateID(), nil
}
