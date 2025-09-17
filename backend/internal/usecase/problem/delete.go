package problem

import (
	"context"
	"fmt"

	hearingRepository "github.com/goda6565/ai-consultant/backend/internal/domain/hearing/repository"
	hearingMessageRepository "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/repository"
	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/errors"
	transaction "github.com/goda6565/ai-consultant/backend/internal/usecase/ports/transaction"
)

type DeleteProblemInputPort interface {
	Execute(ctx context.Context, input DeleteProblemUseCaseInput) error
}

type DeleteProblemUseCaseInput struct {
	ProblemID string
}

type DeleteProblemInteractor struct {
	problemRepository        repository.ProblemRepository
	hearingMessageRepository hearingMessageRepository.HearingMessageRepository
	hearingRepository        hearingRepository.HearingRepository
	adminUnitOfWork          transaction.AdminUnitOfWork
}

func NewDeleteProblemUseCase(problemRepository repository.ProblemRepository, hearingMessageRepository hearingMessageRepository.HearingMessageRepository, hearingRepository hearingRepository.HearingRepository, adminUnitOfWork transaction.AdminUnitOfWork) DeleteProblemInputPort {
	return &DeleteProblemInteractor{problemRepository: problemRepository, hearingMessageRepository: hearingMessageRepository, hearingRepository: hearingRepository, adminUnitOfWork: adminUnitOfWork}
}

func (i *DeleteProblemInteractor) Execute(ctx context.Context, input DeleteProblemUseCaseInput) error {
	// validate and create problem ID
	problemID, err := sharedValue.NewID(input.ProblemID)
	if err != nil {
		return fmt.Errorf("failed to create problem id: %w", err)
	}

	// check if problem exists
	problem, err := i.problemRepository.FindById(ctx, problemID)
	if err != nil {
		return fmt.Errorf("failed to find problem: %w", err)
	}
	if problem == nil {
		return errors.NewUseCaseError(errors.NotFoundError, "problem not found")
	}

	// check if hearing exists
	hearing, err := i.hearingRepository.FindByProblemId(ctx, problemID)
	if err != nil {
		return fmt.Errorf("failed to find hearing: %w", err)
	}

	// delete all
	err = i.adminUnitOfWork.WithTx(ctx, func(ctx context.Context) error {
		if hearing != nil {
			// delete hearing messages
			_, err = i.hearingMessageRepository.DeleteByHearingID(ctx, hearing.GetID())
			if err != nil {
				return fmt.Errorf("failed to delete hearing messages: %w", err)
			}

			// delete hearing
			_, err = i.hearingRepository.DeleteByProblemID(ctx, problem.GetID())
			if err != nil {
				return fmt.Errorf("failed to delete hearing: %w", err)
			}
		}

		// delete problem
		numDeleted, err := i.problemRepository.Delete(ctx, problem.GetID())
		if err != nil {
			return fmt.Errorf("failed to delete problem: %w", err)
		}
		if numDeleted == 0 {
			return errors.NewUseCaseError(errors.NotFoundError, "problem not found")
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to delete problem with admin unit of work: %w", err)
	}

	return nil
}
