package xtype

import (
	"fmt"
	"reflect"
	"strconv"
)

type AuthResult struct {
	AppID  string
	UserID UserID
}

type UserIDTypes interface {
	~int64 | ~string
}

type UserID interface {
	Int() (int64, bool)
	String() string
	Value() any
}

type userIDImpl[T UserIDTypes] struct {
	id T

	intVal    *int64
	stringVal string
}

func (u *userIDImpl[T]) Value() any {
	return u.id
}

func (u *userIDImpl[T]) Int() (int64, bool) {
	if u.intVal == nil {
		return 0, false
	}
	return *u.intVal, true
}

func (u *userIDImpl[T]) String() string {
	return u.stringVal
}

func NewUserID[T ~int64 | ~string](id T) UserID {
	uid := &userIDImpl[T]{id: id}
	v := reflect.ValueOf(id)

	if v.Type().ConvertibleTo(reflect.TypeFor[int64]()) {
		intVal := v.Int()
		uid.intVal = &intVal
		uid.stringVal = fmt.Sprint(intVal)
	} else {
		uid.stringVal = v.String()
		i, err := strconv.ParseInt(uid.stringVal, 10, 64)
		if err == nil {
			uid.intVal = &i
		}
	}
	return uid
}
