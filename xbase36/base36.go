package xbase36

import (
	"errors"
	"math/big"

	"github.com/google/uuid"
)

const base = 36

func EncodeToString(src []byte) string {
	if len(src) == 0 {
		return ""
	}
	var i big.Int
	i.SetBytes(src)
	return i.Text(base)
}

func DecodeString(s string) ([]byte, error) {
	if s == "" {
		return nil, nil
	}
	var i big.Int
	_, ok := i.SetString(s, base)
	if !ok {
		return nil, errors.New("illegal base36 data")
	}
	return i.Bytes(), nil
}

func NewUUIDString() string {
	id := uuid.New()
	return EncodeToString(id[:])
}

func UUIDFromString(s string) (id uuid.UUID, err error) {
	b, err := DecodeString(s)
	if err != nil {
		return id, err
	}
	copy(id[len(id)-len(b):], b[:])
	//return uuid.FromBytes(b)
	return id, nil
}
