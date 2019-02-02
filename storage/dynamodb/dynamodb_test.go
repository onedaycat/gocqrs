package dynamodb

import (
	"sync"
	"testing"
	"time"

	"github.com/onedaycat/gocqrs/internal/clock"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/onedaycat/gocqrs"
	"github.com/onedaycat/gocqrs/example/ecom/domain/stock"
	"github.com/stretchr/testify/require"
)

var _db *DynamoDBEventStore

func getDB() *DynamoDBEventStore {
	if _db == nil {
		sess, err := session.NewSession(&aws.Config{
			Credentials: credentials.NewEnvCredentials(),
			Region:      aws.String("ap-southeast-1"),
		})
		if err != nil {
			panic(err)
		}

		_db = NewDynamoDBEventStore(sess)
		err = _db.CreateSchema(true)
		if err != nil {
			panic(err)
		}
	}

	return _db
}

func TestSaveAndGet(t *testing.T) {
	db := getDB()
	db.TruncateTables()

	es := gocqrs.NewEventStore(db, nil)

	now1 := time.Now().UTC().Add(time.Second * -10)
	now2 := time.Now().UTC().Add(time.Second * -5)

	st := stock.NewStockItem()
	st.Create("1", 0)
	st.Add(10)
	st.Sub(5)
	st.Add(2)
	st.Add(3)

	clock.Freeze(now1)
	err := es.Save(st)
	require.NoError(t, err)

	// Get
	st2 := stock.NewStockItem()
	err = es.Get(st.GetAggregateID(), st2)
	require.NoError(t, err)
	require.Equal(t, st, st2)

	// Get By Time
	st.Add(2)
	st.Remove()

	clock.Freeze(now2)
	err = es.Save(st)
	require.NoError(t, err)

	st3 := stock.NewStockItem()
	events, err := es.GetByTime(st.GetAggregateID(), now2.Unix(), st3)
	require.NoError(t, err)
	require.Len(t, events, 2)
	require.True(t, st.IsRemoved())
	require.Equal(t, stock.StockItemUpdatedEvent, events[0].Type)
	require.Equal(t, 6, events[0].Version)
	require.Equal(t, stock.StockItemRemovedEvent, events[1].Type)
	require.Equal(t, 7, events[1].Version)

	// GetSnapshot
	st4 := stock.NewStockItem()
	err = es.GetSnapshot(st.GetAggregateID(), st4)
	require.NoError(t, err)
	require.Equal(t, st4, st)

	// GetByEventType
	events, err = es.GetByEventType(stock.StockItemUpdatedEvent, now2.Unix())
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, stock.StockItemUpdatedEvent, events[0].Type)
	require.Equal(t, 6, events[0].Version)

	// GetByAggregateType
	events, err = es.GetByAggregateType(st.GetAggregateType(), now2.Unix())
	require.NoError(t, err)
	require.Len(t, events, 2)
	require.True(t, st.IsRemoved())
	require.Equal(t, stock.StockItemUpdatedEvent, events[0].Type)
	require.Equal(t, 6, events[0].Version)
	require.Equal(t, stock.StockItemRemovedEvent, events[1].Type)
	require.Equal(t, 7, events[1].Version)
}

func TestConcurency(t *testing.T) {
	db := getDB()

	db.TruncateTables()
	es := gocqrs.NewEventStore(db, nil)

	wg := sync.WaitGroup{}
	wg.Add(2)

	var err1 error
	var err2 error
	go func() {
		st := stock.NewStockItem()
		st.SetAggregateID("a1")
		st.Create("1", 0)
		st.Add(10)
		st.Sub(5)
		st.Add(2)
		st.Add(3)

		err1 = es.Save(st)

		wg.Done()
	}()

	go func() {
		st := stock.NewStockItem()
		st.SetAggregateID("a1")
		st.Create("1", 0)
		st.Remove()

		err2 = es.Save(st)

		wg.Done()
	}()

	wg.Wait()
	require.Equal(t, gocqrs.ErrVersionInconsistency, err1)
	require.Nil(t, err2)
}
