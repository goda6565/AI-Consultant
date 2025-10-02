package hearing

import (
	"context"
	"fmt"

	hearingEntity "github.com/goda6565/ai-consultant/backend/internal/domain/hearing/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing/repository"
	hearingService "github.com/goda6565/ai-consultant/backend/internal/domain/hearing/service"
	problemValue "github.com/goda6565/ai-consultant/backend/internal/domain/problem/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/uuid"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/errors"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/ports/transaction"
)

type CreateHearingInputPort interface {
	Create(ctx context.Context, input CreateHearingUseCaseInput, logger logger.Logger) (*CreateHearingOutput, error)
}

type CreateHearingUseCaseInput struct {
	ProblemID string
}

type CreateHearingOutput struct {
	Hearing *hearingEntity.Hearing
}

type CreateHearingInteractor struct {
	hearingRepository       repository.HearingRepository
	duplicateCheckerService *hearingService.DuplicateCheckerService
	adminUnitOfWork         transaction.AdminUnitOfWork
}

func NewCreateHearingUseCase(
	hearingRepository repository.HearingRepository,
	duplicateCheckerService *hearingService.DuplicateCheckerService,
	adminUnitOfWork transaction.AdminUnitOfWork,
) CreateHearingInputPort {
	return &CreateHearingInteractor{
		hearingRepository:       hearingRepository,
		duplicateCheckerService: duplicateCheckerService,
		adminUnitOfWork:         adminUnitOfWork,
	}
}

func (i *CreateHearingInteractor) Create(ctx context.Context, input CreateHearingUseCaseInput, logger logger.Logger) (*CreateHearingOutput, error) {
	// validate and create problem ID
	problemID, err := sharedValue.NewID(input.ProblemID)
	if err != nil {
		return nil, fmt.Errorf("failed to create problem id: %w", err)
	}

	var hearing *hearingEntity.Hearing
	err = i.adminUnitOfWork.WithTx(ctx, func(txCtx context.Context) error {
		// check if hearing already exists for this problem (inside transaction)
		hearingRepo := i.adminUnitOfWork.HearingRepository(txCtx)
		existingHearing, err := hearingRepo.FindByProblemId(txCtx, problemID)
		if err != nil {
			return fmt.Errorf("failed to check duplicate hearing: %w", err)
		}
		if existingHearing != nil {
			return errors.NewUseCaseError(errors.DuplicateError, "hearing already exists for this problem")
		}

		// create value objects
		id, err := sharedValue.NewID(uuid.NewUUID())
		if err != nil {
			return fmt.Errorf("failed to create id: %w", err)
		}

		// create hearing
		hearing = hearingEntity.NewHearing(id, problemID, nil)

		// save hearing using transaction-aware repository
		err = hearingRepo.Create(txCtx, hearing)
		if err != nil {
			return fmt.Errorf("failed to save hearing: %w", err)
		}

		// update problem status using transaction-aware repository
		problemRepo := i.adminUnitOfWork.ProblemRepository(txCtx)
		err = problemRepo.UpdateStatus(txCtx, problemID, problemValue.StatusHearing)
		if err != nil {
			return fmt.Errorf("failed to update problem status: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &CreateHearingOutput{
		Hearing: hearing,
	}, nil
}
