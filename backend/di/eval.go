package di

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/evaluate"
)

type Eval struct {
	Evaluator evaluate.BaseEvaluator
}

func (a *Eval) Evaluate(ctx context.Context) {
	a.Evaluator.Evaluate(ctx)
}
