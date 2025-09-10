package publish

import (
	pubsubClient "cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/ports/pubsub"
)

type Publisher struct {
	client *pubsubClient.Client
	env    *environment.Environment
}

func NewPublisher(client *pubsubClient.Client, e *environment.Environment) pubsub.Publisher {
	return &Publisher{
		client: client,
		env:    e,
	}
}

func (p *Publisher) Publish(ctx context.Context, message pubsub.PubsubMessage) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to marshal message: %v", err))
	}
	topic := p.client.Topic(p.env.TopicName)
	result := topic.Publish(ctx, &pubsubClient.Message{
		Data: payload,
	})
	_, err = result.Get(ctx)
	if err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to publish pubsub message: %v", err))
	}
	return err
}
