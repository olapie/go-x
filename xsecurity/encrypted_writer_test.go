package xsecurity

import (
	"bytes"
	"io"
	"testing"
	"time"

	"go.olapie.com/x/xtest"
)

func TestEncryptedWriter(t *testing.T) {
	raw := []byte(SHA1(time.Now().String()))
	enc := &bytes.Buffer{}
	w := NewEncryptedWriter(enc, "123")
	n, err := io.Copy(w, bytes.NewReader(raw))
	xtest.NoError(t, err)
	t.Log(n)
	data, err := Encrypt(raw, "123")
	xtest.NoError(t, err)
	xtest.Equal(t, enc.Bytes(), data)
}
