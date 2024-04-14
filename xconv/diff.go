package xconv

import (
	"fmt"
)

func diffSlice[E comparable](expected, got []E) string {
	if len(expected) != len(got) {
		return fmt.Sprintf("expect len %d, got len %d", len(expected), len(got))
	}

	for i, e := range expected {
		if e != got[i] {
			return fmt.Sprintf("index %d: expect %v, got %v", i, e, got[i])
		}
	}

	return ""
}
