package hearing

import (
	"context"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing/repository"
	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing/service"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/uuid"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/errors"
)

type CreateHearingInputPort interface {
	Execute(ctx context.Context, input CreateHearingUseCaseInput) (*CreateHearingOutput, error)
}

type CreateHearingUseCaseInput struct {
	ProblemID string
}

type CreateHearingOutput struct {
	Hearing *entity.Hearing
}

type CreateHearingInteractor struct {
	hearingRepository       repository.HearingRepository
	duplicateCheckerService *service.DuplicateCheckerService
}

func NewCreateHearingUseCase(hearingRepository repository.HearingRepository, duplicateCheckerService *service.DuplicateCheckerService) CreateHearingInputPort {
	return &CreateHearingInteractor{hearingRepository: hearingRepository, duplicateCheckerService: duplicateCheckerService}
}

func (i *CreateHearingInteractor) Execute(ctx context.Context, input CreateHearingUseCaseInput) (*CreateHearingOutput, error) {
	// validate and create problem ID
	problemID, err := sharedValue.NewID(input.ProblemID)
	if err != nil {
		return nil, fmt.Errorf("failed to create problem id: %w", err)
	}

	// check if hearing already exists for this problem
	isDuplicate, err := i.duplicateCheckerService.Execute(ctx, problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to check duplicate hearing: %w", err)
	}
	if isDuplicate {
		return nil, errors.NewUseCaseError(errors.DuplicateError, "hearing already exists for this problem")
	}

	// create value objects
	id, err := sharedValue.NewID(uuid.NewUUID())
	if err != nil {
		return nil, fmt.Errorf("failed to create id: %w", err)
	}

	// create hearing
	hearing := entity.NewHearing(id, problemID, nil)

	// save hearing
	err = i.hearingRepository.Create(ctx, hearing)
	if err != nil {
		return nil, fmt.Errorf("failed to save hearing: %w", err)
	}

	return &CreateHearingOutput{Hearing: hearing}, nil
}
