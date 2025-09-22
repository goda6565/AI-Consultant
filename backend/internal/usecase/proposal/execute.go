package proposal

import (
	"context"
	"fmt"

	problemRepository "github.com/goda6565/ai-consultant/backend/internal/domain/problem/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	stateEntity "github.com/goda6565/ai-consultant/backend/internal/domain/state/entity"
	stateService "github.com/goda6565/ai-consultant/backend/internal/domain/state/service"
	stateValue "github.com/goda6565/ai-consultant/backend/internal/domain/state/value"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/uuid"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/errors"
)

const InitialStatusMessage = "提案作成を開始します。"

type ExecuteProposalInputPort interface {
	Execute(ctx context.Context, input ExecuteProposalUseCaseInput) error
}

type ExecuteProposalUseCaseInput struct {
	ProblemID string
}

type ExecuteProposalInteractor struct {
	problemRepository      problemRepository.ProblemRepository
	hearingSummaryPhase    *stateService.HearingSummaryPhase
	problemDefinitionPhase *stateService.ProblemDefinitionPhase
}

func NewExecuteAgentUseCase(
	problemRepository problemRepository.ProblemRepository,
	hearingSummaryPhase *stateService.HearingSummaryPhase,
	problemDefinitionPhase *stateService.ProblemDefinitionPhase,
) ExecuteProposalInputPort {
	return &ExecuteProposalInteractor{
		problemRepository:      problemRepository,
		hearingSummaryPhase:    hearingSummaryPhase,
		problemDefinitionPhase: problemDefinitionPhase,
	}
}

func (i *ExecuteProposalInteractor) Execute(ctx context.Context, input ExecuteProposalUseCaseInput) error {
	logger := logger.GetLogger(ctx)

	// initialize state
	id, err := sharedValue.NewID(uuid.NewUUID())
	if err != nil {
		return fmt.Errorf("failed to create id: %w", err)
	}
	problemID, err := sharedValue.NewID(input.ProblemID)
	if err != nil {
		return fmt.Errorf("failed to create problem id: %w", err)
	}
	problem, err := i.problemRepository.FindById(ctx, problemID)
	if err != nil {
		return fmt.Errorf("failed to find problem: %w", err)
	}
	if problem == nil {
		return errors.NewUseCaseError(errors.NotFoundError, "problem not found")
	}
	statusMessage := stateValue.NewStatusMessage(InitialStatusMessage)
	state := stateEntity.NewState(id, problem, stateValue.InitialPhase, *statusMessage, nil)
	// execute hearing summary phase
	logger.Info("executing hearing summary phase")
	state.ProceedPhase()
	state, err = i.hearingSummaryPhase.Execute(ctx, state)
	if err != nil {
		return fmt.Errorf("failed to execute hearing summary phase: %w", err)
	}
	// execute problem definition phase
	logger.Info("executing problem definition phase")
	state.ProceedPhase()
	_, err = i.problemDefinitionPhase.Execute(ctx, state)
	if err != nil {
		return fmt.Errorf("failed to execute problem definition phase: %w", err)
	}
	return nil
}
