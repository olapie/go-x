package xlog

import "log/slog"

func Int64[T ~int64 | ~int | ~int32 | ~int16 | ~int8](key string, value T) slog.Attr {
	return slog.Int64(key, int64(value))
}

func Int[T ~int | ~int32 | ~int16 | ~int8](key string, value T) slog.Attr {
	return slog.Int(key, int(value))
}

func Float64[T ~float64 | ~float32](key string, value T) slog.Attr {
	return slog.Float64(key, float64(value))
}

func String[T ~string](key string, value T) slog.Attr {
	return slog.String(key, string(value))
}

func Bool[T ~bool](key string, value T) slog.Attr {
	return slog.Bool(key, bool(value))
}

func Err(err error) slog.Attr {
	return slog.String("err", err.Error())
}
