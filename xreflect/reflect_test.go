package xreflect

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestIndirectKind(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		k := IndirectKind(nil)
		if diff := diffKinds(reflect.Invalid, k); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("Struct", func(t *testing.T) {
		var p time.Time
		k := IndirectKind(p)
		if diff := diffKinds(reflect.Struct, k); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("PointerToStruct", func(t *testing.T) {
		var p *time.Time
		k := IndirectKind(p)
		if diff := diffKinds(reflect.Struct, k); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("PointerToPointerToStruct", func(t *testing.T) {
		var p **time.Time
		k := IndirectKind(p)
		if diff := diffKinds(reflect.Struct, k); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("Map", func(t *testing.T) {
		var p map[string]any
		k := IndirectKind(p)
		if diff := diffKinds(reflect.Map, k); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("PointerToMap", func(t *testing.T) {
		var p map[string]any
		k := IndirectKind(p)
		if k != reflect.Map {

		}
		if diff := diffKinds(reflect.Map, k); diff != "" {
			t.Fatal(diff)
		}
	})
}

func diffKinds(expected, got reflect.Kind) string {
	if expected != got {
		return fmt.Sprintf("expect %v, got %v", expected, got)
	}
	return ""
}
