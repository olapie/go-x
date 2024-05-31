package xconv

import (
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"unsafe"

	"go.olapie.com/x/xreflect"
)

const (
	// MaxInt represents maximum int
	MaxInt = 1<<(8*unsafe.Sizeof(int(0))-1) - 1
	// MinInt represents minimum int
	MinInt = -1 << (8*unsafe.Sizeof(int(0)) - 1)
	// MaxUint represents maximum uint
	MaxUint = 1<<(8*unsafe.Sizeof(uint(0))) - 1
)

// ToInt converts i to int
func ToInt(i any) (int, error) {
	n, err := parseInt64(i)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %#v of type %T to int64", i, i)
	}
	if n > MaxInt || n < MinInt {
		return 0, strconv.ErrRange
	}
	return int(n), nil
}

// ToInt8 converts i to int8
func ToInt8(i any) (int8, error) {
	n, err := parseInt64(i)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %#v of type %T to int8", i, i)
	}
	if n > math.MaxInt8 || n < math.MinInt8 {
		return 0, strconv.ErrRange
	}
	return int8(n), nil
}

// ToInt16 converts i to int16
func ToInt16(i any) (int16, error) {
	n, err := parseInt64(i)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %#v of type %T to int16", i, i)
	}
	if n > math.MaxInt16 || n < math.MinInt16 {
		return 0, strconv.ErrRange
	}
	return int16(n), nil
}

func ToInt32(i any) (int32, error) {
	n, err := parseInt64(i)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %#v of type %T to int32", i, i)
	}
	if n > math.MaxInt32 || n < math.MinInt32 {
		return 0, strconv.ErrRange
	}
	return int32(n), nil
}

func MustInt32(i any) int32 {
	v, err := ToInt32(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func ToInt64(i any) (int64, error) {
	return parseInt64(i)
}

func ToUint(i any) (uint, error) {
	n, err := parseUint64(i)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %#v of type %T to uint", i, i)
	}
	if n > MaxUint {
		return 0, strconv.ErrRange
	}
	return uint(n), nil
}

func ToUint8(i any) (uint8, error) {
	n, err := parseUint64(i)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %#v of type %T to uint8", i, i)
	}
	if n > math.MaxUint8 {
		return 0, strconv.ErrRange
	}
	return uint8(n), nil
}

func ToUint16(i any) (uint16, error) {
	n, err := parseUint64(i)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %#v of type %T to uint16", i, i)
	}
	if n > math.MaxUint16 {
		return 0, strconv.ErrRange
	}
	return uint16(n), nil
}

func ToUint32(i any) (uint32, error) {
	n, err := parseUint64(i)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %#v of type %T to uint32", i, i)
	}
	if n > math.MaxUint32 {
		return 0, strconv.ErrRange
	}
	return uint32(n), nil
}

func ToUint64(i any) (uint64, error) {
	return parseUint64(i)
}

func ToIntSlice(i any) ([]int, error) {
	i = xreflect.Indirect(i)
	if i == nil {
		return nil, nil
	}
	if l, ok := i.([]int); ok {
		return l, nil
	}
	switch v := reflect.ValueOf(i); v.Kind() {
	case reflect.Slice, reflect.Array:
		num := v.Len()
		res := make([]int, num)
		var err error
		for j := 0; j < num; j++ {
			res[j], err = ToInt(v.Index(j).Interface())
			if err != nil {
				return nil, fmt.Errorf("convert index %d: %w", j, err)
			}
		}
		return res, nil
	default:
		if k, err := ToInt(i); err == nil {
			return []int{k}, nil
		}
		return nil, fmt.Errorf("cannot convert %v to slice", v.Kind())
	}
}

func ToInt64Slice(i any) ([]int64, error) {
	i = xreflect.Indirect(i)
	if i == nil {
		return nil, nil
	}
	if l, ok := i.([]int64); ok {
		return l, nil
	}

	switch v := reflect.ValueOf(i); v.Kind() {
	case reflect.Slice, reflect.Array:
		num := v.Len()
		res := make([]int64, num)
		var err error
		for j := 0; j < num; j++ {
			res[j], err = parseInt64(v.Index(j).Interface())
			if err != nil {
				return nil, fmt.Errorf("convert element at index %d: %w", i, err)
			}
		}
		return res, nil
	default:
		if k, err := ToInt64(i); err == nil {
			return []int64{k}, nil
		}
		return nil, fmt.Errorf("cannot convert %#v of type %T to []int64", i, i)
	}
}

func ToUintSlice(i any) ([]uint, error) {
	i = xreflect.Indirect(i)
	if i == nil {
		return nil, nil
	}
	if l, ok := i.([]uint); ok {
		return l, nil
	}

	switch v := reflect.ValueOf(i); v.Kind() {
	case reflect.Slice, reflect.Array:
		num := v.Len()
		res := make([]uint, num)
		var err error
		for j := 0; j < num; j++ {
			res[j], err = ToUint(v.Index(j).Interface())
			if err != nil {
				return nil, fmt.Errorf("convert element at index %d: %w", i, err)
			}
		}
		return res, nil
	default:
		if ui, err := ToUint(i); err == nil {
			return []uint{ui}, nil
		}
		return nil, fmt.Errorf("cannot convert %#v of type %T to []uint", i, i)
	}
}

func ToUint64Slice(i any) ([]uint64, error) {
	i = xreflect.Indirect(i)
	if i == nil {
		return nil, nil
	}
	if l, ok := i.([]uint64); ok {
		return l, nil
	}

	switch v := reflect.ValueOf(i); v.Kind() {
	case reflect.Slice, reflect.Array:
		num := v.Len()
		res := make([]uint64, num)
		var err error
		for j := 0; j < num; j++ {
			res[j], err = parseUint64(v.Index(j).Interface())
			if err != nil {
				return nil, fmt.Errorf("convert element at index %d: %w", i, err)
			}
		}
		return res, nil
	default:
		if ui, err := ToUint64(i); err == nil {
			return []uint64{ui}, nil
		}
		return nil, fmt.Errorf("cannot convert %#v of type %T to []uint64", i, i)
	}
}

// MustToInt panics if ToInt(i) failed
func MustToInt(i any) int {
	v, err := ToInt(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

// MustToInt8 panics if ToInt8(i) failed
func MustToInt8(i any) int8 {
	v, err := ToInt8(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

// MustToInt16 panics if ToInt16(i) failed
func MustToInt16(i any) int16 {
	v, err := ToInt16(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func MustToInt64(i any) int64 {
	v, err := ToInt64(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}
func MustToUint(i any) uint {
	v, err := ToUint(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}
func MustToUint8(i any) uint8 {
	v, err := ToUint8(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func MustToUint16(i any) uint16 {
	v, err := ToUint16(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func MustToUint32(i any) uint32 {
	v, err := ToUint32(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func MustToUint64(i any) uint64 {
	v, err := ToUint64(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func MustToIntSlice(i any) []int {
	v, err := ToIntSlice(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func MustToInt64Slice(i any) []int64 {
	v, err := ToInt64Slice(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func MustToUintSlice(i any) []uint {
	v, err := ToUintSlice(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}
func MustToUint64Slice(i any) []uint64 {
	v, err := ToUint64Slice(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func MustToString(i any) string {
	v, err := ToString(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func MustToStringSlice(i any) []string {
	v, err := ToStringSlice(i)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func CastToIntSlice[T ~int | ~int32 | ~int16 | ~int8](a []T) []int {
	res := make([]int, len(a))
	for i, v := range a {
		res[i] = int(v)
	}
	return res
}

func CastFromIntSlice[T ~int | ~int32 | ~int16 | ~int8](a []int) []T {
	res := make([]T, len(a))
	for i, v := range a {
		res[i] = T(v)
	}
	return res
}

func CastToInt16Slice[T ~int16 | ~int8 | ~int](a []T) []int16 {
	res := make([]int16, len(a))
	for i, v := range a {
		res[i] = int16(v)
	}
	return res
}

func CastFromInt16Slice[T ~int16 | ~int8 | ~int](a []int16) []T {
	res := make([]T, len(a))
	for i, v := range a {
		res[i] = T(v)
	}
	return res
}

func CastToInt32Slice[T ~int32 | ~int16 | ~int8 | ~int](a []T) []int32 {
	res := make([]int32, len(a))
	for i, v := range a {
		res[i] = int32(v)
	}
	return res
}

func CastFromInt32Slice[T ~int32 | ~int16 | ~int8 | ~int](a []int32) []T {
	res := make([]T, len(a))
	for i, v := range a {
		res[i] = T(v)
	}
	return res
}

func CastToInt64Slice[T ~int64 | ~int16 | ~int32 | ~int | ~int8](a []T) []int64 {
	res := make([]int64, len(a))
	for i, v := range a {
		res[i] = int64(v)
	}
	return res
}

func CastFromInt64Slice[T ~int64](a []int64) []T {
	res := make([]T, len(a))
	for i, v := range a {
		res[i] = T(v)
	}
	return res
}

func CastToIntP[T ~int | ~int32](p *T) *int {
	if p == nil {
		return nil
	}
	v := int(*p)
	return &v
}

func CastFromIntP[T ~int | ~int32](p *int) *T {
	if p == nil {
		return nil
	}
	v := T(*p)
	return &v
}

func CastToInt32P[T ~int32 | ~int](p *T) *int32 {
	if p == nil {
		return nil
	}
	v := int32(*p)
	return &v
}

func CastFromInt32P[T ~int32 | ~int](p *int32) *T {
	if p == nil {
		return nil
	}
	v := T(*p)
	return &v
}

func CastToInt64P[T ~int64 | ~int](p *T) *int64 {
	if p == nil {
		return nil
	}
	v := int64(*p)
	return &v
}

func CastFromInt64P[T ~int64 | ~int](p *int64) *T {
	if p == nil {
		return nil
	}
	v := T(*p)
	return &v
}

func parseInt64(i any) (int64, error) {
	i = xreflect.Indirect(i)
	if i == nil {
		return 0, strconv.ErrSyntax
	}
	if b, ok := i.([]byte); ok {
		i = string(b)
	}
	v := reflect.ValueOf(i)
	if xreflect.IsIntValue(v) {
		return v.Int(), nil
	}

	if xreflect.IsUintValue(v) {
		n := v.Uint()
		if n > math.MaxInt64 {
			return 0, strconv.ErrRange
		}
		return int64(n), nil
	}

	if xreflect.IsFloatValue(v) {
		return int64(v.Float()), nil
	}

	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.String:
		n, err := strconv.ParseInt(v.String(), 0, 64)
		if err == nil {
			return n, nil
		}
		if errors.Is(err, strconv.ErrRange) {
			return 0, err
		}
		if f, fErr := strconv.ParseFloat(v.String(), 64); fErr == nil {
			return int64(f), nil
		}
		return 0, err
	default:
		return 0, strconv.ErrSyntax
	}
}

func parseUint64(i any) (uint64, error) {
	i = xreflect.Indirect(i)
	if i == nil {
		return 0, strconv.ErrSyntax
	}
	if b, ok := i.([]byte); ok {
		i = string(b)
	}
	v := reflect.ValueOf(i)
	if xreflect.IsIntValue(v) {
		n := v.Int()
		if n < 0 {
			return 0, strconv.ErrRange
		}
		return uint64(n), nil
	}

	if xreflect.IsUintValue(v) {
		return v.Uint(), nil
	}

	if xreflect.IsFloatValue(v) {
		f := v.Float()
		if f < 0 {
			return 0, strconv.ErrRange
		}
		return uint64(f), nil
	}

	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.String:
		n, err := strconv.ParseInt(v.String(), 0, 64)
		if err == nil {
			if n < 0 {
				return 0, strconv.ErrRange
			}
			return uint64(n), nil
		}
		if errors.Is(err, strconv.ErrRange) {
			return 0, err
		}
		if f, fErr := strconv.ParseFloat(v.String(), 64); fErr == nil {
			if f < 0 {
				return 0, strconv.ErrRange
			}
			return uint64(f), nil
		}
		return 0, err
	default:
		return 0, strconv.ErrSyntax
	}
}
