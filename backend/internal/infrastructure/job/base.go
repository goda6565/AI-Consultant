package job

import (
	"context"
	"os"

	"github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
)

type JobApplication interface {
	Execute(ctx context.Context) error
}

type Job interface {
	Run(ctx context.Context)
}

type BaseJob struct {
	logger      logger.Logger
	application JobApplication
}

func NewBaseJob(ctx context.Context, logger logger.Logger, application JobApplication) Job {
	return &BaseJob{logger: logger, application: application}
}

func (b *BaseJob) Run(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case error:
				b.logger.Error("error", "error", r)
			default:
				b.logger.Error("panic", "panic", r)
			}
			os.Exit(1)
		}
	}()
	// set logger to context
	ctx = logger.WithLogger(ctx, b.logger)
	err := b.application.Execute(ctx)
	if err != nil {
		b.logger.Error("error", "error", err)
	}
}
