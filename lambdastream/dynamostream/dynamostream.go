package dynamostream

import (
	"context"
)

type DyanmoStream struct {
}

func New() *DyanmoStream {
	return &DyanmoStream{}
}

func (s *DyanmoStream) Run(ctx context.Context, event interface{}) (interface{}, error) {
	return nil, nil
}
