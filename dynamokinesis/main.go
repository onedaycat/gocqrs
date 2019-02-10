package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/onedaycat/gocqrs"
	"github.com/onedaycat/gocqrs/eventbus/kinesis"
	"github.com/onedaycat/gocqrs/lambdastream/dynamostream"
	"github.com/rs/zerolog/log"
)

var (
	ks *kinesis.KinesisEventBus
)

func handler(ctx context.Context, stream *dynamostream.DynamoDBStreamEvent) error {
	events := make([]*gocqrs.EventMessage, 0, len(stream.Records))
	for _, record := range stream.Records {
		if record.EventName == dynamostream.EventInsert {
			events = append(events, record.DynamoDB.NewImage.EventMessage)
		}
	}

	if len(events) == 0 {
		return nil
	}

	err := ks.Publish(events)
	if err != nil {
		log.Error().Msg("Publish error: " + err.Error())
		return err
	}

	return nil
}

func init() {
	log.Info().Msg("Start init")
	sess, err := session.NewSession()
	if err != nil {
		log.Panic().Msg("AWS Session error: " + err.Error())
	}

	ks = kinesis.NewKinesisEventBus(sess, os.Getenv("KINESIS_STREAM_NAME"))
	log.Info().Msg("Done init")
}

func main() {
	log.Info().Msg("Start Dynamodb to Kinesis")
	lambda.Start(handler)
}
