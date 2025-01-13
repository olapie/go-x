package xconv

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"go.olapie.com/x/xreflect"
)

func ToFloat32(i any) (float32, error) {
	v, err := ToFloat64(i)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %#v of type %T to float32", i, i)
	}
	if v > math.MaxFloat32 || v < -math.MaxFloat32 {
		return 0, strconv.ErrRange
	}
	return float32(v), nil
}

func ToFloat64(i any) (float64, error) {
	i = xreflect.Indirect(i)
	if i == nil {
		return 0, strconv.ErrSyntax
	}

	if b, ok := i.([]byte); ok {
		i = string(b)
	}
	v := reflect.ValueOf(i)
	if xreflect.IsIntValue(v) {
		return float64(v.Int()), nil
	}

	if xreflect.IsUintValue(v) {
		return float64(v.Uint()), nil
	}

	if xreflect.IsFloatValue(v) {
		return v.Float(), nil
	}

	switch v.Kind() {
	case reflect.String:
		return strconv.ParseFloat(v.String(), 64)
	case reflect.Bool:
		if v.Bool() {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %#v of type %T to float64", i, i)
	}
}

func ToFloat32Slice(i any) ([]float32, error) {
	i = xreflect.Indirect(i)
	if i == nil {
		return nil, nil
	}
	if l, ok := i.([]float32); ok {
		return l, nil
	}
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, fmt.Errorf("cannot convert %#v of type %T to []float32", i, i)
	}
	num := v.Len()
	res := make([]float32, num)
	var err error
	for j := 0; j < num; j++ {
		res[j], err = ToFloat32(v.Index(j).Interface())
		if err != nil {
			return nil, fmt.Errorf("convert index %d: %w", i, err)
		}
	}
	return res, nil
}

func ToFloat64Slice(i any) ([]float64, error) {
	i = xreflect.Indirect(i)
	if i == nil {
		return nil, nil
	}
	if l, ok := i.([]float64); ok {
		return l, nil
	}
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, fmt.Errorf("cannot convert %#v of type %T to []float64", i, i)
	}
	num := v.Len()
	res := make([]float64, num)
	var err error
	for j := 0; j < num; j++ {
		res[j], err = ToFloat64(v.Index(j).Interface())
		if err != nil {
			return nil, fmt.Errorf("convert index %d: %w", i, err)
		}
	}
	return res, nil
}
