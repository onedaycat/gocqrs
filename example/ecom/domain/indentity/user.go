package domain

import "github.com/onedaycat/gocqrs"

type User struct {
	gocqrs.AggregateBase
	ID       string
	Name     string
	Email    string
	Password string
}
