package xcompress

import (
	"bytes"
	"compress/gzip"
	"io"
)

func GzipData(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	w := gzip.NewWriter(buf)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	if err = w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GunzipData(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if err = r.Close(); err != nil {
		return nil, err
	}
	return raw, nil
}
