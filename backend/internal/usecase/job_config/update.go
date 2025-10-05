package jobconfig

import (
	"context"
	"fmt"

	jobConfigEntity "github.com/goda6565/ai-consultant/backend/internal/domain/job_config/entity"
	jobConfigRepository "github.com/goda6565/ai-consultant/backend/internal/domain/job_config/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/errors"
)

type UpdateJobConfigInputPort interface {
	Execute(ctx context.Context, input UpdateJobConfigUseCaseInput) (*UpdateJobConfigOutput, error)
}

type UpdateJobConfigUseCaseInput struct {
	ProblemID            string
	EnableInternalSearch bool
}

type UpdateJobConfigOutput struct {
	JobConfig *jobConfigEntity.JobConfig
}

type UpdateJobConfigInteractor struct {
	jobConfigRepository jobConfigRepository.JobConfigRepository
}

func NewUpdateJobConfigUseCase(
	jobConfigRepository jobConfigRepository.JobConfigRepository,
) UpdateJobConfigInputPort {
	return &UpdateJobConfigInteractor{
		jobConfigRepository: jobConfigRepository,
	}
}

func (u *UpdateJobConfigInteractor) Execute(ctx context.Context, input UpdateJobConfigUseCaseInput) (*UpdateJobConfigOutput, error) {
	problemID, err := sharedValue.NewID(input.ProblemID)
	if err != nil {
		return nil, fmt.Errorf("invalid problem id: %w", err)
	}

	existingJobConfig, err := u.jobConfigRepository.FindByProblemID(ctx, problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to find job config: %w", err)
	}

	if existingJobConfig == nil {
		return nil, errors.NewUseCaseError(errors.NotFoundError, "job config not found")
	}

	if input.EnableInternalSearch {
		existingJobConfig.EnableInternalSearch()
	} else {
		existingJobConfig.DisableInternalSearch()
	}

	err = u.jobConfigRepository.Update(ctx, existingJobConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to update job config: %w", err)
	}

	return &UpdateJobConfigOutput{
		JobConfig: existingJobConfig,
	}, nil
}
