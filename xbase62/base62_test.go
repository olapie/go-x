package xbase62

import (
	"math/rand"
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
