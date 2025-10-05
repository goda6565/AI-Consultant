package logger

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
)

type loggerKeyType struct{}

var LoggerKey = loggerKeyType{}

func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

func GetLogger(ctx context.Context) Logger {
	return ctx.Value(LoggerKey).(Logger)
}

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	Panic(msg string, keysAndValues ...interface{})
	LogUsage(llm.Usage)
	Sync() error
}
