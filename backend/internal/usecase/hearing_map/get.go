package hearingmap

import (
	"context"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing_map/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing_map/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/errors"
)

type GetHearingMapInputPort interface {
	Execute(ctx context.Context, input GetHearingMapUseCaseInput) (*GetHearingMapUseCaseOutput, error)
}

type GetHearingMapUseCaseInput struct {
	HearingID string
}

type GetHearingMapUseCaseOutput struct {
	HearingMap *entity.HearingMap
}

type GetHearingMapInteractor struct {
	hearingMapRepository repository.HearingMapRepository
}

func NewGetHearingMapUseCase(hearingMapRepository repository.HearingMapRepository) GetHearingMapInputPort {
	return &GetHearingMapInteractor{hearingMapRepository: hearingMapRepository}
}

func (i *GetHearingMapInteractor) Execute(ctx context.Context, input GetHearingMapUseCaseInput) (*GetHearingMapUseCaseOutput, error) {
	// validate and create hearing ID
	hearingID, err := sharedValue.NewID(input.HearingID)
	if err != nil {
		return nil, fmt.Errorf("failed to create hearing id: %w", err)
	}

	// find hearing map by hearing ID
	hearingMap, err := i.hearingMapRepository.FindByHearingID(ctx, hearingID)
	if err != nil {
		return nil, fmt.Errorf("failed to find hearing map: %w", err)
	}
	if hearingMap == nil {
		return nil, errors.NewUseCaseError(errors.NotFoundError, "hearing map not found")
	}

	return &GetHearingMapUseCaseOutput{HearingMap: hearingMap}, nil
}
