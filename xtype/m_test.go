package xtype

import (
	"testing"

	"go.olapie.com/x/xtest"
)

func TestM_AddStruct(t *testing.T) {
	m := M{}
	m["id"] = 1
	m["name"] = "Smith"
	var foo struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	foo.Name = "Mike"
	foo.Age = 19
	err := m.AddStruct(foo)
	xtest.NoError(t, err)
	xtest.Equal(t, 1, m.Int("id"))
	xtest.Equal(t, 19, m.Int("age"))
	xtest.Equal(t, "Mike", m.String("name"))
}
