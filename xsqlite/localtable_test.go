package xsqlite

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"go.olapie.com/x/xtest"

	"github.com/google/uuid"
)

type localTableItem struct {
	ID     int64
	Text   string
	Number float64
	List   []int
}

func setupLocalTable(t testing.TB) *LocalTable[*localTableItem] {
	t.Log("setupLocalTable")
	if err := os.MkdirAll("testdata", 0755); err != nil {
		t.Fatal(err)
	}
	filename := "testdata/localtable" + fmt.Sprint(time.Now().UnixMilli()) + ".db"
	db, err := Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(
		func() {
			db.Close()
			os.Remove(filename)
		})
	return NewLocalTable[*localTableItem](db, func(opts *LocalTableOptions[*localTableItem]) {
		opts.Password = uuid.NewString()
	})
}

func newLocalTableItem() *localTableItem {
	return &localTableItem{
		ID:     rand.Int63(),
		Text:   time.Now().String(),
		Number: rand.Float64(),
		List:   []int{rand.Int(), rand.Int()},
	}
}

func TestLocalTable_SaveRemote(t *testing.T) {
	t.Log("TestLocalTable_SaveRemote")

	ctx := context.TODO()
	table := setupLocalTable(t)
	item := newLocalTableItem()
	localID := uuid.NewString()
	err := table.SaveRemote(ctx, localID, 0, item, time.Now().Unix())
	xtest.NoError(t, err)
	record, err := table.Get(ctx, localID)
	xtest.NoError(t, err)
	xtest.Equal(t, item, record)

	item.Text = time.Now().String() + "new"
	err = table.SaveRemote(ctx, localID, 0, item, time.Now().Unix())
	xtest.NoError(t, err)
	record, err = table.Get(ctx, localID)
	xtest.NoError(t, err)
	xtest.Equal(t, item, record)

	records, err := table.ListLocals(ctx)
	xtest.NoError(t, err)
	xtest.Equal(t, 0, len(records))

	records, err = table.ListDeletions(ctx)
	xtest.NoError(t, err)
	xtest.Equal(t, 0, len(records))

	records, err = table.ListUpdates(ctx)
	xtest.NoError(t, err)
	xtest.Equal(t, 0, len(records))

	records, err = table.ListRemotes(ctx)
	xtest.NoError(t, err)
	xtest.Equal(t, 1, len(records))

	err = table.Delete(ctx, localID)
	xtest.NoError(t, err)
	records, err = table.ListRemotes(ctx)
	xtest.NoError(t, err)
	xtest.Equal(t, 0, len(records))

	records, err = table.ListDeletions(ctx)
	xtest.NoError(t, err)
	xtest.Equal(t, 1, len(records))
	xtest.Equal(t, item, records[0])

}

func TestLocalTable_SaveLocal(t *testing.T) {
	ctx := context.TODO()
	table := setupLocalTable(t)
	t.Run("SyncedRemote", func(t *testing.T) {
		item := newLocalTableItem()
		localID := uuid.NewString()
		err := table.SaveLocal(ctx, localID, 0, item)
		xtest.NoError(t, err)
		record, err := table.Get(ctx, localID)
		xtest.NoError(t, err)
		xtest.Equal(t, item, record)

		item.Text = time.Now().String() + "new"
		err = table.SaveLocal(ctx, localID, 0, item)
		xtest.NoError(t, err)
		record, err = table.Get(ctx, localID)
		xtest.NoError(t, err)
		xtest.Equal(t, item, record)

		locals, err := table.ListLocals(ctx)
		xtest.NoError(t, err)
		xtest.Equal(t, 1, len(locals))

		remoteItem := newLocalTableItem()
		err = table.SaveRemote(ctx, localID, 0, remoteItem, time.Now().Unix())
		xtest.NoError(t, err)

		locals, err = table.ListLocals(ctx)
		xtest.NoError(t, err)
		xtest.Equal(t, 0, len(locals))
	})

	t.Run("DeleteLocal", func(t *testing.T) {
		item := newLocalTableItem()
		localID := uuid.NewString()
		err := table.SaveLocal(ctx, localID, 0, item)
		xtest.NoError(t, err)
		record, err := table.Get(ctx, localID)
		xtest.NoError(t, err)
		xtest.Equal(t, item, record)

		item.Text = time.Now().String() + "new"
		err = table.SaveLocal(ctx, localID, 0, item)
		xtest.NoError(t, err)
		record, err = table.Get(ctx, localID)
		xtest.NoError(t, err)
		xtest.Equal(t, item, record)

		locals, err := table.ListLocals(ctx)
		xtest.NoError(t, err)
		xtest.Equal(t, 1, len(locals))

		err = table.Delete(ctx, localID)
		xtest.NoError(t, err)

		locals, err = table.ListLocals(ctx)
		xtest.NoError(t, err)
		xtest.Equal(t, 0, len(locals))

		deletes, err := table.ListDeletions(ctx)
		xtest.NoError(t, err)
		xtest.Equal(t, 0, len(deletes))
	})
}

func TestLocalTable_Update(t *testing.T) {
	ctx := context.TODO()
	t.Run("UpdateRemote", func(t *testing.T) {
		table := setupLocalTable(t)
		item := newLocalTableItem()
		id := uuid.NewString()
		err := table.SaveRemote(ctx, id, 0, item, time.Now().Unix())
		xtest.NoError(t, err)
		got, err := table.Get(ctx, id)
		xtest.NoError(t, err)
		xtest.Equal(t, item, got)

		item.Text = xtest.RandomString(20)
		err = table.UpdateRemote(ctx, id, item)
		xtest.NoError(t, err)
		got, err = table.Get(ctx, id)
		xtest.NoError(t, err)
		xtest.Equal(t, item, got)
	})

	t.Run("UpdateLocal", func(t *testing.T) {
		table := setupLocalTable(t)
		item := newLocalTableItem()
		id := uuid.NewString()
		err := table.SaveLocal(ctx, id, 0, item)
		xtest.NoError(t, err)
		got, err := table.Get(ctx, id)
		xtest.NoError(t, err)
		xtest.Equal(t, item, got)

		item.Text = xtest.RandomString(20)
		err = table.UpdateLocal(ctx, id, item)
		xtest.NoError(t, err)
		got, err = table.Get(ctx, id)
		xtest.NoError(t, err)
		xtest.Equal(t, item, got)
	})

	t.Run("Update", func(t *testing.T) {
		table := setupLocalTable(t)
		t.Run("Remote", func(t *testing.T) {
			item := newLocalTableItem()
			id := uuid.NewString()
			err := table.SaveRemote(ctx, id, 0, item, time.Now().Unix())
			xtest.NoError(t, err)
			got, err := table.Get(ctx, id)
			xtest.NoError(t, err)
			xtest.Equal(t, item, got)

			item.Text = xtest.RandomString(20)
			err = table.Update(ctx, id, item)
			xtest.NoError(t, err)
			got, err = table.Get(ctx, id)
			xtest.NoError(t, err)
			xtest.Equal(t, item, got)
		})

		t.Run("Local", func(t *testing.T) {
			item := newLocalTableItem()
			id := uuid.NewString()
			err := table.SaveLocal(ctx, id, 0, item)
			xtest.NoError(t, err)
			got, err := table.Get(ctx, id)
			xtest.NoError(t, err)
			xtest.Equal(t, item, got)

			item.Text = xtest.RandomString(20)
			err = table.UpdateLocal(ctx, id, item)
			xtest.NoError(t, err)
			got, err = table.Get(ctx, id)
			xtest.NoError(t, err)
			xtest.Equal(t, item, got)
		})
	})
}

func BenchmarkLocalTable_SaveLocal(b *testing.B) {
	ctx := context.TODO()
	table := setupLocalTable(b)
	for i := 0; i < 100; i++ {
		item := newLocalTableItem()
		localID := uuid.NewString()
		err := table.SaveLocal(ctx, localID, 0, item)
		xtest.NoError(b, err)
	}

	var ids []string
	for i := 0; i < 100; i++ {
		item := newLocalTableItem()
		localID := uuid.NewString()
		err := table.SaveRemote(ctx, localID, 0, item, time.Now().Unix())
		xtest.NoError(b, err)
		ids = append(ids, localID)
	}
	for i := 0; i < 30; i++ {
		err := table.Delete(ctx, uuid.NewString())
		xtest.NoError(b, err)
	}

	for i := 0; i < b.N; i++ {
		item := newLocalTableItem()
		localID := uuid.NewString()
		table.SaveLocal(ctx, localID, 0, item)
	}

	for i := 0; i < b.N; i++ {
		table.Get(ctx, ids[i%len(ids)])
	}
}

func BenchmarkLocalTable_Get(b *testing.B) {
	ctx := context.TODO()
	table := setupLocalTable(b)
	for i := 0; i < 100; i++ {
		item := newLocalTableItem()
		localID := uuid.NewString()
		err := table.SaveLocal(ctx, localID, 0, item)
		xtest.NoError(b, err)
	}

	var ids []string
	for i := 0; i < 100; i++ {
		item := newLocalTableItem()
		localID := uuid.NewString()
		err := table.SaveRemote(ctx, localID, 0, item, time.Now().Unix())
		xtest.NoError(b, err)
		ids = append(ids, localID)
	}
	for i := 0; i < 30; i++ {
		err := table.Delete(ctx, uuid.NewString())
		xtest.NoError(b, err)
	}

	for i := 0; i < b.N; i++ {
		table.Get(ctx, ids[i%len(ids)])
	}
}

func BenchmarkLocalTable_ListLocals(b *testing.B) {
	ctx := context.TODO()
	table := setupLocalTable(b)
	for i := 0; i < 100; i++ {
		item := newLocalTableItem()
		localID := uuid.NewString()
		err := table.SaveLocal(ctx, localID, 0, item)
		xtest.NoError(b, err)
	}

	var ids []string
	for i := 0; i < 100; i++ {
		item := newLocalTableItem()
		localID := uuid.NewString()
		err := table.SaveRemote(ctx, localID, 0, item, time.Now().Unix())
		xtest.NoError(b, err)
		ids = append(ids, localID)
	}
	for i := 0; i < 30; i++ {
		err := table.Delete(ctx, uuid.NewString())
		xtest.NoError(b, err)
	}

	for i := 0; i < b.N; i++ {
		table.ListLocals(ctx)
	}
}

func BenchmarkLocalTable_ListRemotes(b *testing.B) {
	ctx := context.TODO()
	table := setupLocalTable(b)
	for i := 0; i < 100; i++ {
		item := newLocalTableItem()
		localID := uuid.NewString()
		err := table.SaveLocal(ctx, localID, 0, item)
		xtest.NoError(b, err)
	}

	var ids []string
	for i := 0; i < 100; i++ {
		item := newLocalTableItem()
		localID := uuid.NewString()
		err := table.SaveRemote(ctx, localID, 0, item, time.Now().Unix())
		xtest.NoError(b, err)
		ids = append(ids, localID)
	}
	for i := 0; i < 30; i++ {
		err := table.Delete(ctx, uuid.NewString())
		xtest.NoError(b, err)
	}

	for i := 0; i < b.N; i++ {
		table.ListRemotes(ctx)
	}
}
