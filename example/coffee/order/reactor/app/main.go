package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/onedaycat/gocqrs"
	"github.com/onedaycat/gocqrs/example/coffee/order/order/command"
	"github.com/onedaycat/gocqrs/example/coffee/order/reactor"
	"github.com/onedaycat/gocqrs/lambdastream/dynamostream"
	"github.com/onedaycat/gocqrs/storage/dynamodb"
	"github.com/rs/zerolog/log"
)

var (
	ds *dynamostream.DyanmoStream
)

func init() {
	storage := dynamodb.New(session.New(), "eventstore", "eventsnapshot")
	es := gocqrs.NewEventStore(storage, nil)
	cmd := command.NewService(es)
	rt := reactor.NewHandler(cmd)
	ds = dynamostream.New()
	ds.OnApplyEventMessage(rt.Apply)
	ds.OnError(func(msg *gocqrs.EventMessage, err error) {
		log.Error().Fields(map[string]interface{}{"msg": msg}).Msg(err.Error())
	})
}

func main() {
	log.Info().Msg("Init success")
	lambda.Start(ds.Run)
}
