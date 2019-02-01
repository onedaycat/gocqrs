package kinesis

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/onedaycat/gocqrs"
)

var (
	streamName  = "eventsource"
	sstreamName = aws.String("eventsource")
)

type KinesisEventBus struct {
	kin *kinesis.Kinesis
}

func NewKinesisEventBus(sess *session.Session) *KinesisEventBus {
	return &KinesisEventBus{
		kin: kinesis.New(sess),
	}
}

func (k *KinesisEventBus) Publish(events []*gocqrs.EventMessage) error {
	records := make([]*kinesis.PutRecordsRequestEntry, len(events))
	wg := sync.WaitGroup{}
	wg.Add(len(events))

	for i := 0; i < len(events); i++ {
		go func(index int, event *gocqrs.EventMessage) {
			data, _ := json.Marshal(event)
			records[index] = &kinesis.PutRecordsRequestEntry{
				Data:         data,
				PartitionKey: aws.String(event.AggregateID),
			}
			wg.Done()
		}(i, events[i])
	}
	wg.Wait()

	out, err := k.kin.PutRecords(&kinesis.PutRecordsInput{
		Records:    records,
		StreamName: sstreamName,
	})

	if *out.FailedRecordCount > 0 {
		log.Println("Kinesis Put Failed: ", *out.FailedRecordCount, "records")
	}

	return err
}
