package job

import "context"

type JobInput struct {
	JobName   string
	ProblemID string
}

type Job interface {
	CallJob(ctx context.Context, input JobInput) error
}
