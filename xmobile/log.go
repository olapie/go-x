package xmobile

import "log/slog"

var logger = slog.Default()

func SetLogger(l *slog.Logger) {
	logger = l
}

func LogDebug(s string) {
	logger.Debug(s)
}

func LogInfo(s string) {
	logger.Info(s)
}

func LogWarn(s string) {
	logger.Warn(s)
}

func LogError(s string) {
	logger.Error(s)
}
