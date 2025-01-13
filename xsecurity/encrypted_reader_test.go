package xsecurity

import (
	"bytes"
	"io"
	"testing"
	"time"

	"go.olapie.com/x/xtest"
)

func TestEncryptedReader(t *testing.T) {
	raw := []byte(SHA1(time.Now().String()))
	enc := &bytes.Buffer{}
	{
		r := NewEncryptedReader(bytes.NewReader(raw), "123")
		n, err := io.Copy(enc, r)
		xtest.NoError(t, err)
		t.Log(n)
	}
	{
		data, err := Encrypt(raw, "123")
		xtest.NoError(t, err)
		xtest.Equal(t, enc.Bytes(), data)
	}
}
