package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
)

func ProvidePubsubClient(ctx context.Context, e *environment.Environment) *pubsub.Client {
	client, err := pubsub.NewClient(ctx, e.ProjectID)
	if err != nil {
		panic(err)
	}
	return client
}
