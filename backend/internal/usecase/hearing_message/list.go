package hearing_message

import (
	"context"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type ListHearingMessageInputPort interface {
	Execute(ctx context.Context, input ListHearingMessageUseCaseInput) (*ListHearingMessageOutput, error)
}

type ListHearingMessageUseCaseInput struct {
	HearingID string
}

type ListHearingMessageOutput struct {
	HearingMessages []entity.HearingMessage
}

type ListHearingMessageInteractor struct {
	hearingMessageRepository repository.HearingMessageRepository
}

func NewListHearingMessageUseCase(hearingMessageRepository repository.HearingMessageRepository) ListHearingMessageInputPort {
	return &ListHearingMessageInteractor{hearingMessageRepository: hearingMessageRepository}
}

func (i *ListHearingMessageInteractor) Execute(ctx context.Context, input ListHearingMessageUseCaseInput) (*ListHearingMessageOutput, error) {
	// validate and create hearing ID
	hearingID, err := sharedValue.NewID(input.HearingID)
	if err != nil {
		return nil, fmt.Errorf("failed to create hearing id: %w", err)
	}

	// find hearing messages by hearing ID
	hearingMessages, err := i.hearingMessageRepository.FindByHearingID(ctx, hearingID)
	if err != nil {
		return nil, fmt.Errorf("failed to find hearing messages: %w", err)
	}

	return &ListHearingMessageOutput{HearingMessages: hearingMessages}, nil
}
