package xbase62

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/google/uuid"
)

const base = 62

func EncodeToString(src []byte) string {
	if len(src) == 0 {
		return ""
	}
	var i big.Int
	i.SetBytes(src) // leading zeros will be ignored, so DecodeString won't output same bytes
	return i.Text(base)
}

func DecodeString(s string) ([]byte, error) {
	if s == "" {
		return nil, nil
	}
	var i big.Int
	_, ok := i.SetString(s, base)
	if !ok {
		return nil, errors.New("illegal format")
	}
	return i.Bytes(), nil
}

func Itoa(i int64) string {
	return big.NewInt(i).Text(base)
}

func Atoi(s string) (int64, error) {
	if s == "" {
		return 0, errors.New("error parsing empty string into int64")
	}
	var i big.Int
	_, ok := i.SetString(s, base)
	if !ok {
		return 0, errors.New("illegal format")
	}
	if i.IsInt64() {
		return i.Int64(), nil
	}
	return 0, fmt.Errorf("error parsing %s into int64", s)
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
