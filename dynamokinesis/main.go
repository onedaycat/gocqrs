package main

import (
	"context"
	"errors"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/service/kinesis"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/onedaycat/gocqrs"
	"github.com/onedaycat/gocqrs/lambdastream/dynamostream"
	"github.com/rs/zerolog/log"
)

var (
	ks         *kinesis.Kinesis
	streamName = os.Getenv("KINESIS_STREAM_NAME")
)

func handler(ctx context.Context, stream *dynamostream.DynamoDBStreamEvent) error {
	n := len(stream.Records)
	dataList := make([]*kinesis.PutRecordsRequestEntry, n)
	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(index int, event *gocqrs.EventMessage) {
			data, _ := event.Marshal()
			dataList[index] = &kinesis.PutRecordsRequestEntry{
				Data:         data,
				PartitionKey: &event.PartitionKey,
			}
			wg.Done()
		}(i, stream.Records[i].DynamoDB.NewImage.EventMessage)
	}
	wg.Wait()

	out, err := ks.PutRecords(&kinesis.PutRecordsInput{
		Records:    dataList,
		StreamName: &streamName,
	})

	if err != nil {
		return err
	}

	if out.FailedRecordCount != nil && *out.FailedRecordCount > 0 {
		return errors.New("One or more events published failed")
	}

	return nil
}

func init() {
	log.Info().Msg("Start init")
	sess, err := session.NewSession()
	if err != nil {
		log.Panic().Msg("AWS Session error: " + err.Error())
	}

	ks = kinesis.New(sess)
	log.Info().Msg("Done init")
}

func main() {
	log.Info().Msg("Start Dynamodb to Kinesis")
	lambda.Start(handler)
}
