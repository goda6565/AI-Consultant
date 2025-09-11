package queue

import "context"

type SyncQueueMessage struct {
	DocumentID string `json:"documentId"`
}

type SyncQueue interface {
	Enqueue(ctx context.Context, message SyncQueueMessage) error
}
