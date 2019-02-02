package dynamostream

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDynamoDBPayload(t *testing.T) {
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
							"S": "10001"
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
			}
		]
	}
	`)

	type pdata struct {
		ID string `json:"id"`
	}

	p := &DynamoDBPayload{}
	err := json.Unmarshal(payload, p)
	require.NoError(t, err)
	require.Len(t, p.Records, 1)
	require.Equal(t, eventInsert, p.Records[0].EventName)
	require.Equal(t, "domain.aggregate", p.Records[0].DynamoDB.Payload.EventMessage.AggregateType)

	pp := &pdata{}
	err = p.Records[0].DynamoDB.Payload.EventMessage.Payload.UnmarshalPayload(pp)
	require.NoError(t, err)
	require.Equal(t, &pdata{"1"}, pp)
}

func BenchmarkA(b *testing.B) {
	data := make(map[string]interface{})
	for i := 0; i < 10; i++ {
		data[strconv.Itoa(i)] = 10
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := data["1"]
		_ = d
	}
}

func BenchmarkB(b *testing.B) {
	data := "1"
	x := "1"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		switch data {
		case "0":
		case "1":
			d := x
			_ = d
		case "2":
		case "3":
		case "4":
		case "5":
		case "6":
		case "7":
		case "8":
		case "9":
		case "10":
		}
	}
}
