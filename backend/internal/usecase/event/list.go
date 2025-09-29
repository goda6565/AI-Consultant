package event

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/event/entity"
	repository "github.com/goda6565/ai-consultant/backend/internal/domain/event/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type ListEventInputPort interface {
	Execute(ctx context.Context, input ListEventUseCaseInput) (*ListEventOutput, error)
}

type ListEventUseCaseInput struct {
	ProblemID string
}

type ListEventOutput struct {
	Events []entity.Event
}

type ListEventInteractor struct {
	repository repository.EventRepository
}

func NewListEventUseCase(repository repository.EventRepository) ListEventInputPort {
	return &ListEventInteractor{repository: repository}
}

func (i *ListEventInteractor) Execute(ctx context.Context, input ListEventUseCaseInput) (*ListEventOutput, error) {
	problemID, err := sharedValue.NewID(input.ProblemID)
	if err != nil {
		return nil, err
	}
	events, err := i.repository.FindAllByProblemID(ctx, problemID)
	if err != nil {
		return nil, err
	}
	return &ListEventOutput{Events: events}, nil
}
