package cloudtasks

import (
	"context"
	"encoding/json"
	"fmt"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"

	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	queuePorts "github.com/goda6565/ai-consultant/backend/internal/usecase/ports/queue"
)

type CloudTasksClient struct {
	client *cloudtasks.Client
	env    *environment.Environment
}

func NewCloudTasksClient(ctx context.Context, e *environment.Environment) (queuePorts.SyncQueue, func()) {
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	return &CloudTasksClient{client: client, env: e}, func() {
		err = client.Close()
		if err != nil {
			panic(err)
		}
	}
}

func (c *CloudTasksClient) Enqueue(ctx context.Context, message queuePorts.SyncQueueMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to marshal message: %v", err))
	}

	request := &taskspb.CreateTaskRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s/queues/%s", c.env.ProjectID, c.env.QueueLocation, c.env.QueueName),
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        c.env.TargetURL,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: body,
				},
			},
		},
	}

	_, err = c.client.CreateTask(ctx, request)
	if err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to enqueue task: %v", err))
	}
	return nil
}
