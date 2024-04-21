package xlog

import (
	"io"
	"log"
	"log/slog"
	"os"
	"strings"
)

const (
	EnvDebug = "OLA_LOG_DEBUG" // 1, true, or enabled
)

// Init default logger writing to stderr and filename if it's not empty
func Init(filename string, opts ...func(options *slog.HandlerOptions)) io.Closer {
	options := new(slog.HandlerOptions)
	for _, opt := range opts {
		opt(options)
	}

	if IsDebugging() {
		options.Level = slog.LevelDebug
	}

	if filename == "" {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, options)))
		return nil
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("open file %s: %v\n", filename, err)
	}

	h := MultiHandler(slog.NewJSONHandler(f, options),
		slog.NewTextHandler(os.Stderr, options))
	slog.SetDefault(slog.New(h))
	return f
}

func IsDebugging() bool {
	debug := strings.ToLower(os.Getenv(EnvDebug))
	return debug == "1" || debug == "true" || debug == "enabled"
}
