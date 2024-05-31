package xconv

import (
	"container/list"
	"reflect"

	"go.olapie.com/x/xreflect"
)

// ToList creates list.List
// i can be nil, *list.List, or array/slice
func ToList(i any) *list.List {
	if i == nil {
		return list.New()
	}

	if l, ok := i.(*list.List); ok {
		return l
	}

	lt := reflect.TypeOf((*list.List)(nil))
	if it := reflect.TypeOf(i); it.ConvertibleTo(lt) {
		return reflect.ValueOf(i).Convert(lt).Interface().(*list.List)
	}

	l := list.New()
	v := reflect.ValueOf(xreflect.Indirect(i))
	if v.IsValid() && (v.Kind() == reflect.Slice || v.Kind() == reflect.Array) && !v.IsNil() {
		for j := 0; j < v.Len(); j++ {
			l.PushBack(v.Index(j).Interface())
		}
	} else {
		l.PushBack(i)
	}
	return l
}

// SliceToList converts slice to list.List
func SliceToList[E any](a []E) *list.List {
	l := list.New()
	for _, e := range a {
		l.PushBack(e)
	}
	return l
}
