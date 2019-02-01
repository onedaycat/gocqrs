package local

import (
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/onedaycat/gocqrs"
)

type FakeEventBus struct {
	kin *kinesis.Kinesis
}

func FakecalEventBus() *FakeEventBus {
	return &FakeEventBus{}
}

func (k *FakeEventBus) Publish(events []*gocqrs.EventMessage) error {
	return nil
}
