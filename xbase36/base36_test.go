package xbase36

import (
	"testing"

	"go.olapie.com/x/xtest"
)

func TestUUIDFromString(t *testing.T) {
	for range 200 {
		s1 := NewUUIDString()
		id, err := UUIDFromString(s1)
		xtest.NoError(t, err)

		s2 := EncodeToString(id[:])
		xtest.Equal(t, s1, s2)
	}
}
