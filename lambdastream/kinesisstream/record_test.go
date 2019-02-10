package kinesisstream

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseKinesisStreamEvent(t *testing.T) {
	_ = []byte(`
	{
		"a": "a1",
		"b": "domain.aggregate",
		"s": 10,
		"e": "domain.aggregate.event",
		"x": 10001,
		"p": {
			"id": "1"
		}
	}
	`)

	payload := []byte(`
	{
		"Records": [
			{
				"kinesis": {
					"kinesisSchemaVersion": "1.0",
					"partitionKey": "1",
					"sequenceNumber": "49590338271490256608559692538361571095921575989136588898",
					"data": "Cgl7CgkJImEiOiAiYTEiLAoJCSJiIjogImRvbWFpbi5hZ2dyZWdhdGUiLAoJCSJ2IjogMTAsCgkJImUiOiAiZG9tYWluLmFnZ3JlZ2F0ZS5ldmVudCIsCgkJInMiOiAxMDAwMSwKCQkicCI6IHsKCQkJImlkIjogIjEiCgkJfQoJfQoJ",
					"approximateArrivalTimestamp": 1545084650.987
				},
				"eventSource": "aws:kinesis",
				"eventVersion": "1.0",
				"eventID": "shardId-000000000006:49590338271490256608559692538361571095921575989136588898",
				"eventName": "aws:kinesis:record",
				"invokeIdentityArn": "arn:aws:iam::123456789012:role/lambda-role",
				"awsRegion": "us-east-2",
				"eventSourceARN": "arn:aws:kinesis:us-east-2:123456789012:stream/lambda-stream"
			},
			{
				"kinesis": {
					"kinesisSchemaVersion": "1.0",
					"partitionKey": "1",
					"sequenceNumber": "49590338271490256608559692540925702759324208523137515618",
					"data": "Cgl7CgkJImEiOiAiYTEiLAoJCSJiIjogImRvbWFpbi5hZ2dyZWdhdGUiLAoJCSJ2IjogMTEsCgkJImUiOiAiZG9tYWluLmFnZ3JlZ2F0ZS5ldmVudCIsCgkJInMiOiAxMDAwMSwKCQkicCI6IHsKCQkJImlkIjogIjEiCgkJfQoJfQoJ",
					"approximateArrivalTimestamp": 1545084711.166
				},
				"eventSource": "aws:kinesis",
				"eventVersion": "1.0",
				"eventID": "shardId-000000000006:49590338271490256608559692540925702759324208523137515618",
				"eventName": "aws:kinesis:record",
				"invokeIdentityArn": "arn:aws:iam::123456789012:role/lambda-role",
				"awsRegion": "us-east-2",
				"eventSourceARN": "arn:aws:kinesis:us-east-2:123456789012:stream/lambda-stream"
			}
		]
	}
	`)

	type pdata struct {
		ID string `json:"id"`
	}

	event := &KinesisStreamEvent{}
	err := json.Unmarshal(payload, event)
	require.NoError(t, err)
	require.Len(t, event.Records, 2)
	require.Equal(t, "domain.aggregate", event.Records[0].Kinesis.Data.EventMessage.AggregateType)
	require.Equal(t, "domain.aggregate.event", event.Records[0].Kinesis.Data.EventMessage.EventType)
	// require.Equal(t, int64(10), event.Records[0].Kinesis.Data.EventMessage.Seq)
	// require.Equal(t, int64(11), event.Records[1].Kinesis.Data.EventMessage.Seq)

	pp := &pdata{}
	err = event.Records[0].Kinesis.Data.EventMessage.Payload.UnmarshalPayload(pp)
	fmt.Println(pp)
	require.NoError(t, err)
	require.Equal(t, &pdata{"1"}, pp)
}
