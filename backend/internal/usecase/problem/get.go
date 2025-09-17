package problem

import (
	"context"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/errors"
)

type GetProblemInputPort interface {
	Execute(ctx context.Context, input GetProblemUseCaseInput) (*GetProblemOutput, error)
}

type GetProblemUseCaseInput struct {
	ProblemID string
}

type GetProblemOutput struct {
	Problem *entity.Problem
}

type GetProblemInteractor struct {
	problemRepository repository.ProblemRepository
}

func NewGetProblemUseCase(problemRepository repository.ProblemRepository) GetProblemInputPort {
	return &GetProblemInteractor{problemRepository: problemRepository}
}

func (i *GetProblemInteractor) Execute(ctx context.Context, input GetProblemUseCaseInput) (*GetProblemOutput, error) {
	// validate and create problem ID
	problemID, err := sharedValue.NewID(input.ProblemID)
	if err != nil {
		return nil, fmt.Errorf("failed to create problem id: %w", err)
	}

	// find problem
	problem, err := i.problemRepository.FindById(ctx, problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to find problem: %w", err)
	}
	if problem == nil {
		return nil, errors.NewUseCaseError(errors.NotFoundError, "problem not found")
	}

	return &GetProblemOutput{Problem: problem}, nil
}
