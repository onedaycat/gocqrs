package dynamodb

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/onedaycat/gocqrs"
)

const eventsnapshot = "eventsnapshot"
const eventstore = "eventstore"
const xid = "id"

var bIDIndex = aws.String("b-id-index")
var eIDIndex = aws.String("e-id-index")
var seventsnapshot = aws.String(eventsnapshot)
var seventstore = aws.String(eventstore)

type DynamoDBEventStore struct {
	db *dynamodb.DynamoDB
}

func NewDynamoDBEventStore(sess *session.Session) *DynamoDBEventStore {
	return &DynamoDBEventStore{
		db: dynamodb.New(sess),
	}
}

func (d *DynamoDBEventStore) TruncateTables() {
	output, err := d.db.Scan(&dynamodb.ScanInput{
		TableName: aws.String(eventstore),
	})
	if err != nil {
		panic(err)
	}
	if len(output.Items) == 0 {
		return
	}

	keyStores := make([]*dynamodb.WriteRequest, len(output.Items))
	for i := 0; i < len(output.Items); i++ {
		keyStores[i] = &dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{
				Key: map[string]*dynamodb.AttributeValue{
					"a":  &dynamodb.AttributeValue{S: output.Items[i]["a"].S},
					"id": &dynamodb.AttributeValue{S: output.Items[i]["id"].S},
				},
			},
		}
	}
	_, err = d.db.BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			eventstore: keyStores,
		},
	})
	if err != nil {
		panic(err)
	}

	output, err = d.db.Scan(&dynamodb.ScanInput{
		TableName: aws.String(eventsnapshot),
	})
	if err != nil {
		panic(err)
	}
	if len(output.Items) == 0 {
		return
	}
	keyStores = make([]*dynamodb.WriteRequest, len(output.Items))
	for i := 0; i < len(output.Items); i++ {
		keyStores[i] = &dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{
				Key: map[string]*dynamodb.AttributeValue{
					"id": &dynamodb.AttributeValue{S: output.Items[i]["id"].S},
				},
			},
		}
	}
	_, err = d.db.BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			eventsnapshot: keyStores,
		},
	})
	if err != nil {
		panic(err)
	}
}

func (d *DynamoDBEventStore) CreateSchema(enableStream bool) error {
	_, err := d.db.CreateTable(&dynamodb.CreateTableInput{
		BillingMode: aws.String("PAY_PER_REQUEST"),
		StreamSpecification: &dynamodb.StreamSpecification{
			StreamEnabled:  aws.Bool(enableStream),
			StreamViewType: aws.String("NEW_IMAGE"),
		},
		TableName: seventstore,
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("a"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("e"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("b"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("v"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("a"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("v"),
				KeyType:       aws.String("RANGE"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("a-id-index"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("a"),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String("id"),
						KeyType:       aws.String("RANGE"),
					},
				},
			},
			{
				IndexName: aws.String("e-id-index"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("e"),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String("id"),
						KeyType:       aws.String("RANGE"),
					},
				},
			},
			{
				IndexName: aws.String("b-id-index"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("b"),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String("id"),
						KeyType:       aws.String("RANGE"),
					},
				},
			},
		},
	})

	if err != nil {
		aerr, _ := err.(awserr.Error)
		if aerr.Code() != dynamodb.ErrCodeResourceInUseException {
			return err
		}
	}

	_, err = d.db.CreateTable(&dynamodb.CreateTableInput{
		BillingMode: aws.String("PAY_PER_REQUEST"),
		TableName:   seventsnapshot,
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("b"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("b-id-index"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("b"),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String("id"),
						KeyType:       aws.String("RANGE"),
					},
				},
			},
		},
	})

	if err != nil {
		aerr, _ := err.(awserr.Error)
		if aerr.Code() != dynamodb.ErrCodeResourceInUseException {
			return err
		}
	}

	return nil
}

func (d *DynamoDBEventStore) Get(id string, time int64, limit int, nextToken string) ([]*gocqrs.EventMessage, string, error) {
	return nil, "", nil
}

func (d *DynamoDBEventStore) GetByEventType(eventType gocqrs.EventType, time int64, limit int, nextToken string) ([]*gocqrs.EventMessage, error) {
	return nil, nil
}

func (d *DynamoDBEventStore) GetByAggregateType(aggType gocqrs.AggregateType, time int64, limit int, nextToken string) ([]*gocqrs.EventMessage, error) {
	return nil, nil
}

func (d *DynamoDBEventStore) GetSnapshot(id string) (*gocqrs.Snapshot, error) {
	output, err := d.db.GetItem(&dynamodb.GetItemInput{
		TableName: seventsnapshot,
		Key: map[string]*dynamodb.AttributeValue{
			xid: &dynamodb.AttributeValue{S: aws.String(id)},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(output.Item) == 0 {
		return nil, gocqrs.ErrNotFound
	}

	snapshot := &gocqrs.Snapshot{}
	if err = dynamodbattribute.UnmarshalMap(output.Item, snapshot); err != nil {
		return nil, err
	}

	return snapshot, nil
}

func (d *DynamoDBEventStore) GetSnapshotsByAggregateType(aggType gocqrs.AggregateType, limit int, nextToken string) ([]*gocqrs.Snapshot, string, error) {
	var newNextToken string
	keyCon := aws.String("b=:agg")

	exValue := map[string]*dynamodb.AttributeValue{
		":agg": &dynamodb.AttributeValue{S: aws.String(string(aggType))},
	}

	if nextToken != "" {
		keyCon = aws.String("b=:agg and id>:nextToken")
		exValue[":nextToken"] = &dynamodb.AttributeValue{S: aws.String(nextToken)}
	}

	output, err := d.db.Query(&dynamodb.QueryInput{
		TableName:                 seventsnapshot,
		IndexName:                 bIDIndex,
		Limit:                     aws.Int64(int64(limit)),
		KeyConditionExpression:    keyCon,
		ExpressionAttributeValues: exValue,
	})

	if len(output.Items) == 0 {
		return nil, "", nil
	}

	snapshots := make([]*gocqrs.Snapshot, 0, len(output.Items))
	if err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &snapshots); err != nil {
		return nil, newNextToken, err
	}

	if len(output.LastEvaluatedKey) > 0 {
		newNextToken = *output.LastEvaluatedKey["id"].S
	}

	return snapshots, newNextToken, nil
}

func (d *DynamoDBEventStore) BeginTx(fn func(ctx context.Context) error) error {
	return fn(context.Background())
}

func (d *DynamoDBEventStore) Save(ctx context.Context, payloads []*gocqrs.EventMessage, snapshot *gocqrs.Snapshot) error {
	var err error
	var snapshotReq map[string]*dynamodb.AttributeValue
	snapshotReq, err = dynamodbattribute.MarshalMap(snapshot)
	if err != nil {
		return err
	}

	var payloadReq map[string]*dynamodb.AttributeValue
	putES := make([]*dynamodb.TransactWriteItem, 0, len(payloads)+1)

	putES = append(putES, &dynamodb.TransactWriteItem{
		Put: &dynamodb.Put{
			TableName:           seventsnapshot,
			ConditionExpression: aws.String("attribute_not_exists(id) or v<:v"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":v": &dynamodb.AttributeValue{N: aws.String(strconv.Itoa(snapshot.Version))},
			},
			Item: snapshotReq,
		},
	})

	for i := 0; i < len(payloads); i++ {
		payloadReq, err = dynamodbattribute.MarshalMap(payloads[i])
		if err != nil {
			return err
		}

		putES = append(putES, &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				TableName: seventstore,
				// ConditionExpression: aws.String("v <> :v"),
				// ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				// 	":v": &dynamodb.AttributeValue{N: aws.String(strconv.Itoa(payloads[i].Version))},
				// },
				Item: payloadReq,
			},
		})
	}

	_, err = d.db.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: putES,
	})

	return err
}
