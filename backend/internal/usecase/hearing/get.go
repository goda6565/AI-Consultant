package hearing

import (
	"context"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/errors"
)

type GetHearingInputPort interface {
	Execute(ctx context.Context, input GetHearingUseCaseInput) (*GetHearingOutput, error)
}

type GetHearingUseCaseInput struct {
	ProblemID string
}

type GetHearingOutput struct {
	Hearing *entity.Hearing
}

type GetHearingInteractor struct {
	hearingRepository repository.HearingRepository
}

func NewGetHearingUseCase(hearingRepository repository.HearingRepository) GetHearingInputPort {
	return &GetHearingInteractor{hearingRepository: hearingRepository}
}

func (i *GetHearingInteractor) Execute(ctx context.Context, input GetHearingUseCaseInput) (*GetHearingOutput, error) {
	// validate and create problem ID
	problemID, err := sharedValue.NewID(input.ProblemID)
	if err != nil {
		return nil, fmt.Errorf("failed to create problem id: %w", err)
	}

	// find hearing by problem ID
	hearing, err := i.hearingRepository.FindByProblemId(ctx, problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to find hearing: %w", err)
	}
	if hearing == nil {
		return nil, errors.NewUseCaseError(errors.NotFoundError, "hearing not found")
	}

	return &GetHearingOutput{Hearing: hearing}, nil
}
