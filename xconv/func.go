package xconv

type FuncE[S any, T any] func(S, error) (T, error)

func ToFuncE[S any, T any](fn func(S) T) FuncE[S, T] {
	return func(s S, err error) (T, error) {
		var t T
		if err != nil {
			return t, err
		}
		return fn(s), nil
	}
}
