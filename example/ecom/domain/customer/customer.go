package customer

import (
	"github.com/onedaycat/gocqrs"
)

type Customer struct {
	gocqrs.AggregateBase
	ID    string
	Name  string
	Email string
}
