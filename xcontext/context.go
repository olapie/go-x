package xcontext

import (
	"context"
	"reflect"

	"go.olapie.com/x/xtype"
)

type activityIncomingContext struct{}
type activityOutgoingContext struct{}

func GetIncomingActivity(ctx context.Context) *Activity {
	a, _ := ctx.Value(activityIncomingContext{}).(*Activity)
	return a
}

func GetOutgoingActivity(ctx context.Context) *Activity {
	a, _ := ctx.Value(activityOutgoingContext{}).(*Activity)
	return a
}

func WithIncomingActivity(ctx context.Context, a *Activity) context.Context {
	return context.WithValue(ctx, activityIncomingContext{}, a)
}

func WithOutgoingActivity(ctx context.Context, a *Activity) context.Context {
	return context.WithValue(ctx, activityOutgoingContext{}, a)
}

func SetIncomingUserID[T xtype.UserIDTypes](ctx context.Context, id T) error {
	a := GetIncomingActivity(ctx)
	if a == nil {
		return ErrNotExist
	}
	a.SetUserID(xtype.NewUserID(id))
	return nil
}

func GetIncomingUserID[T xtype.UserIDTypes](ctx context.Context) T {
	var id T
	a := GetIncomingActivity(ctx)
	if a == nil {
		return id
	}

	if a.userID == nil {
		return id
	}

	v := a.userID.Value()
	if id, ok := v.(T); ok {
		return id
	}

	t := reflect.TypeOf(v)
	idType := reflect.TypeOf(id)
	if t.ConvertibleTo(reflect.TypeOf(id)) {
		id, _ = reflect.ValueOf(v).Convert(idType).Interface().(T)
	}

	return id
}

var systemUserID = xtype.NewUserID[string]("ola-system-user-id")

func SetSystemUser(ctx context.Context) {
	a := GetIncomingActivity(ctx)
	if a == nil {
		panic("no incoming activity")
		return
	}
	a.userID = systemUserID
}

func IsSystemUser(ctx context.Context) bool {
	a := GetIncomingActivity(ctx)
	if a == nil {
		return false
	}
	return a.userID == systemUserID
}
