package xapp

import (
	"fmt"
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
		slog.Info(fmt.Sprintf("received signal %v", received))
		f()
		slog.Info("shutting down server")
		if sig, ok := received.(syscall.Signal); ok {
			os.Exit(int(sig))
		} else {
			os.Exit(0)
		}
	}()
}
