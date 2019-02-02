package dynamostream

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	gocqrs "github.com/onedaycat/gocqrs"
)

type Payload struct {
	EventMessage *gocqrs.EventMessage
}

func (p *Payload) UnmarshalJSON(b []byte) error {
	var err error
	data := make(map[string]*dynamodb.AttributeValue)
	if err = json.Unmarshal(b, &data); err != nil {
		return err
	}

	event := &gocqrs.EventMessage{}
	if err = dynamodbattribute.UnmarshalMap(data, event); err != nil {
		return err
	}

	p.EventMessage = event
	return nil
}

const eventInsert = "INSERT"

type Records = []*Record

type DynamoDBPayload struct {
	Records Records `json:"Records"`
}

type Record struct {
	EventName string          `json:"eventName"`
	DynamoDB  *DynamoDBRecord `json:"dynamodb"`
}

type DynamoDBRecord struct {
	Payload *Payload `json:"NewImage"`
}
