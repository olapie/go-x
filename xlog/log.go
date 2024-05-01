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

type Options struct {
	slog.HandlerOptions
	Filename      string
	ConsoleOutput bool
}

// Init default logger
func Init(opts ...func(options *Options)) io.Closer {
	options := new(Options)
	for _, opt := range opts {
		opt(options)
	}

	if IsDebugging() {
		options.Level = slog.LevelDebug
	}

	if options.Filename == "" { // regardless of ConsoleOutput
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &options.HandlerOptions)))
		return nil
	}

	f, err := os.OpenFile(options.Filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("open file %s: %v\n", options.Filename, err)
	}

	var h slog.Handler
	if options.ConsoleOutput {
		h = MultiHandler(slog.NewJSONHandler(f, &options.HandlerOptions),
			slog.NewTextHandler(os.Stderr, &options.HandlerOptions))
	} else {
		h = slog.NewJSONHandler(f, &options.HandlerOptions)
	}
	slog.SetDefault(slog.New(h))
	return f
}

func IsDebugging() bool {
	debug := strings.ToLower(os.Getenv(EnvDebug))
	return debug == "1" || debug == "true" || debug == "enabled"
}
