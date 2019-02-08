package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/onedaycat/gocqrs"
	"github.com/onedaycat/gocqrs/example/coffee/order/order/query"
	"github.com/onedaycat/gocqrs/example/coffee/order/order/query/repository"
	"github.com/onedaycat/gocqrs/example/coffee/order/projector"
	"github.com/onedaycat/gocqrs/lambdastream/dynamostream"
	"github.com/rs/zerolog/log"
)

var (
	ds *dynamostream.DyanmoStream
)

func init() {
	mgo, err := mongo.NewClient(os.Getenv("APP_MONGODB_ENDPOINT"))
	if err != nil {
		panic(err)
	}

	repo := repository.NewOrderMongoRepository(mgo, os.Getenv("APP_MONGODB_DATABASE"))
	qry := query.NewService(repo)
	pt := projector.NewHandler(qry)
	ds = dynamostream.New()
	ds.OnApplyEventMessage(pt.Apply)
	ds.OnError(func(msg *gocqrs.EventMessage, err error) {
		log.Error().Fields(map[string]interface{}{"msg": msg}).Msg(err.Error())
	})
}

func main() {
	log.Info().Msg("Init success")
	lambda.Start(ds.Run)
}
