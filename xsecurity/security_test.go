package xsecurity

import (
	"crypto/rand"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go.olapie.com/x/xtest"
)

func TestEncrypt(t *testing.T) {
	password := SHA1(time.Now().String())
	testEncrypt(t, 1<<4+9, password)
	testEncrypt(t, 1<<24, password)
}

func testEncrypt(t *testing.T, size int, password string) {
	raw := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, raw[:])
	if err != nil {
		t.Fatal(err)
	}

	enc, err := Encrypt(raw[:], password)
	t.Log(enc[:30])
	xtest.NoError(t, err)
	xtest.True(t, IsEncrypted(enc))
	dec, err := Decrypt(enc, password)
	xtest.NoError(t, err)
	xtest.False(t, IsEncrypted(dec), dec[:HeaderSize])
	xtest.Equal(t, raw, dec)
}

func TestEncryptFile(t *testing.T) {
	err := os.MkdirAll("testdata", 0755)
	if err != nil {
		t.Fatal(err)
	}

	rawFilename := "testdata/rawfile"
	largeFilename := "testdata/largefile"
	t.Cleanup(func() {
		os.RemoveAll(rawFilename)
		os.RemoveAll(largeFilename)
	})

	password := SHA1(time.Now().String())
	var raw [32]byte
	n, err := io.ReadFull(rand.Reader, raw[:])
	xtest.NoError(t, err)
	t.Log(n, raw)
	f, err := os.OpenFile(rawFilename, os.O_CREATE|os.O_WRONLY, 0644)
	xtest.NoError(t, err)

	_, err = f.Write(raw[:])
	f.Close()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rawFilename)
	testEncryptFile(t, rawFilename, password)

	var large [32 * 1024 * 1024]byte
	n, err = io.ReadFull(rand.Reader, large[:])
	xtest.NoError(t, err)

	f, err = os.OpenFile(largeFilename, os.O_CREATE|os.O_WRONLY, 0644)
	xtest.NoError(t, err)

	_, err = f.Write(large[:])
	f.Close()
	xtest.NoError(t, err)

	testEncryptFile(t, largeFilename, password)
}

func testEncryptFile(t *testing.T, rawFilename string, password string) {
	encFilename := rawFilename + ".enc" + filepath.Ext(rawFilename)
	decFilename := rawFilename + ".dec" + filepath.Ext(rawFilename)
	t.Cleanup(func() {
		os.RemoveAll(decFilename)
		os.RemoveAll(encFilename)
	})
	err := EncryptFile(SF(rawFilename), DF(encFilename), password)
	xtest.NoError(t, err)

	xtest.True(t, IsEncryptedFile(encFilename))
	xtest.False(t, IsEncryptedFile(rawFilename))
	err = DecryptFile(SF(encFilename), DF(decFilename), password)
	xtest.NoError(t, err)
	raw, err := os.ReadFile(rawFilename)
	xtest.NoError(t, err)

	enc, err := os.ReadFile(encFilename)
	xtest.NoError(t, err)
	xtest.NotEqual(t, raw, enc)

	dec, err := os.ReadFile(decFilename)
	xtest.NoError(t, err)
	xtest.Equal(t, raw, dec)

	valid := ValidateFilePassword(encFilename, password)
	xtest.True(t, valid)
}
