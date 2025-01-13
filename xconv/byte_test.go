package xconv

import (
	"encoding/json"
	"math/rand"
	"testing"
	"time"
)

type customByteSlice []byte

func TestMarshalCustomBytesType(t *testing.T) {
	var input customByteSlice = []byte(time.Now().String())
	output, err := UnsafeMarshal(input)
	if err != nil {
		t.Fatal(err)
	}
	if diff := diffSlice(input, output); diff != "" {
		t.Fatal(diff)
	}
}

type jsonObject struct {
	ID   int64
	Text string
}

func (o *jsonObject) MarshalJSON() ([]byte, error) {
	type alias jsonObject
	obj := (*alias)(o)
	return json.Marshal(obj)
}

func (o *jsonObject) UnmarshalJSON(data []byte) error {
	type alias jsonObject
	var obj alias
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return err
	}
	*o = jsonObject(obj)
	return nil
}

func TestMarshalJSON(t *testing.T) {
	o := jsonObject{ID: rand.Int63(), Text: time.Now().String()}
	data, err := UnsafeMarshal(&o)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))

	var o2 jsonObject
	err = UnsafeUnmarshal(data, &o2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(o2.ID, o2.Text)
}
