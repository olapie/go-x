package xapp

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func CleanUp(f func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		received := <-c
		slog.Info(received.String())
		f()
		if sig, ok := received.(syscall.Signal); ok {
			os.Exit(int(sig))
		} else {
			os.Exit(0)
		}
	}()
}
