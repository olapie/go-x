package xassert

import "fmt"

// True panics if b is not true
func True(b bool, msgAndArgs ...any) {
	if !b {
		if len(msgAndArgs) == 0 {
			panic("")
		}
		panic(fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...))
	}
}

// False panics if b is not true
func False(b bool, msgAndArgs ...any) {
	if b {
		if len(msgAndArgs) == 0 {
			panic("")
		}
		panic(fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...))
	}
}

// Error panics if v is not nil
func Error(v error, msgAndArgs ...any) {
	if v == nil {
		if len(msgAndArgs) == 0 {
			panic("")
		}
		panic(fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...))
	}
}

// NoError panics if v is nil
func NoError(v error, msgAndArgs ...any) {
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

// NilPtr panics if v is not nil
func NilPtr[T any](v *T, msgAndArgs ...any) {
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

// NotNilPtr panics if v is nil
func NotNilPtr[T any](v *T, msgAndArgs ...any) {
	if v == nil {
		if len(msgAndArgs) == 0 {
			panic("")
		}
		panic(fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...))
	}
}

// Nil panics if v is not nil
func Nil[T any, V []T | *T | map[string]T | map[int]T | map[int64]T](v V, msgAndArgs ...any) {
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

// NotNil panics if v is nil
func NotNil[T any, V []T | *T | map[string]T | map[int]T | map[int64]T](v V, msgAndArgs ...any) {
	if v == nil {
		if len(msgAndArgs) == 0 {
			panic("")
		}
		panic(fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...))
	}
}
