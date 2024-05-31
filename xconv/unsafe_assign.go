package xconv

import (
	"encoding"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"

	"go.olapie.com/x/xreflect"
)

type FieldNameMatcher interface {
	MatchFieldName(srcName, dstName string) bool
}

type fieldNameEqual struct {
}

func (f fieldNameEqual) MatchFieldName(srcName, dstName string) bool {
	return srcName == dstName
}

type UnsafeAssignOptions struct {
	FieldNameMatcher FieldNameMatcher
}

// UnsafeAssign fill src underlying value and fields with dst
func UnsafeAssign(dst any, src any, optFns ...func(options *UnsafeAssignOptions)) error {
	options := &UnsafeAssignOptions{}
	for _, fn := range optFns {
		fn(options)
	}

	if dst == nil {
		return errors.New("dst is nil")
	}

	if src == nil {
		return errors.New("src is nil")
	}

	if options.FieldNameMatcher == nil {
		options.FieldNameMatcher = fieldNameEqual{}
	}

	func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()
		_ = JSONCopy(dst, src)
		_ = GobCopy(dst, src)
	}()

	dv := xreflect.IndirectWritableValue(reflect.ValueOf(dst), false)
	// dv must be a nil pointer or a valid value
	err := unsafeAssign(dv, reflect.ValueOf(src), options)
	if err != nil {
		return fmt.Errorf("cannot assign %T to %T: %w", src, dv.Interface(), err)
	}
	return Validate(dst)
}

// dst is valid value or pointer to value
func unsafeAssign(dst reflect.Value, src reflect.Value, options *UnsafeAssignOptions) error {
	src = xreflect.IndirectReadableValue(src)
	dv := xreflect.IndirectWritableValue(dst, true)
	switch dv.Kind() {
	case reflect.Bool:
		b, err := ToBool(src.Interface())
		if err != nil {
			return fmt.Errorf("parse bool: %w", err)
		}
		dv.SetBool(b)
	case reflect.String:
		s, err := ToString(src.Interface())
		if err != nil {
			return fmt.Errorf("parse string: %w", err)
		}
		dv.SetString(s)
	case reflect.Slice:
		if src.Kind() != reflect.Slice {
			return errors.New("source value is not slice")
		}
		l := reflect.MakeSlice(dv.Type(), src.Len(), src.Cap())
		for i := 0; i < src.Len(); i++ {
			err := unsafeAssign(l.Index(i), src.Index(i), options)
			if err != nil {
				return fmt.Errorf("cannot assign [%d]: %w", i, err)
			}
		}
		dv.Set(l)
	case reflect.Map:
		if src.Kind() != reflect.Map {
			return fmt.Errorf("cannot assign %v to map", src.Kind())
		}
		err := mapToMap(dv, src, options)
		if err != nil {
			return fmt.Errorf("mapToMap: %w", err)
		}
		return nil
	case reflect.Struct:
		err := valueToStruct(dv, src, options)
		if err != nil {
			return fmt.Errorf("valueToStruct: %w", err)
		}
		return nil
	case reflect.Interface:
		// if i is a pointer to an interface, then ValueOf(i).Elem().Kind() is reflect.Interface
		pv := reflect.New(dv.Elem().Type())
		if err := unsafeAssign(pv.Elem(), src, options); err != nil {
			return fmt.Errorf("cannot assign to interface(%v): %w", dv.Elem().Kind(), err)
		}
		dv.Set(pv.Elem())
	default:
		if xreflect.IsIntValue(dv) {
			i, err := ToInt64(src.Interface())
			if err != nil {
				return fmt.Errorf("parse int64: %w", err)
			}
			dv.SetInt(i)
		} else if xreflect.IsUintValue(dv) {
			i, err := ToUint64(src.Interface())
			if err != nil {
				return fmt.Errorf("parse uint64: %w", err)
			}
			dv.SetUint(i)
		} else if xreflect.IsFloatValue(dv) {
			i, err := ToFloat64(src.Interface())
			if err != nil {
				return fmt.Errorf("parse float64: %w", err)
			}
			dv.SetFloat(i)
		} else {
			return fmt.Errorf("unknown kind=%v", dv.Kind())
		}
	}
	return nil
}

func valueToStruct(dst reflect.Value, src reflect.Value, options *UnsafeAssignOptions) error {
	switch k := src.Kind(); k {
	case reflect.Map:
		err := mapToStruct(dst, src, options)
		if err != nil {
			return fmt.Errorf("mapToStruct: %w", err)
		}
		return nil
	case reflect.Struct:
		err := structToStruct(dst, src, options)
		if err != nil {
			return fmt.Errorf("structToStruct: %w", err)
		}
		return nil
	case reflect.String:
		if dst.CanInterface() {
			if u, ok := dst.Interface().(encoding.TextUnmarshaler); ok && u != nil {
				err := u.UnmarshalText([]byte(src.String()))
				if err != nil {
					return fmt.Errorf("cannot unmarshal text into %v: %w", dst.Type(), err)
				}
				return nil
			}

			if dst.CanAddr() && dst.Addr().CanInterface() {
				if u, ok := dst.Addr().Interface().(encoding.TextUnmarshaler); ok && u != nil {
					err := u.UnmarshalText([]byte(src.String()))
					if err != nil {
						return fmt.Errorf("cannot unmarshal text into pointer to %v: %w", dst.Type(), err)
					}
					return nil
				}
			}
		}
		return fmt.Errorf("src is %v instead of struct or map", k)
	case reflect.Ptr:
		if src.IsNil() {
			return nil
		}
		return valueToStruct(dst, src.Elem(), options)
	default:
		return fmt.Errorf("src is %v instead of struct or map", k)
	}
}

func mapToMap(dst reflect.Value, src reflect.Value, options *UnsafeAssignOptions) error {
	if !src.Type().Key().AssignableTo(dst.Type().Key()) {
		if dst.CanAddr() {
			if addr, ok := dst.Addr().Interface().(json.Unmarshaler); ok {
				if dst.IsNil() {
					dst.Set(reflect.MakeMap(dst.Type()))
				}
				err := JSONCopy(addr, src.Interface())
				if err != nil {
					return fmt.Errorf("json copy: %w", err)
				}
				return nil
			}

			if addr, ok := dst.Addr().Interface().(gob.GobDecoder); ok {
				if dst.IsNil() {
					dst.Set(reflect.MakeMap(dst.Type()))
				}
				err := GobCopy(addr, src.Interface())
				if err != nil {
					return fmt.Errorf("gob copy: %w", err)
				}
				return nil
			}
		}
		return fmt.Errorf("cannot assign %s to %s", src.Type().Key(), dst.Type().Key())
	}

	if dst.IsNil() {
		dst.Set(reflect.MakeMap(dst.Type()))
	}

	de := dst.Type().Elem()
	canAssign := src.Type().Elem().AssignableTo(de)
	for _, k := range src.MapKeys() {
		switch {
		case canAssign:
			dst.SetMapIndex(k, src.MapIndex(k))
		case de.Kind() == reflect.Ptr:
			kv := reflect.New(de.Elem())
			err := unsafeAssign(kv, src.MapIndex(k), options)
			if err != nil {
				log.Printf("Cannot unsafeAssign: %v", k.Interface())
				break
			}
			dst.SetMapIndex(k, kv)
		default:
			kv := reflect.New(de)
			err := unsafeAssign(kv, src.MapIndex(k), options)
			if err != nil {
				log.Printf("Cannot unsafeAssign: %v", k.Interface())
				break
			}
			dst.SetMapIndex(k, kv.Elem())
		}
	}
	return nil
}

func mapToStruct(dst reflect.Value, src reflect.Value, options *UnsafeAssignOptions) error {
	if k := src.Type().Key().Kind(); k != reflect.String {
		return fmt.Errorf("src key is %s intead of string", k)
	}

	for i := 0; i < dst.NumField(); i++ {
		fv := dst.Field(i)
		if fv.IsValid() == false || fv.CanSet() == false {
			continue
		}

		ft := dst.Type().Field(i)
		if ft.Anonymous {
			err := unsafeAssign(fv, src, options)
			if err != nil {
				log.Printf("Cannot assign %s: %v", ft.Name, err)
			}
			continue
		}

		for _, key := range src.MapKeys() {
			if !options.FieldNameMatcher.MatchFieldName(key.String(), ft.Name) {
				continue
			}

			fsv := src.MapIndex(key)
			if !fsv.IsValid() {
				log.Printf("Invalid value for %s", ft.Name)
				continue
			}

			if fsv.Interface() == nil {
				continue
			}

			err := unsafeAssign(fv, reflect.ValueOf(fsv.Interface()), options)
			if err != nil {
				log.Printf("Cannot assign %s: %v", ft.Name, err)
			}
			break
		}
	}
	return nil
}

func structToStruct(dst reflect.Value, src reflect.Value, options *UnsafeAssignOptions) error {
	for i := 0; i < dst.NumField(); i++ {
		fv := dst.Field(i)
		if fv.IsValid() == false || fv.CanSet() == false {
			continue
		}

		ft := dst.Type().Field(i)
		if ft.Anonymous {
			if err := unsafeAssign(fv, src, options); err != nil {
				log.Printf("Cannot assign anonymous %s: %v", ft.Name, err)
			}
			continue
		}

		for i := 0; i < src.NumField(); i++ {
			sfv := src.Field(i)
			sfName := src.Type().Field(i).Name
			if !sfv.IsValid() || sfv.Interface() == nil {
				continue
			}

			if !xreflect.IsExported(sfName) || !options.FieldNameMatcher.MatchFieldName(sfName, ft.Name) {
				continue
			}

			err := unsafeAssign(fv, reflect.ValueOf(sfv.Interface()), options)
			if err != nil {
				log.Printf("Cannot assign %s to %s: %v", ft.Name, sfName, err)
			}
			break
		}
	}

	for i := 0; i < src.NumField(); i++ {
		sfv := src.Field(i)
		sfName := src.Type().Field(i).Name
		if !sfv.IsValid() || (sfv.CanInterface() && sfv.Interface() == nil) || sfv.IsZero() || !xreflect.IsExported(sfName) {
			continue
		}

		if src.Type().Field(i).Anonymous {
			_ = unsafeAssign(dst, reflect.ValueOf(sfv.Interface()), options)
		}
	}
	return nil
}

type FieldValidator interface {
	Validate() error
}

func Validate(i any) error {
	if v, ok := i.(FieldValidator); ok {
		return v.Validate()
	}

	v := reflect.ValueOf(i)
	if v.IsValid() && (v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface) && !v.IsNil() {
		v = v.Elem()
		if v.CanInterface() {
			if va, ok := v.Interface().(FieldValidator); ok {
				return va.Validate()
			}
		}
	}

	v = xreflect.IndirectReadableValue(v)
	if v.Kind() == reflect.Struct {
		t := v.Type()
		for j := 0; j < v.NumField(); j++ {
			if !xreflect.IsExported(t.Field(j).Name) {
				continue
			}
			if err := Validate(v.Field(j).Interface()); err != nil {
				return fmt.Errorf("%s:%w", t.Field(j).Name, err)
			}
		}
	}
	return nil
}

func UnsafeSetBytes(target any, b []byte) error {
	if tu, ok := target.(encoding.TextUnmarshaler); ok {
		err := tu.UnmarshalText(b)
		if err != nil {
			return fmt.Errorf("unmarshal text: %w", err)
		}
		return nil
	}

	if bu, ok := target.(encoding.BinaryUnmarshaler); ok {
		err := bu.UnmarshalBinary(b)
		if err != nil {
			return fmt.Errorf("unmarshal binary: %w", err)
		}
		return nil
	}

	if ju, ok := target.(json.Unmarshaler); ok {
		err := ju.UnmarshalJSON(b)
		if err != nil {
			return fmt.Errorf("unmarshal json: %w", err)
		}
		return nil
	}

	v := xreflect.IndirectReadableValue(reflect.ValueOf(target))
	if !v.CanSet() {
		return fmt.Errorf("cannot set value: %T", target)
	}
	if xreflect.IsIntValue(v) {
		i, err := ToInt64(b)
		if err != nil {
			return fmt.Errorf("parse int: %v", err)
		}
		v.SetInt(i)
	}

	if xreflect.IsUintValue(v) {
		i, err := ToUint64(b)
		if err != nil {
			return fmt.Errorf("parse uint: %w", err)
		}
		v.SetUint(i)
	}

	if xreflect.IsFloatValue(v) {
		i, err := ToFloat64(b)
		if err != nil {
			return fmt.Errorf("parse float: %w", err)
		}
		v.SetFloat(i)
	}

	switch v.Kind() {
	case reflect.String:
		v.SetString(string(b))
	case reflect.Bool:
		i, err := ToBool(b)
		if err != nil {
			return fmt.Errorf("parse bool: %w", err)
		}
		v.SetBool(i)
	case reflect.Array:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			v.SetBytes(b)
		} else {
			return fmt.Errorf("cannot assign %T", target)
		}
	default:
		return fmt.Errorf("cannot assign %T", target)
	}
	return nil
}
