package di

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/job"
)

type Job struct {
	Job job.Job
}

func (a *Job) Run(ctx context.Context) {
	a.Job.Run(ctx)
}
