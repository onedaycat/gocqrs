package command

import (
	"github.com/onedaycat/gocqrs"
	"github.com/onedaycat/gocqrs/example/ecom/stock/domain"
)

type StockCommandHandler interface{}

type stockCommandHandler struct {
	es gocqrs.EventStore
}

func NewHandler() StockCommandHandler {
	return &stockCommandHandler{}
}

func (h *stockCommandHandler) CreateStockProduct(cmd *CreateStockProduct) error {
	st := domain.NewStockItem()
	h.es.GetSnapshot(cmd.ProductID, st)
	return nil
}
