package xconv

import (
	"encoding"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"go.olapie.com/x/xreflect"
)

func ToBytes(i any) ([]byte, error) {
	i = xreflect.Indirect(i)
	switch v := i.(type) {
	case []byte:
		return v, nil
	case nil:
		return nil, strconv.ErrSyntax
	case string:
		return []byte(v), nil
	}

	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.Uint8 {
		return v.Bytes(), nil
	}
	return nil, fmt.Errorf("cannot convert %#v of type %T to []byte", i, i)
}

func ToByteArray8[T []byte | string](v T) [8]byte {
	if len(v) > 8 {
		panic("cannot convert into [8]byte")
	}
	var a [8]byte
	copy(a[:], v[:])
	return a
}

func ToByteArray16[T []byte | string](v T) [16]byte {
	if len(v) > 16 {
		panic("cannot convert into [16]byte")
	}
	var a [16]byte
	copy(a[:], v[:])
	return a
}

func ToByteArray32[T []byte | string](v T) [32]byte {
	if len(v) > 32 {
		panic("cannot convert into [32]byte")
	}
	var a [32]byte
	copy(a[:], v[:])
	return a
}

func ToByteArray64[T []byte | string](v T) [64]byte {
	if len(v) > 64 {
		panic("cannot convert into [64]byte")
	}
	var a [64]byte
	copy(a[:], v[:])
	return a
}

func UnsafeMarshal(i any) ([]byte, error) {
	if data, ok := i.([]byte); ok {
		return data, nil
	}

	var data []byte
	srcType := reflect.TypeOf(i)
	dstType := reflect.TypeOf(data)
	if srcType.AssignableTo(dstType) {
		reflect.ValueOf(&data).Elem().Set(reflect.ValueOf(i))
		return data, nil
	}

	if srcType.ConvertibleTo(dstType) {
		reflect.ValueOf(&data).Elem().Set(reflect.ValueOf(i).Convert(dstType))
		return data, nil
	}

	if m, ok := i.(encoding.BinaryMarshaler); ok {
		return m.MarshalBinary()
	}

	if m, ok := i.(encoding.TextMarshaler); ok {
		return m.MarshalText()
	}

	if m, ok := i.(json.Marshaler); ok {
		return m.MarshalJSON()
	}

	if m, ok := i.(gob.GobEncoder); ok {
		return m.GobEncode()
	}

	return nil, errors.New("cannot convert ")
}

func UnsafeUnmarshal(data []byte, i any) (err error) {
	if reflect.ValueOf(i).Kind() != reflect.Pointer {
		return fmt.Errorf("cannot unmarshal to non pointer type: %T", i)
	}

	if p, ok := i.(*[]byte); ok {
		*p = data
		return nil
	}

	srcType := reflect.TypeOf(data)
	dstType := reflect.TypeOf(i).Elem()
	if srcType.AssignableTo(dstType) {
		reflect.ValueOf(i).Elem().Set(reflect.ValueOf(data))
		return nil
	}

	if srcType.ConvertibleTo(dstType) {
		reflect.ValueOf(i).Elem().Set(reflect.ValueOf(data).Convert(dstType))
		return nil
	}

	// i is a pointer
	// v is pointer of the same type
	v := xreflect.DeepNew(reflect.TypeOf(i).Elem())
	defer func() {
		if err == nil {
			// assign v to i
			// i is a parameter, it cannot be set, the value it points to can be set
			// assign v.Elem to i.Elem
			reflect.ValueOf(i).Elem().Set(v.Elem())
		}
	}()

	for p := v; p.Kind() == reflect.Pointer && !p.IsNil(); p = p.Elem() {
		if u, ok := p.Interface().(encoding.BinaryUnmarshaler); ok {
			return u.UnmarshalBinary(data)
		}

		if u, ok := p.Interface().(encoding.TextUnmarshaler); ok {
			return u.UnmarshalText(data)
		}

		if u, ok := p.Interface().(json.Unmarshaler); ok {
			return u.UnmarshalJSON(data)
		}

		if d, ok := p.Interface().(gob.GobDecoder); ok {
			return d.GobDecode(data)
		}
	}

	return fmt.Errorf("cannot unmarshal into: %T", i)
}
