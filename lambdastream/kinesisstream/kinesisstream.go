package kinesisstream

import (
	"context"

	"github.com/onedaycat/gocqrs"
)

type EventMessage = gocqrs.EventMessage
type EventMessages = []*gocqrs.EventMessage

type LambdaHandler func(ctx context.Context, event *KinesisStreamEvent) (interface{}, error)
type EventMessageHandler func(msg *EventMessage) error
type EventMessagesHandler func(msgs EventMessages) error
type EventMessageErrorHandler func(msg *EventMessage, err error)
type EventMessagesErrorHandler func(msgs EventMessages, err error)

type KinesisStream struct{}

func New() *KinesisStream {
	return &KinesisStream{}
}

func (s *KinesisStream) CreateIteratorHandler(handler EventMessageHandler, onError EventMessageErrorHandler) LambdaHandler {
	return func(ctx context.Context, event *KinesisStreamEvent) (interface{}, error) {
		if handler == nil {
			return nil, nil
		}
		if onError == nil {
			onError = func(msg *EventMessage, err error) {}
		}

		var err error
		var msg *gocqrs.EventMessage
		for _, record := range event.Records {
			msg = record.Kinesis.Data.EventMessage
			if err = handler(record.Kinesis.Data.EventMessage); err != nil {
				onError(msg, err)
			}
		}

		return nil, nil
	}
}

func (s *KinesisStream) CreateConcurencyHandler(handler EventMessageHandler, onError EventMessageErrorHandler) LambdaHandler {
	return func(ctx context.Context, event *KinesisStreamEvent) (interface{}, error) {
		if handler == nil {
			return nil, nil
		}
		if onError == nil {
			onError = func(msg *EventMessage, err error) {}
		}

		cm := newConcurrencyManager(len(event.Records))

		for _, record := range event.Records {
			cm.Send(record, handler, onError)
		}

		cm.Wait()

		return nil, nil
	}
}

func (s *KinesisStream) CreateGroupConcurencyHandler(handler EventMessagesHandler, onError EventMessagesErrorHandler) LambdaHandler {
	return func(ctx context.Context, event *KinesisStreamEvent) (interface{}, error) {
		if handler == nil {
			return nil, nil
		}
		if onError == nil {
			onError = func(msgs EventMessages, err error) {}
		}

		cm := newGroupConcurrencyManager(len(event.Records))

		cm.Send(event.Records, handler, onError)
		cm.Wait()

		return nil, nil
	}
}
