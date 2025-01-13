package xhttp

import (
	"net/http"
	"reflect"
	"testing"
)

func TestToMap(t *testing.T) {
	t.Run("HeaderToMap", func(t *testing.T) {
		h := http.Header{}
		h.Set("k1", "v1")
		h.Set("k2", "v2")
		h.Add("k2", "v22")
		m := ToMapAny(h)
		expected := map[string]any{"K1": "v1", "K2": []string{"v2", "v22"}}
		if !reflect.DeepEqual(expected, m) {
			t.Fatalf("expect %v, got %v", expected, m)
		}
	})
}
