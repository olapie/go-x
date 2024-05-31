package xerror

import (
	"fmt"
	"reflect"

	"go.olapie.com/x/xreflect"
)

type String string

func (s String) Error() string {
	return string(s)
}

func Wrap(err error, format string, a ...any) error {
	if err == nil {
		return nil
	}
	a = append(a, err)
	return fmt.Errorf(format+":%w", a...)
}

// Cause returns the root cause error
func Cause(err error) error {
	for {
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			break
		}
		err = u.Unwrap()
	}
	return err
}

func GetCode(err error) int {
	err = Cause(err)
	if err == nil {
		return 0
	}

	if s, ok := err.(interface{ Code() int }); ok {
		return s.Code()
	}

	if s, ok := err.(interface{ GetCode() int }); ok {
		return s.GetCode()
	}

	// ------------
	// int32

	if s, ok := err.(interface{ Code() int32 }); ok {
		return int(s.Code())
	}

	if s, ok := err.(interface{ GetCode() int32 }); ok {
		return int(s.GetCode())
	}

	if s, ok := err.(interface{ StatusCode() int }); ok {
		return s.StatusCode()
	}

	if s, ok := err.(interface{ GetStatusCode() int }); ok {
		return s.GetStatusCode()
	}

	if s, ok := err.(interface{ Status() int }); ok {
		return s.Status()
	}

	if s, ok := err.(interface{ GetStatus() int }); ok {
		return s.GetStatus()
	}

	// ------------
	// int32

	if s, ok := err.(interface{ StatusCode() int32 }); ok {
		return int(s.StatusCode())
	}

	if s, ok := err.(interface{ GetStatusCode() int32 }); ok {
		return int(s.GetStatusCode())
	}

	if s, ok := err.(interface{ Status() int32 }); ok {
		return int(s.Status())
	}

	if s, ok := err.(interface{ GetStatus() int32 }); ok {
		return int(s.GetStatus())
	}

	v := reflect.ValueOf(xreflect.Indirect(err))
	t := v.Type()
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			ft := t.Field(i)
			switch ft.Name {
			case "Code", "Status", "StatusCode", "ErrorCode":
				fv := v.Field(i)
				if fv.CanInt() {
					return int(fv.Int())
				}

				if fv.CanUint() {
					return int(fv.Uint())
				}

				return 0
			}
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			if k.Kind() != reflect.String {
				continue
			}
			switch k.String() {
			case "Code", "Status", "StatusCode", "ErrorCode":
				vv := v.MapIndex(k)
				if vv.CanInt() {
					return int(vv.Int())
				}

				if vv.CanUint() {
					return int(vv.Uint())
				}

				return 0
			}
		}
	default:
		break
	}
	return 0
}
