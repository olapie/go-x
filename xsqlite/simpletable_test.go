package xsqlite

import (
	"database/sql"
	"errors"
	"math/rand"
	"strings"
	"testing"

	"github.com/google/uuid"
	"go.olapie.com/x/xtest"
)

func createTable[K SimpleKey, R any](t *testing.T, pkFn func(R) K) *SimpleTable[K, R] {
	t.Log("createTable")

	db, err := sql.Open("sqlite", "file::memory:")
	if err != nil {
		t.Fatal(err)
	}

	name := "test" + strings.ReplaceAll(uuid.NewString(), "-", "")
	tbl, err := NewSimpleTable[K, R](db, name, pkFn)
	xtest.NoError(t, err)
	return tbl
}

type IntItem struct {
	ID    int64
	Name  string
	Score float64
}

func newIntItem() *IntItem {
	return &IntItem{
		ID:    rand.Int63(),
		Name:  uuid.NewString(),
		Score: float64(rand.Int63()) / float64(3),
	}
}

type StringItem struct {
	ID    string
	Name  string
	Score float64
}

func (i *StringItem) PrimaryKey() string {
	return i.ID
}

func newStringItem() *StringItem {
	return &StringItem{
		ID:    uuid.NewString(),
		Name:  uuid.NewString(),
		Score: float64(rand.Int63()) / float64(3),
	}
}

func TestIntTable(t *testing.T) {
	t.Log("TestIntTable")
	tbl := createTable[int64, *IntItem](t, func(item *IntItem) int64 {
		return item.ID
	})
	var items []*IntItem
	item := newIntItem()
	items = append(items, item)
	err := tbl.Insert(item)
	xtest.NoError(t, err)
	v, err := tbl.Get(item.ID)
	xtest.NoError(t, err)
	xtest.Equal(t, item, v)

	item = newIntItem()
	item.ID = items[0].ID + 1
	err = tbl.Insert(item)
	xtest.NoError(t, err)
	items = append(items, item)

	l, err := tbl.ListAll()
	xtest.NoError(t, err)
	xtest.True(t, len(l) != 0)
	xtest.Equal(t, items, l)

	l, err = tbl.ListGreaterThan(item.ID, 10)
	xtest.NoError(t, err)
	xtest.True(t, len(l) == 0)

	l, err = tbl.ListLessThan(item.ID+1, 10)
	xtest.NoError(t, err)
	xtest.Equal(t, 2, len(l))
	//t.Log(l[0].ID, l[1].ID)
	xtest.True(t, l[0].ID < l[1].ID)

	err = tbl.Delete(item.ID)
	xtest.NoError(t, err)

	v, err = tbl.Get(item.ID)
	xtest.Error(t, err)
	xtest.Equal(t, true, errors.Is(err, sql.ErrNoRows))
}

func TestStringTable(t *testing.T) {
	t.Log("TestStringTable")

	tbl := createTable[string, *StringItem](t, func(item *StringItem) string {
		return item.ID
	})
	var items []*StringItem
	item := newStringItem()
	t.Log(item.PrimaryKey())
	items = append(items, item)
	err := tbl.Insert(item)
	xtest.NoError(t, err)
	v, err := tbl.Get(item.ID)
	xtest.NoError(t, err)
	xtest.Equal(t, item, v)

	item = newStringItem()
	t.Log(item.PrimaryKey())
	err = tbl.Insert(item)
	xtest.NoError(t, err)
	items = append(items, item)

	l, err := tbl.ListAll()
	t.Log(len(l), err)
	xtest.NoError(t, err, "ListAll")
	xtest.NotEqual(t, 0, len(l))
	xtest.Equal(t, items, l)

	l, err = tbl.ListGreaterThan("\x01", 10)
	xtest.NoError(t, err, "ListGreaterThan")
	xtest.Equal(t, len(l), 2)

	l, err = tbl.ListLessThan("\xFF", 10)
	xtest.NoError(t, err)
	xtest.Equal(t, 2, len(l))
	//t.Log(l[0].ID, l[1].ID)
	xtest.True(t, l[0].ID < l[1].ID, "ListLessThan")

	err = tbl.Delete(item.ID)
	xtest.NoError(t, err)

	v, err = tbl.Get(item.ID)
	xtest.Error(t, err)
	xtest.Equal(t, true, errors.Is(err, sql.ErrNoRows))
}
