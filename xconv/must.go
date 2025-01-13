package xconv

// MustGet eliminates nil err and panics if err isn't nil
func MustGet[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetTwo eliminates nil err and panics if err isn't nil
func MustGetTwo[T1 any, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err)
	}
	return v1, v2
}
