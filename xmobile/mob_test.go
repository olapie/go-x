package xmobile

import (
	"encoding/json"
	"testing"

	"go.olapie.com/x/xtest"
)

func TestMarshalStringList(t *testing.T) {
	l := new(StringList)
	l.Add("a")
	l.Add("b")
	data, err := json.Marshal(l)
	xtest.NoError(t, err)
	t.Log(string(data))

	var l2 StringList
	err = json.Unmarshal(data, &l2)
	xtest.NoError(t, err)
	t.Log(l2.Len())
}
