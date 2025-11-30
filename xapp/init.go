package xapp

import (
	"fmt"
	"log/slog"

	"go.olapie.com/x/xlog"
)

type Config struct {
	AppName        string
	LogFilename    string
	HTTPServerAddr string
	GRPCServerAddr string
}

func Initialize(cfg *Config) {
	CheckVersionArgument(cfg.AppName)
	closer := xlog.Init(func(options *xlog.Options) {
		options.Filename = cfg.LogFilename
	})

	CleanUp(func() {
		if closer != nil {
			_ = closer.Close()
		}
	})

	go func() {
		if cfg.HTTPServerAddr != "" {
			slog.Info(fmt.Sprintf("http server is running on %v", cfg.HTTPServerAddr))
		}

		if cfg.GRPCServerAddr != "" {
			slog.Info(fmt.Sprintf("grpc server is running on %v", cfg.GRPCServerAddr))
		}

		if closer != nil {
			fmt.Printf("log file: %s\n", cfg.LogFilename)

			if cfg.HTTPServerAddr != "" {
				fmt.Printf("http server is running on %v\n", cfg.HTTPServerAddr)
			}

			if cfg.GRPCServerAddr != "" {
				fmt.Printf("grpc server is running on %v\n", cfg.GRPCServerAddr)
			}
		}
	}()
}
