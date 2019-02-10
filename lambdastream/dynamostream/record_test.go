package dynamostream

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDynamoDBStreamEvent(t *testing.T) {
	payload := []byte(`
	{
		"Records": [
			{
				"eventID": "7de3041dd709b024af6f29e4fa13d34c",
				"eventName": "INSERT",
				"eventVersion": "1.1",
				"eventSource": "aws:dynamodb",
				"awsRegion": "us-west-2",
				"dynamodb": {
					"ApproximateCreationDateTime": 1479499740,
					"Keys": {
						"Timestamp": {
							"S": "2016-11-18:12:09:36"
						},
						"Username": {
							"S": "John Doe"
						}
					},
					"NewImage": {
						"a": {
							"S": "a1"
						},
						"b": {
							"S": "domain.aggregate"
						},
						"v": {
							"N": "10"
						},
						"e": {
							"S": "domain.aggregate.event"
						},
						"s": {
							"N": "10001"
						},
						"p": {
							"M": {
								"id": {
									"S": "1"
								}
							}
						}
					},
					"SequenceNumber": "13021600000000001596893679",
					"SizeBytes": 112,
					"StreamViewType": "NEW_IMAGE"
				},
				"eventSourceARN": "arn:aws:dynamodb:us-east-1:123456789012:table/BarkTable/stream/2016-11-16T20:42:48.104"
			},
			{
				"eventID":"3",
				"eventName":"REMOVE",
				"eventVersion":"1.0",
				"eventSource":"aws:dynamodb",
				"awsRegion":"us-east-1",
				"dynamodb":{
				   "Keys":{
					  "Id":{
						 "N":"101"
					  }
				   },
				   "NewImage": {
					"a": {
						"S": "a1"
					},
					"b": {
						"S": "domain.aggregate"
					},
					"v": {
						"N": "10"
					},
					"e": {
						"S": "domain.aggregate.event"
					},
					"s": {
						"N": "10001"
					},
					"p": {
						"M": {
							"id": {
								"S": "1"
							}
						}
					}
				},
				   "SequenceNumber":"333",
				   "SizeBytes":38,
				   "StreamViewType":"NEW_AND_OLD_IMAGES"
				},
				"eventSourceARN":"stream-ARN"
			 }
		]
	}
	`)

	type pdata struct {
		ID string `json:"id"`
	}

	event := &DynamoDBStreamEvent{}
	err := json.Unmarshal(payload, event)
	require.NoError(t, err)
	require.Len(t, event.Records, 2)
	require.Equal(t, EventInsert, event.Records[0].EventName)
	require.Equal(t, eventRemove, event.Records[1].EventName)
	require.Equal(t, "domain.aggregate", event.Records[0].DynamoDB.NewImage.EventMessage.AggregateType)

	xx, _ := json.Marshal(event.Records[0].DynamoDB.NewImage.EventMessage)
	fmt.Println(string(xx))

	pp := &pdata{}
	err = event.Records[0].DynamoDB.NewImage.EventMessage.Payload.UnmarshalPayload(pp)
	require.NoError(t, err)
	require.Equal(t, &pdata{"1"}, pp)
}
