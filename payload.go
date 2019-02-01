package gocqrs

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/mongodb/mongo-go-driver/bson"
)

type Payload struct {
	jsonData   json.RawMessage
	bsonData   bson.Raw
	dynamoData *dynamodb.AttributeValue
	data       interface{}
}

func NewPayload(data interface{}) *Payload {
	return &Payload{
		data: data,
	}
}

func (p *Payload) UnmarshalPayload(v interface{}) error {
	switch {
	case p.jsonData != nil:
		return json.Unmarshal(p.jsonData, v)
	case p.bsonData != nil:
		return bson.Unmarshal(p.bsonData, v)
	case p.dynamoData != nil:
		return dynamodbattribute.Unmarshal(p.dynamoData, v)
	}

	return ErrEncodingNotSupported
}

func (p *Payload) UnmarshalJSON(b []byte) error {
	p.jsonData = b

	return nil
}

func (p *Payload) UnmarshalBSON(b []byte) error {
	p.bsonData = b

	return nil
}

func (p *Payload) UnmarshalDynamoDBAttributeValue(v *dynamodb.AttributeValue) error {
	p.dynamoData = v

	return nil
}

func (p *Payload) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.data)
}

func (p *Payload) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	xav, err := dynamodbattribute.Marshal(p.data)
	*av = *xav

	return err
}

func (p *Payload) MarshalBSON() ([]byte, error) {
	return bson.Marshal(p.data)
}
