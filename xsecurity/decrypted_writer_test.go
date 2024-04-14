package xsecurity

import (
	"bytes"
	"io"
	"testing"
	"time"

	"go.olapie.com/security/internal/testutil"
)

func TestDecryptedWriter(t *testing.T) {
	raw := []byte(SHA1(time.Now().String()))
	enc, err := Encrypt(raw, "123")
	testutil.NoError(t, err)
	dec := &bytes.Buffer{}
	w := NewDecryptedWriter(dec, "123")
	n, err := io.Copy(w, bytes.NewReader(enc))
	t.Log(n)
	testutil.NoError(t, err)
	testutil.Equal(t, raw, dec.Bytes())
}
