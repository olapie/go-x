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
		slog.Info(fmt.Sprintf("received signal %v, will shutdown server", received))
		f()
		if sig, ok := received.(syscall.Signal); ok {
			os.Exit(int(sig))
		} else {
			os.Exit(0)
		}
	}()
}
