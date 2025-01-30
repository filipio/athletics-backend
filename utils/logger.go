package utils

import (
	"log/slog"
	"os"
)

func SetupLogger() {
	var slogHandler slog.Handler

	if os.Getenv("APP_ENV") != "prod" {
		slogHandler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		logPath := os.Getenv("LOG_PATH")
		if logPath == "" {
			slog.Error("Failed to determine log file path")
			os.Exit(1)
		}

		logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		defer logFile.Close()

		slogHandler = slog.NewJSONHandler(logFile, nil)
	}

	logger := slog.New(slogHandler)
	slog.SetDefault(logger)
}
