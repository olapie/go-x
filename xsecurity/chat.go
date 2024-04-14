package xsecurity

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"

	"go.olapie.com/x/xbase62"
)

func EncryptChatMessage(message string, key [32]byte) (string, error) {
	if message == "" {
		return "", nil
	}

	var output bytes.Buffer
	w := gzip.NewWriter(&output)
	_, err := w.Write([]byte(message))
	if err != nil {
		return "", fmt.Errorf("gzip.Write: %w", err)
	}
	err = w.Close()
	if err != nil {
		return "", fmt.Errorf("gzip.Close: %w", err)
	}

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", fmt.Errorf("aes.NewCipher: %w", err)
	}
	stream := cipher.NewCTR(block, key[16:])
	data := output.Bytes()
	stream.XORKeyStream(data, data)

	return xbase62.EncodeToString(data), nil
}

func DecryptChatMessage(message string, key [32]byte) (string, error) {
	if message == "" {
		return "", nil
	}

	data, err := xbase62.DecodeString(message)
	if err != nil {
		return "", fmt.Errorf("xbase62.DecodeString: %w", err)
	}

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", fmt.Errorf("aes.NewCipher: %w", err)
	}
	stream := cipher.NewCTR(block, key[16:])
	stream.XORKeyStream(data, data)

	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("gzip.NewReader: %w", err)
	}
	data, err = io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("io.ReadAll: %w", err)
	}
	err = r.Close()
	if err != nil {
		return "", fmt.Errorf("gzip.Close: %w", err)
	}

	return string(data), nil
}
