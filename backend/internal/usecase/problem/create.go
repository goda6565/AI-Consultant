package problem

import (
	"context"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/repository"
	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/service"
	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/uuid"
)

type CreateProblemInputPort interface {
	Execute(ctx context.Context, input CreateProblemUseCaseInput) (*CreateProblemOutput, error)
}

type CreateProblemUseCaseInput struct {
	Description string
}

type CreateProblemOutput struct {
	Problem *entity.Problem
}

type CreateProblemInteractor struct {
	generateTitleService *service.GenerateTitleService
	problemRepository    repository.ProblemRepository
}

func NewCreateProblemUseCase(generateTitleService *service.GenerateTitleService, problemRepository repository.ProblemRepository) CreateProblemInputPort {
	return &CreateProblemInteractor{generateTitleService: generateTitleService, problemRepository: problemRepository}
}

func (i *CreateProblemInteractor) Execute(ctx context.Context, input CreateProblemUseCaseInput) (*CreateProblemOutput, error) {
	description, err := value.NewDescription(input.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create description: %w", err)
	}

	// generate title
	title, err := i.generateTitleService.Execute(ctx, service.GenerateTitleServiceInput{Description: input.Description})
	if err != nil {
		return nil, fmt.Errorf("failed to generate title: %w", err)
	}

	// create value objects
	id, err := sharedValue.NewID(uuid.NewUUID())
	if err != nil {
		return nil, fmt.Errorf("failed to create id: %w", err)
	}
	titleValue, err := value.NewTitle(title.Title)
	if err != nil {
		return nil, fmt.Errorf("failed to create title: %w", err)
	}

	// create problem
	problem := entity.NewProblem(id, *titleValue, *description, value.StatusPending, nil)

	// save problem
	err = i.problemRepository.Create(ctx, problem)
	if err != nil {
		return nil, fmt.Errorf("failed to save problem: %w", err)
	}

	return &CreateProblemOutput{Problem: problem}, nil
}
