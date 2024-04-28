package xbase36

import (
	"testing"

	"github.com/google/uuid"
	"go.olapie.com/x/xtest"
)

func TestEncodeToString(t *testing.T) {
	t.Run("Decode", func(t *testing.T) {
		for range 20 {
			id := uuid.New()
			idStr := EncodeToString(id[:])
			parsed, err := DecodeString(idStr)
			xtest.NoError(t, err)
			xtest.Equal(t, id[:], parsed)
		}
	})

	t.Run("Parse", func(t *testing.T) {
		for range 20 {
			idStr := NewUUIDString()
			_, err := UUIDFromString(idStr)
			xtest.NoError(t, err)
		}
	})
}
