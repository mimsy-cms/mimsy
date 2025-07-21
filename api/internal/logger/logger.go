package logger

import "log/slog"

const (
	LogFormatText = "text"
	LogFormatJSON = "json"
)

// LevelToSlogLevel converts a string log level to slog.Level.
// Supported levels are "debug", "info", "warn", "error".
func LevelToSlogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
