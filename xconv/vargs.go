package xconv

import (
	"errors"
	"fmt"
)

func FromVarargs(keyValues ...any) (keys []string, values []any, err error) {
	n := len(keyValues)
	if n%2 != 0 {
		err = errors.New("keyValues should be pairs of (string, any)")
		return
	}

	keys, values = make([]string, 0, n/2), make([]any, 0, n/2)
	for i := 0; i < n/2; i++ {
		if k, ok := keyValues[2*i].(string); !ok {
			err = fmt.Errorf("keyValues[%d] isn't convertible to string", i)
			return
		} else if keyValues[2*i+1] == nil {
			err = fmt.Errorf("keyValues[%d] is nil", 2*i+1)
			return
		} else {
			keys = append(keys, k)
			values = append(values, keyValues[2*i+1])
		}
	}
	return
}
