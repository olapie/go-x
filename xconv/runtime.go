package xconv

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
)

func Dereference[T any](p *T) T {
	if p != nil {
		return *p
	}
	var zero T
	return zero
}

func GobCopy(dst, src any) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(src)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	dec := gob.NewDecoder(&b)
	err = dec.Decode(dst)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	return nil
}

func JSONCopy(dst, src any) error {
	b, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	err = json.Unmarshal(b, dst)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	return nil
}
