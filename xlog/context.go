package xlog

import (
	"context"
	"log/slog"
)

type loggerKeyType int

const keyLogger loggerKeyType = iota

// Deprecated: use NewContext
func NewCtx(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, keyLogger, l)
}

func NewContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, keyLogger, l)
}

// Deprecated: use FromContext
func FromCtx(ctx context.Context) *slog.Logger {
	l, _ := ctx.Value(keyLogger).(*slog.Logger)
	if l == nil {
		l = slog.Default()
	}
	return l
}

func FromContext(ctx context.Context) *slog.Logger {
	l, _ := ctx.Value(keyLogger).(*slog.Logger)
	if l == nil {
		l = slog.Default()
	}
	return l
}
