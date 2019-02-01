package mongo

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/onedaycat/gocqrs"
	"github.com/onedaycat/gocqrs/example/ecom/domain/stock"
	"github.com/stretchr/testify/require"
)

func TestGet2(t *testing.T) {
	client, err := mongo.NewClient(os.Getenv("MONGODB_ENDPOINT"))
	require.NoError(t, err)
	ctx := context.Background()
	client.Connect(ctx)

	db := NewMongoEventStore(client, "eventsourcing_dev")

	db.DropSchema()
	err = db.CreateSchema()
	require.NoError(t, err)

	es := gocqrs.NewEventStore(db, nil)
	wg := sync.WaitGroup{}

	wg.Add(2)

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

	go func() {
		st := stock.NewStockItem("1")
		st.Remove()
		if err := es.Save(st); err != nil {
			panic(err)
		}
		wg.Done()
	}()

	wg.Wait()
}

func TestGet(t *testing.T) {
	client, err := mongo.NewClient(os.Getenv("MONGODB_ENDPOINT"))
	require.NoError(t, err)
	ctx := context.Background()
	client.Connect(ctx)

	db := NewMongoEventStore(client, "eventsourcing_dev")

	db.DropSchema()
	err = db.CreateSchema()
	require.NoError(t, err)

	es := gocqrs.NewEventStore(db, nil)

	st := stock.NewStockItem("1")
	st.Add(10)
	st.Sub(5)
	st.Add(2)
	st.Add(3)
	st.Remove()

	err = es.Save(st)
	require.NoError(t, err)

	expSt := &stock.StockItem{
		AggregateBase: gocqrs.InitAggregate("1"),
		ProductID:     "1",
		Qty:           10,
	}
	expSt.SetVersion(5)
	expSt.MarkAsRemoved()

	err = es.GetSnapshot("1", expSt)
	require.NoError(t, err)
	require.Equal(t, expSt, st)
}
