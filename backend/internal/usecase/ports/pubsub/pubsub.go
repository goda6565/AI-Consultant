package pubsub

import "context"

type PubsubMessage struct {
	DocumentID string `json:"documentId"`
}

type Publisher interface {
	Publish(ctx context.Context, message PubsubMessage) error
}
