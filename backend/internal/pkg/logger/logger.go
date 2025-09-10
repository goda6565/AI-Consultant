package logger

import (
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
)

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
