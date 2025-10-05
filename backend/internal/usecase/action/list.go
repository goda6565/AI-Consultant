package action

import (
	"context"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/action/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/action/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type ListActionInputPort interface {
	Execute(ctx context.Context, input ListActionUseCaseInput) (*ListActionOutput, error)
}

type ListActionUseCaseInput struct {
	ProblemID string
}

type ListActionOutput struct {
	Actions []entity.Action
}

type ListActionInteractor struct {
	actionRepository repository.ActionRepository
}

func NewListActionUseCase(actionRepository repository.ActionRepository) ListActionInputPort {
	return &ListActionInteractor{actionRepository: actionRepository}
}

func (i *ListActionInteractor) Execute(ctx context.Context, input ListActionUseCaseInput) (*ListActionOutput, error) {
	// validate and create problem ID
	problemID, err := sharedValue.NewID(input.ProblemID)
	if err != nil {
		return nil, fmt.Errorf("failed to create problem id: %w", err)
	}

	// find actions by problem ID
	actions, err := i.actionRepository.FindByProblemID(ctx, problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to find actions: %w", err)
	}

	return &ListActionOutput{Actions: actions}, nil
}
