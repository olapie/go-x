package xtype

import (
	"fmt"
	"github.com/google/uuid"
	"go.olapie.com/x/xtest"
	"math/rand/v2"
	"testing"
)

func TestNewUserID_Int64(t *testing.T) {
	id := rand.Int64()
	uid := NewUserID(id)
	got, _ := uid.Int()
	xtest.Equal(t, id, got)
	xtest.Equal(t, fmt.Sprint(id), uid.String())
	xtest.True(t, id == uid.Value())
}

func TestNewUserID_String(t *testing.T) {
	id := uuid.NewString()
	uid := NewUserID(id)
	_, ok := uid.Int()
	xtest.False(t, ok)
	xtest.Equal(t, id, uid.String())
	xtest.True(t, id == uid.Value())
}

func TestNewUserID_IntString(t *testing.T) {
	id := fmt.Sprint(rand.Int64())
	uid := NewUserID(id)
	got, _ := uid.Int()
	xtest.Equal(t, id, fmt.Sprint(got))
	xtest.Equal(t, id, uid.String())
	xtest.True(t, id == uid.Value())
}

func TestNewUserID_Int64_CustomType(t *testing.T) {
	type IDType int64
	id := IDType(rand.Int64())
	uid := NewUserID(id)
	got, _ := uid.Int()
	xtest.Equal(t, int64(id), got)
	xtest.Equal(t, fmt.Sprint(id), uid.String())
	xtest.True(t, id == uid.Value())
}

func TestNewUserID_String_CustomType(t *testing.T) {
	type IDType string
	id := IDType(uuid.NewString())
	uid := NewUserID(id)
	_, ok := uid.Int()
	xtest.False(t, ok)
	xtest.Equal(t, string(id), uid.String())
	xtest.True(t, id == uid.Value())
}
