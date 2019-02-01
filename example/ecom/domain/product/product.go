package domain

import "github.com/onedaycat/gocqrs"

type Product struct {
	gocqrs.AggregateBase
	ID          string
	Name        string
	Description string
}
