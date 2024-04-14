package xbase62

import (
	"encoding/base64"
	"math/rand"
	"strings"
	"testing"

	"github.com/google/uuid"
	"go.olapie.com/x/xtest"
)

func TestEncodeToString(t *testing.T) {
	t.Run("UUID", func(t *testing.T) {
		id := uuid.New()
		idStr := EncodeToString(id[:])
		t.Log(idStr)
		t.Log(base64.StdEncoding.EncodeToString(id[:]))
		t.Log(strings.ReplaceAll(id.String(), "-", ""))
		t.Log(id.String())
		parsed, err := DecodeString(idStr)
		xtest.NoError(t, err)
		xtest.Equal(t, id[:], parsed)
	})
}

func TestAtoi(t *testing.T) {
	i := rand.Int63()
	s := Itoa(i)
	t.Log(i, s)
	got, err := Atoi(s)
	if err != nil {
		t.Fatal(err)
	}
	if got != i {
		t.Fatalf("got %d, want %d", got, i)
	}
}
