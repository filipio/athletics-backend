package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB
var dbOnce sync.Once

type SlogLogger struct {
	logLevel logger.LogLevel
}

func (l *SlogLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

func (l *SlogLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Info {
		slog.Info(msg, "data", data)
	}
}

func (l *SlogLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Warn {
		slog.Warn(msg, "data", data)
	}
}

func (l *SlogLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Error {
		slog.Error(msg, "data", data)
	}
}

// Trace logs SQL statements with duration
func (l *SlogLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin).Milliseconds()
	sql, rows := fc()

	switch {
	case err != nil && l.logLevel >= logger.Error:
		slog.Error("Trace", "sql", sql, "rows", rows, "elapsed(ms)", elapsed, "error", err)
	case elapsed > 200 && l.logLevel >= logger.Warn: // Slow query threshold
		slog.Warn("Trace (slow query)", "sql", sql, "rows", rows, "elapsed(ms)", elapsed)
	case l.logLevel >= logger.Info:
		slog.Info("Trace", "sql", sql, "rows", rows, "elapsed(ms)", elapsed)
	}
}

func DatabaseConnection() *gorm.DB {

	dbOnce.Do(func() {
		dbInstance, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn(),
			PreferSimpleProtocol: false,
		}), &gorm.Config{
			Logger: &SlogLogger{logLevel: logger.Info},
		})
		if err != nil {
			panic("failed to connect database")
		}
		db = dbInstance
	})

	return db
}

func dsn() string {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	result := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=allow TimeZone=UTC",
		host, user, password, dbname, port)

	return result
}
