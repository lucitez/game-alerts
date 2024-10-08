package logger

import (
	"log/slog"
	"os"
)

func Init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}
