package jobconfig

import (
	"context"
	"fmt"

	jobConfigEntity "github.com/goda6565/ai-consultant/backend/internal/domain/job_config/entity"
	jobConfigRepository "github.com/goda6565/ai-consultant/backend/internal/domain/job_config/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/errors"
)

type GetJobConfigInputPort interface {
	Execute(ctx context.Context, input GetJobConfigUseCaseInput) (*GetJobConfigOutput, error)
}

type GetJobConfigUseCaseInput struct {
	ProblemID string
}

type GetJobConfigOutput struct {
	JobConfig *jobConfigEntity.JobConfig
}

type GetJobConfigInteractor struct {
	jobConfigRepository jobConfigRepository.JobConfigRepository
}

func NewGetJobConfigUseCase(jobConfigRepository jobConfigRepository.JobConfigRepository) GetJobConfigInputPort {
	return &GetJobConfigInteractor{jobConfigRepository: jobConfigRepository}
}

func (i *GetJobConfigInteractor) Execute(ctx context.Context, input GetJobConfigUseCaseInput) (*GetJobConfigOutput, error) {
	// validate and create problem ID
	problemID, err := sharedValue.NewID(input.ProblemID)
	if err != nil {
		return nil, fmt.Errorf("failed to create problem id: %w", err)
	}

	// find job config by problem ID
	jobConfig, err := i.jobConfigRepository.FindByProblemID(ctx, problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to find job config: %w", err)
	}
	if jobConfig == nil {
		return nil, errors.NewUseCaseError(errors.NotFoundError, "job config not found")
	}

	return &GetJobConfigOutput{JobConfig: jobConfig}, nil
}
