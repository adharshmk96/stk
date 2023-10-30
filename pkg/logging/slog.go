package logging

import (
	"log/slog"
	"os"
)

func NewSlogLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return logger
}
