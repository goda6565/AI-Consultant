package problem

import (
	"context"
	"fmt"

	actionRepository "github.com/goda6565/ai-consultant/backend/internal/domain/action/repository"
	hearingRepository "github.com/goda6565/ai-consultant/backend/internal/domain/hearing/repository"
	hearingMessageRepository "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/repository"
	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/repository"
	reportRepository "github.com/goda6565/ai-consultant/backend/internal/domain/report/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
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
	actionRepository         actionRepository.ActionRepository
	adminUnitOfWork          transaction.AdminUnitOfWork
	reportRepository         reportRepository.ReportRepository
}

func NewDeleteProblemUseCase(
	problemRepository repository.ProblemRepository,
	hearingMessageRepository hearingMessageRepository.HearingMessageRepository,
	hearingRepository hearingRepository.HearingRepository,
	actionRepository actionRepository.ActionRepository,
	adminUnitOfWork transaction.AdminUnitOfWork,
	reportRepository reportRepository.ReportRepository,
) DeleteProblemInputPort {
	return &DeleteProblemInteractor{
		problemRepository:        problemRepository,
		hearingMessageRepository: hearingMessageRepository,
		hearingRepository:        hearingRepository,
		actionRepository:         actionRepository,
		adminUnitOfWork:          adminUnitOfWork,
		reportRepository:         reportRepository,
	}
}

func (i *DeleteProblemInteractor) Execute(ctx context.Context, input DeleteProblemUseCaseInput) error {
	logger := logger.GetLogger(ctx)
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

	// check if job config exists
	jobConfig, err := i.adminUnitOfWork.JobConfigRepository(ctx).FindByProblemID(ctx, problemID)
	if err != nil {
		return fmt.Errorf("failed to find job config: %w", err)
	}

	// check if hearing exists
	hearings, err := i.hearingRepository.FindAllByProblemId(ctx, problemID)
	if err != nil {
		return fmt.Errorf("failed to find hearing: %w", err)
	}

	// check if actions exists
	actions, err := i.actionRepository.FindByProblemID(ctx, problemID)
	if err != nil {
		return fmt.Errorf("failed to find actions: %w", err)
	}

	// check if report exists
	report, err := i.reportRepository.FindByProblemID(ctx, problemID)
	if err != nil {
		return fmt.Errorf("failed to find report: %w", err)
	}

	// delete all
	err = i.adminUnitOfWork.WithTx(ctx, func(ctx context.Context) error {
		// delete reports if exists
		if report != nil {
			_, err = i.adminUnitOfWork.ReportRepository(ctx).DeleteByProblemID(ctx, problemID)
			if err != nil {
				return fmt.Errorf("failed to delete reports: %w", err)
			}
		}
		// delete actions if exists
		if len(actions) > 0 {
			_, err = i.adminUnitOfWork.ActionRepository(ctx).DeleteByProblemID(ctx, problemID)
			if err != nil {
				return fmt.Errorf("failed to delete actions: %w", err)
			}
		}
		// delete hearing messages first (to avoid foreign key constraint violation)
		if len(hearings) > 0 {
			for _, hearing := range hearings {
				numDeleted, err := i.adminUnitOfWork.HearingMessageRepository(ctx).DeleteByHearingID(ctx, hearing.GetID())
				logger.Info("deleted hearing messages", "hearing_id", hearing.GetID(), "num_deleted", numDeleted)
				if err != nil {
					return fmt.Errorf("failed to delete hearing messages: %w", err)
				}
			}
		}

		// delete hearing
		if len(hearings) > 0 {
			_, err = i.adminUnitOfWork.HearingRepository(ctx).DeleteByProblemID(ctx, problem.GetID())
			if err != nil {
				return fmt.Errorf("failed to delete hearing: %w", err)
			}
		}

		// delete job config
		if jobConfig != nil {
			_, err = i.adminUnitOfWork.JobConfigRepository(ctx).DeleteByProblemID(ctx, problem.GetID())
			if err != nil {
				return fmt.Errorf("failed to delete job config: %w", err)
			}
		}

		// delete problem fields (hearing_messages may reference problem_fields)
		_, err = i.adminUnitOfWork.ProblemFieldRepository(ctx).DeleteByProblemID(ctx, problem.GetID())
		if err != nil {
			return fmt.Errorf("failed to delete problem fields: %w", err)
		}

		// delete problem
		numDeleted, err := i.adminUnitOfWork.ProblemRepository(ctx).Delete(ctx, problem.GetID())
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
