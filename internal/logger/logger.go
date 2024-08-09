package logger

import (
	"log/slog"
	"os"
)

func Init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: replacer}))
	slog.SetDefault(logger)
}

// replacer renames attribute keys to match Cloud Logging structured log format
// https://cloud.google.com/logging/docs/structured-logging
// https://cloud.google.com/stackdriver/docs/instrumentation/setup/go#config-structured-logging
func replacer(groups []string, a slog.Attr) slog.Attr {

	switch a.Key {
	case slog.LevelKey:
		a.Key = "severity"
		// Map slog.Level string values to Cloud Logging LogSeverity
		// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogSeverity
		if level := a.Value.Any().(slog.Level); level == slog.LevelWarn {
			a.Value = slog.StringValue("WARNING")
		}
	case slog.TimeKey:
		a.Key = "timestamp"
	case slog.MessageKey:
		a.Key = "message"
	}
	return a
}
