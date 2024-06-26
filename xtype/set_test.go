package xtype

import (
	"encoding/json"
	"sort"
	"testing"

	"go.olapie.com/x/xtest"
)

func TestMarshalJSON(t *testing.T) {
	s1 := NewSet[int](10)
	a0 := []int{1, 2, 3, 5, 9}
	for _, v := range a0 {
		s1.Add(v)
	}
	d1, err := s1.MarshalJSON()
	xtest.NoError(t, err)
	var s2 *Set[int]
	err = json.Unmarshal(d1, &s2)
	xtest.NoError(t, err)
	a1 := s1.Slice()
	a2 := s2.Slice()
	sort.IntSlice(a1).Sort()
	sort.IntSlice(a2).Sort()
	xtest.Equal(t, a0, a1)
	xtest.Equal(t, a1, a2)
}
