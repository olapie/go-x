package xapp

import (
	"flag"
	"fmt"
	"log/slog"

	"go.olapie.com/x/xlog"
)

func Initialize(appName string, httpServerAddr, grpcServerAddr *string) {
	logFilename := flag.String("logfile", "", "log filename")
	flag.Parse()
	CheckVersionArgument(appName)
	closer := xlog.Init(func(options *xlog.Options) {
		options.Filename = *logFilename
	})

	CleanUp(func() {
		if closer != nil {
			_ = closer.Close()
		}
	})

	go func() {
		if httpServerAddr != nil && *httpServerAddr != "" {
			slog.Info(fmt.Sprintf("http server is running on %v", *httpServerAddr))
		}

		if grpcServerAddr != nil && *grpcServerAddr != "" {
			slog.Info(fmt.Sprintf("grpc server is running on %v", *grpcServerAddr))
		}

		if closer != nil {
			fmt.Printf("log file: %s\n", *logFilename)

			if httpServerAddr != nil && *httpServerAddr != "" {
				fmt.Printf("http server is running on %v\n", *httpServerAddr)
			}

			if grpcServerAddr != nil && *grpcServerAddr != "" {
				fmt.Printf("grpc server is running on %v\n", *grpcServerAddr)
			}
		}
	}()
}
