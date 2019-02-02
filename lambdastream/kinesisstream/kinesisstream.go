package kinesisstream

import (
	"context"

	"github.com/onedaycat/gocqrs"
)

type OnApplyEventMessageHandler func(msg *gocqrs.EventMessage) error
type OnErrorHandler func(msg *gocqrs.EventMessage, err error)

type KinesisStream struct {
	onError             OnErrorHandler
	onApplyEventMessage OnApplyEventMessageHandler
}

func New() *KinesisStream {
	return &KinesisStream{
		onError: func(msg *gocqrs.EventMessage, err error) {},
	}
}

func (s *KinesisStream) OnError(fn OnErrorHandler) {
	s.onError = fn
}

func (s *KinesisStream) OnApplyEventMessage(fn OnApplyEventMessageHandler) {
	s.onApplyEventMessage = fn
}

func (s *KinesisStream) Run(ctx context.Context, event *KinesisStreamEvent) (interface{}, error) {
	if s.onApplyEventMessage == nil {
		return nil, nil
	}

	var err error
	var msg *gocqrs.EventMessage
	for _, record := range event.Records {
		msg = record.Kinesis.Payload.EventMessage
		if err = s.onApplyEventMessage(record.Kinesis.Payload.EventMessage); err != nil {
			s.onError(msg, err)
		}
	}

	return nil, nil
}
