package payment

import (
	"github.com/onedaycat/gocqrs"
)

type CreditCard struct {
	gocqrs.AggregateBase
	Name         string
	Number       string
	CVV          string
	ExpiredMonth string
	ExpiredYear  string
	OrderID      string
	Amount       int64
}
