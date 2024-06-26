package xtype

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sync"

	"go.olapie.com/naming"
)

var initNameToPrototypeOnce sync.Once
var nameToPrototypeMutex sync.RWMutex
var nameToPrototype map[string]reflect.Type

func doInitPrototypeOnce() {
	initNameToPrototypeOnce.Do(
		func() {
			nameToPrototype = map[string]reflect.Type{
				"int":     reflect.TypeOf(int(1)),
				"int8":    reflect.TypeOf(int8(1)),
				"int16":   reflect.TypeOf(int16(1)),
				"int32":   reflect.TypeOf(int32(1)),
				"int64":   reflect.TypeOf(int64(1)),
				"uint":    reflect.TypeOf(uint(1)),
				"uint8":   reflect.TypeOf(uint8(1)),
				"uint16":  reflect.TypeOf(uint16(1)),
				"uint32":  reflect.TypeOf(uint32(1)),
				"uint64":  reflect.TypeOf(uint64(1)),
				"float32": reflect.TypeOf(float32(1)),
				"float64": reflect.TypeOf(float64(1)),
				"bool":    reflect.TypeOf(true),
				"string":  reflect.TypeOf(""),
			}

			nameToPrototype[getAnyTypeName(&Audio{})] = reflect.TypeOf(&Audio{})
			nameToPrototype[getAnyTypeName(&Image{})] = reflect.TypeOf(&Image{})
			nameToPrototype[getAnyTypeName(&Video{})] = reflect.TypeOf(&Video{})
		})
}

type AnyType interface {
	AnyType() string
}

// RegisterAnyType bind typ with prototype
// E.g.
//
//	contents.Register("image", &contents.Image{})
func RegisterAnyType(prototype any) {
	doInitPrototypeOnce()
	name := getAnyTypeName(prototype)
	nameToPrototypeMutex.Lock()
	defer nameToPrototypeMutex.Unlock()
	if _, ok := nameToPrototype[name]; ok {
		log.Fatalf("Duplicate name %s", name)
	}

	nameToPrototype[name] = reflect.TypeOf(prototype)
}

func getAnyTypeName(prototype any) string {
	if a, ok := prototype.(AnyType); ok {
		return a.AnyType()
	}

	p := reflect.TypeOf(prototype)
	for p.Kind() == reflect.Ptr {
		p = p.Elem()
	}
	return naming.ToSnake(p.Name())
}

func getProtoType(typ string) (reflect.Type, bool) {
	nameToPrototypeMutex.RLock()
	defer nameToPrototypeMutex.RUnlock()
	if prototype, ok := nameToPrototype[typ]; ok {
		return prototype, true
	}
	return nil, false
}

type Any struct {
	val     any
	jsonStr string
}

func NewAny(v any) *Any {
	a := &Any{}
	a.SetValue(v)
	return a
}

func (a *Any) Value() any {
	return a.val
}

func (a *Any) SetValue(v any) {
	a.val = v
	a.jsonStr = ""
}

func (a *Any) JSONString() string {
	if len(a.jsonStr) == 0 {
		data, _ := json.Marshal(a)
		a.jsonStr = string(data)
	}
	return a.jsonStr
}

func (a *Any) Int64() int64 {
	v, _ := a.val.(int64)
	return v
}

func (a *Any) Float64() float64 {
	i, _ := a.val.(float64)
	return i
}

func (a *Any) Text() string {
	s, _ := a.val.(string)
	return s
}

const (
	keyAnyType = "@t"
	keyAnyVal  = "@Int"
)

type anyJsonObject struct {
	Type  string          `json:"@t"`
	Value json.RawMessage `json:"@Int"`
}

func (a *Any) UnmarshalJSON(b []byte) error {
	var obj anyJsonObject
	if err := json.Unmarshal(b, &obj); err != nil {
		return err
	}

	doInitPrototypeOnce()
	pt, found := getProtoType(obj.Type)
	if !found {
		return fmt.Errorf("unkonwn type: %s", obj.Type)
	}

	if len(obj.Value) != 0 {
		b = obj.Value
	}

	var ptrVal = reflect.New(pt)

	for val := ptrVal; val.Kind() == reflect.Ptr && val.CanSet(); val = val.Elem() {
		val.Set(reflect.New(val.Elem().Type()))
	}

	err := json.Unmarshal(b, ptrVal.Interface())
	if err != nil {
		return err
	}
	a.SetValue(ptrVal.Elem().Interface())
	return nil
}

func (a *Any) MarshalJSON() ([]byte, error) {
	if a == nil || a.val == nil {
		return json.Marshal(nil)
	}

	var m = make(map[string]any)

	t := reflect.TypeOf(a.val)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Struct || t.Kind() == reflect.Map {
		b, err := json.Marshal(a.val)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(b, &m)
		if err != nil {
			return nil, err
		}
	} else {
		m[keyAnyVal] = a.val
	}

	m[keyAnyType] = a.TypeName()
	return json.Marshal(m)
}

func (a *Any) TypeName() string {
	return getAnyTypeName(a.val)
}
