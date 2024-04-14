package xconv

func UniqueSlice[T comparable](a []T) []T {
	m := make(map[T]struct{}, len(a))
	l := make([]T, 0, len(a))
	for _, v := range a {
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		l = append(l, v)
	}
	return l
}

func SliceToSet[E comparable](a []E) map[E]bool {
	m := make(map[E]bool)
	for _, v := range a {
		m[v] = true
	}
	return m
}

func TransformToSlice[E1 any, E2 any](a []E1, transformer func(e E1) E2) []E2 {
	res := make([]E2, len(a))
	for i, e := range a {
		res[i] = transformer(e)
	}
	return res
}

func TransformToSliceE[E1 any, E2 any](a []E1, transformer func(e E1) (E2, error)) ([]E2, error) {
	res := make([]E2, len(a))
	for i, e := range a {
		e2, err := transformer(e)
		if err != nil {
			return nil, err
		}
		res[i] = e2
	}
	return res, nil
}
