package dynamostream

import (
	"context"

	"github.com/onedaycat/gocqrs"
)

type OnApplyEventMessageHandler func(msg *gocqrs.EventMessage) error
type OnErrorHandler func(msg *gocqrs.EventMessage, err error)

type DyanmoStream struct {
	onError             OnErrorHandler
	onApplyEventMessage OnApplyEventMessageHandler
}

func New() *DyanmoStream {
	return &DyanmoStream{
		onError: func(msg *gocqrs.EventMessage, err error) {},
	}
}

func (s *DyanmoStream) OnError(fn OnErrorHandler) {
	s.onError = fn
}

func (s *DyanmoStream) OnApplyEventMessage(fn OnApplyEventMessageHandler) {
	s.onApplyEventMessage = fn
}

func (s *DyanmoStream) Run(ctx context.Context, event *DynamoDBStreamEvent) (interface{}, error) {
	if s.onApplyEventMessage == nil {
		return nil, nil
	}

	var err error
	var msg *gocqrs.EventMessage
	for _, record := range event.Records {
		if eventInsert != record.EventName {
			continue
		}

		msg = record.DynamoDB.Payload.EventMessage
		if err = s.onApplyEventMessage(record.DynamoDB.Payload.EventMessage); err != nil {
			s.onError(msg, err)
		}
	}

	return nil, nil
}
