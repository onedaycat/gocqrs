package cart

import (
	"github.com/onedaycat/gocqrs"
)

type Cart struct {
	gocqrs.AggregateBase
	CustomerID   string
	ProductItems []*ProductItem
}
