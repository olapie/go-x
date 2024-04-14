package xsecurity

import (
	"bytes"
	"io"
	"testing"
	"time"

	"go.olapie.com/x/xtest"
)

func TestDecryptedReader(t *testing.T) {
	raw := []byte(SHA1(time.Now().String()))
	enc, err := Encrypt(raw, "123")
	xtest.NoError(t, err)
	r := NewDecryptedReader(bytes.NewReader(enc), "123")
	dec := &bytes.Buffer{}
	n, err := io.Copy(dec, r)
	xtest.NoError(t, err)
	t.Log(n)
	xtest.Equal(t, raw, dec.Bytes())
}

func BenchmarkDecryptedReader(b *testing.B) {
	raw := xtest.RandomBytes(int(4 * (1 << 20)))
	enc, err := Encrypt(raw, "123")
	xtest.NoError(b, err)
	for i := 0; i < b.N; i++ {
		Decrypt(enc, "123")
	}
}
