package main

import (
	"context"

	"github.com/onedaycat/amuro/appsync"
	"github.com/onedaycat/gocqrs/example/coffee/order/order/command"
	"github.com/onedaycat/gocqrs/example/coffee/order/order/query"
)

type handler struct {
	cmdSrv command.Service
	qrySrv query.Service
}

func NewHandler(cmdSrv command.Service, qrySrv query.Service) *handler {
	return &handler{
		cmdSrv: cmdSrv,
		qrySrv: qrySrv,
	}
}

func (h *handler) MutationPlaceOrder(ctx context.Context, event *appsync.InvokeEvent) *appsync.Result {
	err := h.cmdSrv.PlaceOrder(&command.PlaceOrder{})
	if err != nil {
		return event.ErrorResult(err)
	}

	return nil
}

func (h *handler) QueryOrder(ctx context.Context, event *appsync.InvokeEvent) *appsync.Result {
	return nil
}
