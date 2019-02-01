package dynamodb

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/onedaycat/gocqrs"
	"github.com/onedaycat/gocqrs/example/ecom/domain/stock"
	"github.com/stretchr/testify/require"
)

// func TestGet(t *testing.T) {
// 	sess, err := session.NewSession(&aws.Config{
// 		Credentials: credentials.NewEnvCredentials(),
// 		Region:      aws.String("ap-southeast-1"),
// 	})
// 	require.NoError(t, err)

// 	db := NewDynamoDBEventStore(sess)

// 	err = db.CreateSchema(true)
// 	require.NoError(t, err)

// 	db.TruncateTables()
// 	st := stock.NewStockItem("1")
// 	st.Add(10)
// 	st.Sub(5)
// 	st.Add(2)
// 	st.Add(3)

// 	es := gocqrs.NewEventStore(db, nil)

// 	err = es.Save(st)
// 	require.NoError(t, err)

// 	st2 := stock.NewStockItem("1")
// 	err = es.GetSnapshot(st2)
// 	require.NoError(t, err)
// 	require.Equal(t, "1", st2.ProductID)
// 	require.Equal(t, stock.Qty(10), st2.Qty)

// 	st2.Remove()
// 	err = es.Save(st2)
// 	require.NoError(t, err)
// 	err = es.GetSnapshot(st2)
// 	require.NoError(t, err)
// 	require.True(t, st2.IsRemovedAggregate())
// }

// func TestGetSnapshotsByAggregateType(t *testing.T) {
// 	sess, err := session.NewSession(&aws.Config{
// 		Credentials: credentials.NewEnvCredentials(),
// 		Region:      aws.String("ap-southeast-1"),
// 	})
// 	require.NoError(t, err)

// 	db := NewDynamoDBEventStore(sess)

// 	err = db.CreateSchema(true)
// 	require.NoError(t, err)

// 	db.TruncateTables()
// 	es := gocqrs.NewEventStore(db, nil)

// 	st1 := stock.NewStockItem("1")
// 	st2 := stock.NewStockItem("2")
// 	st3 := stock.NewStockItem("3")

// 	es.Save(st1)
// 	es.Save(st2)
// 	es.Save(st3)

// 	sts, nextToken, err := es.GetSnapshotsByAggregateType("domain.subdomain.aggregate", 2, "")
// 	require.NoError(t, err)
// 	require.Len(t, sts, 2)
// 	require.Equal(t, "1", sts[0].ID)
// 	require.Equal(t, "2", sts[1].ID)
// 	require.Equal(t, "2", nextToken)

// 	sts, nextToken, err = es.GetSnapshotsByAggregateType("domain.subdomain.aggregate", 2, nextToken)
// 	require.NoError(t, err)
// 	require.Len(t, sts, 1)
// 	require.Equal(t, "3", sts[0].ID)
// 	require.Equal(t, "", nextToken)

// 	sts, nextToken, err = es.GetSnapshotsByAggregateType("notfound", 2, "")
// 	require.NoError(t, err)
// 	require.Nil(t, sts)
// 	require.Equal(t, "", nextToken)
// }

func TestGetSnapshot(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewEnvCredentials(),
		Region:      aws.String("ap-southeast-1"),
	})
	require.NoError(t, err)

	db := NewDynamoDBEventStore(sess)

	err = db.CreateSchema(true)
	require.NoError(t, err)

	db.TruncateTables()
	es := gocqrs.NewEventStore(db, nil)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		st := stock.NewStockItem("1")
		st.Add(10)
		st.Sub(5)
		st.Add(2)
		st.Add(3)

		if err := es.Save(st); err != nil {
			panic(err)
		}

		wg.Done()
	}()

	// go func() {
	// 	st := stock.NewStockItem("1")
	// 	st.Remove()

	// 	if err := es.Save(st); err != nil {
	// 		panic(err)
	// 	}

	// 	wg.Done()
	// }()

	wg.Wait()

	require.NoError(t, err)

	// st.Remove()
	// st.SetVersion(1)
	// err = es.Save(st)
	// require.NoError(t, err)

	// expSt := stock.NewStockItem("1")

	// err = es.GetSnapshot("1", expSt)
	// require.NoError(t, err)
	// require.Equal(t, expSt, st)
}
