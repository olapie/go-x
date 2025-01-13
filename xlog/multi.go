package xlog

import (
	"context"
	"errors"
	"log/slog"
)

func MultiHandler(handlers ...slog.Handler) slog.Handler {
	return &multiHandler{
		handlers: handlers,
	}
}

type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, record slog.Record) error {
	errs := make([]error, len(m.handlers))
	for i, h := range m.handlers {
		errs[i] = h.Handle(ctx, record)
	}
	return errors.Join(errs...)
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	res := &multiHandler{
		handlers: make([]slog.Handler, len(m.handlers)),
	}
	for i, h := range m.handlers {
		res.handlers[i] = h.WithAttrs(attrs)
	}
	return res
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	res := &multiHandler{
		handlers: make([]slog.Handler, len(m.handlers)),
	}
	for i, h := range m.handlers {
		res.handlers[i] = h.WithGroup(name)
	}
	return res
}
