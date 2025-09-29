package event

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/event/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/event/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type StreamEventInputPort interface {
	Execute(ctx context.Context, input StreamEventUseCaseInput) (*StreamEventOutput, error)
}

type StreamEventUseCaseInput struct {
	ProblemID string
}

type StreamEventOutput struct {
	Stream <-chan entity.Event
}

type StreamEventInteractor struct {
	stream repository.EventRepository
}

func NewStreamEventUseCase(stream repository.EventRepository) StreamEventInputPort {
	return &StreamEventInteractor{stream: stream}
}

func (i *StreamEventInteractor) Execute(ctx context.Context, input StreamEventUseCaseInput) (*StreamEventOutput, error) {
	problemID, err := sharedValue.NewID(input.ProblemID)
	if err != nil {
		return nil, err
	}
	stream, err := i.stream.FindAllByProblemIDAsStream(ctx, problemID)
	if err != nil {
		return nil, err
	}
	return &StreamEventOutput{Stream: stream}, nil
}
