package xcontext

import (
	"context"
	"fmt"
	"go.olapie.com/x/xtest"
	"math/rand/v2"
	"testing"
)

func TestGetIncomingUserID(t *testing.T) {
	ctx := context.TODO()
	id := rand.Int64()
	ctx = WithIncomingActivity(ctx, new(Activity))
	err := SetIncomingUserID[int64](ctx, id)
	if err != nil {
		t.Error(err)
		return
	}

	intID := GetIncomingUserID[int64](ctx)
	xtest.Equal(t, id, intID)

	strID := GetIncomingUserID[string](ctx)
	xtest.Equal(t, fmt.Sprint(id), strID)
}

func TestGetIncomingUserID_String(t *testing.T) {
	ctx := context.TODO()
	id := fmt.Sprint(rand.Int64())
	ctx = WithIncomingActivity(ctx, new(Activity))
	err := SetIncomingUserID[string](ctx, id)
	if err != nil {
		t.Error(err)
		return
	}

	intID := GetIncomingUserID[int64](ctx)
	xtest.Equal(t, id, fmt.Sprint(intID))

	strID := GetIncomingUserID[string](ctx)
	xtest.Equal(t, id, strID)
}
