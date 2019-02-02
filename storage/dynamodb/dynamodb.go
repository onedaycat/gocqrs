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

const (
	xid                  = "id"
	emptyStr             = ""
	seqKV                = ":s"
	getKV                = ":a"
	getByEventTypeKV     = ":et"
	getByAggregateTypeKV = ":b"
)

var (
	eSIndex                        = aws.String("e-s-index")
	bSIndex                        = aws.String("b-s-index")
	saveCond                       = aws.String("attribute_not_exists(s)")
	getCond                        = aws.String("a=:a")
	getCondWithTime                = aws.String("a=:a and s > :s")
	getByEventTypeCond             = aws.String("e=:et")
	getByEventTypeWithTimeCond     = aws.String("e=:et and s > :s")
	getByAggregateTypeCond         = aws.String("b=:b")
	getByAggregateTypeWithTimeCond = aws.String("b=:b and s > :s")
)

type DynamoDBEventStore struct {
	db              *dynamodb.DynamoDB
	eventstoreTable string
	snapshotTable   string
}

func New(sess *session.Session, eventstoreTable, snapshotTable string) *DynamoDBEventStore {
	return &DynamoDBEventStore{
		db:              dynamodb.New(sess),
		eventstoreTable: eventstoreTable,
		snapshotTable:   snapshotTable,
	}
}

func (d *DynamoDBEventStore) TruncateTables() {
	output, err := d.db.Scan(&dynamodb.ScanInput{
		TableName: &d.eventstoreTable,
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
					"a": &dynamodb.AttributeValue{S: output.Items[i]["a"].S},
					"s": &dynamodb.AttributeValue{N: output.Items[i]["s"].N},
				},
			},
		}
	}
	_, err = d.db.BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			d.eventstoreTable: keyStores,
		},
	})
	if err != nil {
		panic(err)
	}

	output, err = d.db.Scan(&dynamodb.ScanInput{
		TableName: &d.snapshotTable,
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
			d.snapshotTable: keyStores,
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
		TableName: &d.eventstoreTable,
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("a"),
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
				AttributeName: aws.String("s"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("a"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("s"),
				KeyType:       aws.String("RANGE"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("e-s-index"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("e"),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String("s"),
						KeyType:       aws.String("RANGE"),
					},
				},
			},
			{
				IndexName: aws.String("b-s-index"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("b"),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String("s"),
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
		TableName:   &d.snapshotTable,
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
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

func (d *DynamoDBEventStore) Get(aggID string, seq, limit int64) ([]*gocqrs.EventMessage, error) {
	keyCond := getCond
	exValue := map[string]*dynamodb.AttributeValue{
		getKV: &dynamodb.AttributeValue{S: &aggID},
	}

	if seq > 0 {
		exValue[seqKV] = &dynamodb.AttributeValue{N: aws.String(strconv.FormatInt(seq, 10))}
		keyCond = getCondWithTime
	}

	output, err := d.db.Query(&dynamodb.QueryInput{
		TableName:                 &d.eventstoreTable,
		KeyConditionExpression:    keyCond,
		Limit:                     &limit,
		ExpressionAttributeValues: exValue,
	})

	if err != nil {
		return nil, err
	}

	if len(output.Items) == 0 {
		return nil, nil
	}

	snapshots := make([]*gocqrs.EventMessage, 0, len(output.Items))
	if err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &snapshots); err != nil {
		return nil, err
	}

	return snapshots, nil
}

func (d *DynamoDBEventStore) GetByEventType(eventType gocqrs.EventType, seq, limit int64) ([]*gocqrs.EventMessage, error) {
	keyCond := getByEventTypeCond
	exValue := map[string]*dynamodb.AttributeValue{
		getByEventTypeKV: &dynamodb.AttributeValue{S: &eventType},
	}

	if seq > 0 {
		exValue[seqKV] = &dynamodb.AttributeValue{N: aws.String(strconv.FormatInt(seq, 10))}
		keyCond = getByEventTypeWithTimeCond
	}

	output, err := d.db.Query(&dynamodb.QueryInput{
		TableName:                 &d.eventstoreTable,
		IndexName:                 eSIndex,
		Limit:                     &limit,
		KeyConditionExpression:    keyCond,
		ExpressionAttributeValues: exValue,
	})

	if err != nil {
		return nil, err
	}

	if len(output.Items) == 0 {
		return nil, nil
	}

	snapshots := make([]*gocqrs.EventMessage, 0, len(output.Items))
	if err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &snapshots); err != nil {
		return nil, err
	}

	return snapshots, nil
}

func (d *DynamoDBEventStore) GetByAggregateType(aggType gocqrs.AggregateType, seq, limit int64) ([]*gocqrs.EventMessage, error) {
	keyCond := getByAggregateTypeCond
	exValue := map[string]*dynamodb.AttributeValue{
		getByAggregateTypeKV: &dynamodb.AttributeValue{S: &aggType},
	}

	if seq > 0 {
		exValue[seqKV] = &dynamodb.AttributeValue{N: aws.String(strconv.FormatInt(seq, 10))}
		keyCond = getByAggregateTypeWithTimeCond
	}

	output, err := d.db.Query(&dynamodb.QueryInput{
		TableName:                 &d.eventstoreTable,
		IndexName:                 bSIndex,
		Limit:                     &limit,
		KeyConditionExpression:    keyCond,
		ExpressionAttributeValues: exValue,
	})

	if err != nil {
		return nil, err
	}

	if len(output.Items) == 0 {
		return nil, nil
	}

	snapshots := make([]*gocqrs.EventMessage, 0, len(output.Items))
	if err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &snapshots); err != nil {
		return nil, err
	}

	return snapshots, nil
}

func (d *DynamoDBEventStore) GetSnapshot(aggID string) (*gocqrs.Snapshot, error) {
	output, err := d.db.GetItem(&dynamodb.GetItemInput{
		TableName: &d.snapshotTable,
		Key: map[string]*dynamodb.AttributeValue{
			xid: &dynamodb.AttributeValue{S: &aggID},
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
			TableName: &d.snapshotTable,
			Item:      snapshotReq,
		},
	})

	for i := 0; i < len(payloads); i++ {
		payloadReq, err = dynamodbattribute.MarshalMap(payloads[i])
		if err != nil {
			return err
		}

		putES = append(putES, &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				TableName:           &d.eventstoreTable,
				ConditionExpression: saveCond,
				Item:                payloadReq,
			},
		})
	}

	_, err = d.db.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: putES,
	})

	if err != nil {
		aerr := err.(awserr.Error)
		if aerr.Code() == dynamodb.ErrCodeTransactionCanceledException {
			return gocqrs.ErrVersionInconsistency
		}

		return err
	}

	return nil
}
