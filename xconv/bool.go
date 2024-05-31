package xconv

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	"go.olapie.com/x/xreflect"
)

// ToBool converts i to bool
// i can be bool, integer or string
func ToBool(i any) (bool, error) {
	i = xreflect.Indirect(i)
	switch v := i.(type) {
	case bool:
		return v, nil
	case nil:
		return false, strconv.ErrSyntax
	case string:
		return strconv.ParseBool(v)
	}

	if b, ok := i.([]byte); ok {
		i = string(b)
	}
	switch v := reflect.ValueOf(i); v.Kind() {
	case reflect.Bool:
		return v.Bool(), nil
	case reflect.String:
		return strconv.ParseBool(v.String())
	}

	n, err := parseInt64(i)
	if err != nil {
		return false, fmt.Errorf("cannot convert %#v of type %T to bool", i, i)
	}
	return n != 0, nil
}

// ToBoolSlice converts i to []bool
// i is an array or slice with elements convertible to bool
func ToBoolSlice(i any) ([]bool, error) {
	i = xreflect.Indirect(i)
	if i == nil {
		return nil, nil
	}
	if l, ok := i.([]bool); ok {
		return l, nil
	}
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, fmt.Errorf("cannot convert %#v of type %T to []bool", i, i)
	}
	num := v.Len()
	res := make([]bool, num)
	var err error
	for j := 0; j < num; j++ {
		res[j], err = ToBool(v.Index(j).Interface())
		if err != nil {
			return nil, fmt.Errorf("convert index %d: %w", i, err)
		}
	}
	return res, nil
}

// MustToBoolSlice converts i to []bool, will panic if failed
func MustToBoolSlice(i any) []bool {
	v, err := ToBoolSlice(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

// MustToBool converts i to bool, will panic if failed
func MustToBool(i any) bool {
	v, err := ToBool(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}
