package cloudrunjob

import (
	"context"
	"fmt"

	run "cloud.google.com/go/run/apiv2"
	runpb "cloud.google.com/go/run/apiv2/runpb"

	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/ports/job"
)

type CloudRunJobClient struct {
	client *run.JobsClient
	env    *environment.Environment
}

func NewCloudRunJobClient(ctx context.Context, env *environment.Environment) (job.Job, error) {
	client, err := run.NewJobsClient(ctx)
	if err != nil {
		return nil, err
	}
	return &CloudRunJobClient{client: client, env: env}, nil
}

func (c *CloudRunJobClient) CallJob(ctx context.Context, input job.JobInput) error {
	req := &runpb.RunJobRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/jobs/%s",
			c.env.ProjectID,
			c.env.JobRegion,
			input.JobName,
		),
		Overrides: &runpb.RunJobRequest_Overrides{
			ContainerOverrides: []*runpb.RunJobRequest_Overrides_ContainerOverride{
				{
					Env: []*runpb.EnvVar{
						{Name: "PROBLEM_ID", Values: &runpb.EnvVar_Value{Value: input.ProblemID}},
					},
				},
			},
		},
	}

	_, err := c.client.RunJob(ctx, req)
	if err != nil {
		return err
	}
	return err
}
