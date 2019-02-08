package repository

import (
	"context"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/onedaycat/gocqrs/example/coffee/order/order/query"
)

type OrderMongoRepository struct {
	db         *mongo.Database
	dbName     string
	orderTable *mongo.Collection
}

func NewOrderMongoRepository(mongoClient *mongo.Client, dbName string) *OrderMongoRepository {
	return &OrderMongoRepository{
		db:         mongoClient.Database(dbName),
		dbName:     dbName,
		orderTable: mongoClient.Database(dbName).Collection("order"),
	}
}

func (m *OrderMongoRepository) GetOrder(id string) (*query.Order, error) {
	cursor := m.orderTable.FindOne(context.Background(), bson.D{
		{"_id", id},
	})

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	order := &query.Order{}
	if err := cursor.Decode(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (m *OrderMongoRepository) SaveOrder(order *query.Order) error {
	_, err := m.orderTable.ReplaceOne(context.Background(), bson.D{{"_id", order.ID}}, order)

	return err
}
