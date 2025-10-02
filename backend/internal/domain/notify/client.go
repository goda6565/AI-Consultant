package notify

import "context"

type NotifyInput struct {
	ProblemID string
	Phase     string
	Message   string
}

type Client interface {
	Notify(ctx context.Context, input NotifyInput) error
}
