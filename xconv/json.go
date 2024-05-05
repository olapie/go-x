package xconv

import (
	"encoding/json"
	"fmt"
)

func MustToJSONBytes(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Errorf("json marshal: %w", err))
	}
	return data
}

func MustFromJSONBytes(b []byte, v any) {
	err := json.Unmarshal(b, v)
	if err != nil {
		panic(fmt.Errorf("unmarshal json: %w", err))
	}
}

func MustToJSONString(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Errorf("json marshal: %w", err))
	}
	return string(b)
}

func MustFromJSONString(s string, v any) {
	err := json.Unmarshal([]byte(s), v)
	if err != nil {
		panic(fmt.Errorf("unmarshal json: %w", err))
	}
}
