package evaluate

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
)

type Evaluator interface {
	Execute(ctx context.Context) error
}

type BaseEvaluator interface {
	Evaluate(ctx context.Context)
}

type baseEvaluator struct {
	logger    logger.Logger
	evaluator Evaluator
}

func NewBaseEvaluator(logger logger.Logger, evaluator Evaluator) BaseEvaluator {
	return &baseEvaluator{logger: logger, evaluator: evaluator}
}

func (e *baseEvaluator) Evaluate(ctx context.Context) {
	ctx = logger.WithLogger(ctx, e.logger)
	defer func() {
		if r := recover(); r != nil {
			e.logger.Error("panic", "panic", r)
		}
	}()
	err := e.evaluator.Execute(ctx)
	if err != nil {
		e.logger.Error("error", "error", err)
	}
}
