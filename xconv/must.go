package xconv

import (
	"fmt"
)

// MustGet eliminates nil err and panics if err isn't nil
func MustGet[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetTwo eliminates nil err and panics if err isn't nil
func MustGetTwo[T1 any, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err)
	}
	return v1, v2
}

// MustTrue panics if b is not true
func MustTrue(b bool, msgAndArgs ...any) {
	if !b {
		if len(msgAndArgs) == 0 {
			panic("")
		}
		panic(fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...))
	}
}

// MustFalse panics if b is not true
func MustFalse(b bool, msgAndArgs ...any) {
	if b {
		if len(msgAndArgs) == 0 {
			panic("")
		}
		panic(fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...))
	}
}

// MustError panics if v is not nil
func MustError(v error, msgAndArgs ...any) {
	if v == nil {
		if len(msgAndArgs) == 0 {
			panic("")
		}
		panic(fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...))
	}
}

// MustNoError panics if v is nil
func MustNoError(v error, msgAndArgs ...any) {
	if v == nil {
		return
	}

	s := fmt.Sprintf("%v", v)
	if len(msgAndArgs) == 0 {
		panic(s)
	}
	format := s + " " + fmt.Sprint(msgAndArgs[0])
	panic(fmt.Sprintf(format, msgAndArgs[1:]...))
}

// MustNilPtr panics if v is not nil
func MustNilPtr[T any](v *T, msgAndArgs ...any) {
	if v == nil {
		return
	}

	s := fmt.Sprintf("%v", v)
	if len(msgAndArgs) == 0 {
		panic(s)
	}
	format := s + " " + fmt.Sprint(msgAndArgs[0])
	panic(fmt.Sprintf(format, msgAndArgs[1:]...))
}

// MustNotNilPtr panics if v is nil
func MustNotNilPtr[T any](v *T, msgAndArgs ...any) {
	if v == nil {
		if len(msgAndArgs) == 0 {
			panic("")
		}
		panic(fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...))
	}
}

// MustNil panics if v is not nil
func MustNil[T any, V []T | *T | map[string]T | map[int]T | map[int64]T](v V, msgAndArgs ...any) {
	if v == nil {
		return
	}

	s := fmt.Sprintf("%v", v)
	if len(msgAndArgs) == 0 {
		panic(s)
	}
	format := s + " " + fmt.Sprint(msgAndArgs[0])
	panic(fmt.Sprintf(format, msgAndArgs[1:]...))
}

// MustNotNil panics if v is nil
func MustNotNil[T any, V []T | *T | map[string]T | map[int]T | map[int64]T](v V, msgAndArgs ...any) {
	if v == nil {
		if len(msgAndArgs) == 0 {
			panic("")
		}
		panic(fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...))
	}
}
