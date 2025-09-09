package zap

import (
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapSugarLogger struct {
	logger *zap.SugaredLogger
}

func (z *zapSugarLogger) Info(msg string, kv ...interface{})  { z.logger.Infow(msg, kv...) }
func (z *zapSugarLogger) Debug(msg string, kv ...interface{}) { z.logger.Debugw(msg, kv...) }
func (z *zapSugarLogger) Warn(msg string, kv ...interface{})  { z.logger.Warnw(msg, kv...) }
func (z *zapSugarLogger) Error(msg string, kv ...interface{}) { z.logger.Errorw(msg, kv...) }
func (z *zapSugarLogger) Fatal(msg string, kv ...interface{}) { z.logger.Fatalw(msg, kv...) }
func (z *zapSugarLogger) Panic(msg string, kv ...interface{}) { z.logger.Panicw(msg, kv...) }
func (z *zapSugarLogger) Sync() error                         { return z.logger.Sync() }
func (z *zapSugarLogger) LogUsage(usage llm.Usage) {
	z.logger.Infow("LLM Usage", "inputTokens", usage.InputTokens, "outputTokens", usage.OutputTokens, "totalTokens", usage.TotalTokens)
}

func ProvideZapLogger(e *environment.Environment) (logger.Logger, func()) {
	var raw *zap.Logger
	if e.Env == "development" {
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		raw = zap.Must(cfg.Build())
	} else {
		cfg := zap.NewProductionConfig()
		raw = zap.Must(cfg.Build())
	}
	sugar := raw.Sugar()
	z := &zapSugarLogger{logger: sugar}
	return z, func() {
		_ = z.Sync()
	}
}
