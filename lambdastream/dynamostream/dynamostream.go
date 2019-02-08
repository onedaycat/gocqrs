package dynamostream

import (
	"context"

	"github.com/onedaycat/gocqrs"
)

type OnEventMessageHandler func(msg *gocqrs.EventMessage) error
type OnErrorHandler func(msg *gocqrs.EventMessage, err error)

type DyanmoStream struct {
	onError              OnErrorHandler
	onNewEventMessage    OnEventMessageHandler
	onRemoveEventMessage OnEventMessageHandler
}

func New() *DyanmoStream {
	return &DyanmoStream{
		onError: func(msg *gocqrs.EventMessage, err error) {},
	}
}

func (s *DyanmoStream) OnError(fn OnErrorHandler) {
	s.onError = fn
}

func (s *DyanmoStream) OnNewEventMessage(fn OnEventMessageHandler) {
	s.onNewEventMessage = fn
}

func (s *DyanmoStream) OnRemoveEventMessage(fn OnEventMessageHandler) {
	s.onRemoveEventMessage = fn
}

func (s *DyanmoStream) Run(ctx context.Context, event *DynamoDBStreamEvent) (interface{}, error) {
	if s.onNewEventMessage == nil {
		return nil, nil
	}

	var err error
	var msg *gocqrs.EventMessage
	for _, record := range event.Records {
		if eventInsert != record.EventName {
			continue
		}

		switch record.EventName {
		case eventInsert:
			msg = record.DynamoDB.NewImage.EventMessage
			if err = s.onNewEventMessage(record.DynamoDB.NewImage.EventMessage); err != nil {
				s.onError(msg, err)
			}
		case eventRemove:
			msg = record.DynamoDB.NewImage.EventMessage
			if err = s.onRemoveEventMessage(record.DynamoDB.NewImage.EventMessage); err != nil {
				s.onError(msg, err)
			}
		default:
			continue
		}

	}

	return nil, nil
}
