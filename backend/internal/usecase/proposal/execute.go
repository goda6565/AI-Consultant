package proposal

import (
	"context"
	"fmt"

	actionEntity "github.com/goda6565/ai-consultant/backend/internal/domain/action/entity"
	actionRepository "github.com/goda6565/ai-consultant/backend/internal/domain/action/repository"
	actionService "github.com/goda6565/ai-consultant/backend/internal/domain/action/service"
	actionValue "github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	agentService "github.com/goda6565/ai-consultant/backend/internal/domain/agent/service"
	"github.com/goda6565/ai-consultant/backend/internal/domain/agent/state"
	"github.com/goda6565/ai-consultant/backend/internal/domain/agent/value"
	hearingRepository "github.com/goda6565/ai-consultant/backend/internal/domain/hearing/repository"
	hearingMessageEntity "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/entity"
	hearingMessageRepository "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/repository"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
	problemEntity "github.com/goda6565/ai-consultant/backend/internal/domain/problem/entity"
	problemRepository "github.com/goda6565/ai-consultant/backend/internal/domain/problem/repository"
	problemFieldEntity "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/entity"
	problemFieldRepository "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
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
	problemRepository        problemRepository.ProblemRepository
	problemFieldRepository   problemFieldRepository.ProblemFieldRepository
	hearingRepository        hearingRepository.HearingRepository
	hearingMessageRepository hearingMessageRepository.HearingMessageRepository
	actionRepository         actionRepository.ActionRepository
	orchestrator             *agentService.Orchestrator
	summarizeService         *agentService.SummarizeService
	actionFactory            *actionService.ActionFactory
}

func NewExecuteProposalUseCase(
	problemRepository problemRepository.ProblemRepository,
	problemFieldRepository problemFieldRepository.ProblemFieldRepository,
	hearingRepository hearingRepository.HearingRepository,
	hearingMessageRepository hearingMessageRepository.HearingMessageRepository,
	actionRepository actionRepository.ActionRepository,
	orchestrator *agentService.Orchestrator,
	summarizeService *agentService.SummarizeService,
	actionFactory *actionService.ActionFactory,
) ExecuteProposalInputPort {
	return &ExecuteProposalInteractor{
		problemRepository:        problemRepository,
		problemFieldRepository:   problemFieldRepository,
		hearingRepository:        hearingRepository,
		hearingMessageRepository: hearingMessageRepository,
		actionRepository:         actionRepository,
		orchestrator:             orchestrator,
		summarizeService:         summarizeService,
		actionFactory:            actionFactory,
	}
}

func (i *ExecuteProposalInteractor) Execute(ctx context.Context, input ExecuteProposalUseCaseInput) error {
	logger := logger.GetLogger(ctx)
	problemID, err := sharedValue.NewID(input.ProblemID)
	if err != nil {
		return fmt.Errorf("failed to create problem id: %w", err)
	}
	problem, problemFields, hearingMessages, err := i.preFetch(ctx, problemID)
	if err != nil {
		return fmt.Errorf("failed to pre-fetch: %w", err)
	}
	state := state.NewState(*problem, *value.NewContent(""), problemFields, hearingMessages, *value.NewHistory(""))
	for {
		nextAction, err := i.orchestrator.Execute(ctx, agentService.OrchestratorInput{State: *state})
		logger.Debug("nextAction", "nextAction", nextAction.NextAction.Value())
		if err != nil {
			return fmt.Errorf("failed to execute orchestrator: %w", err)
		}
		if nextAction.NextAction.Equals(actionValue.ActionTypeDone) {
			break
		}
		tmpl, err := i.actionFactory.GetActionTemplate(nextAction.NextAction)
		if err != nil {
			return fmt.Errorf("failed to get action: %w", err)
		}
		output, err := tmpl.Execute(ctx, actionService.ActionTemplateInput{
			State: *state,
		})
		if err != nil {
			return fmt.Errorf("failed to execute action: %w", err)
		}
		err = i.saveAction(ctx, output.Action)
		if err != nil {
			return fmt.Errorf("failed to save action: %w", err)
		}
		state.SetContent(output.Content)
		state.AddHistory(nextAction.NextAction, output.Action.ToHistory())
		history := state.GetHistory()
		logger.Debug("history", "history", history.GetValue())
		summarizeNeeded, err := i.summarizeService.IsSummarizeNeeded(ctx, agentService.SummarizeServiceInput{
			History:   history.GetValue(),
			LLMConfig: llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		})
		if err != nil {
			return fmt.Errorf("failed to check if summarize is needed: %w", err)
		}
		if summarizeNeeded {
			summarizedHistory, err := i.summarizeService.Summarize(ctx, agentService.SummarizeServiceInput{
				History:   history.GetValue(),
				LLMConfig: llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
			})
			if err != nil {
				return fmt.Errorf("failed to summarize history: %w", err)
			}
			state.SetHistory(*value.NewHistory(summarizedHistory.SummarizedHistory))
		}
	}
	// debug
	logger.Info("state", "state", state)

	// update problem status
	// err = i.problemRepository.UpdateStatus(ctx, problemID, problemValue.StatusDone)
	// if err != nil {
	// 	return fmt.Errorf("failed to update problem status: %w", err)
	// }
	return nil
}

func (i *ExecuteProposalInteractor) preFetch(ctx context.Context, problemID sharedValue.ID) (*problemEntity.Problem, []problemFieldEntity.ProblemField, []hearingMessageEntity.HearingMessage, error) {
	problem, err := i.problemRepository.FindById(ctx, problemID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to find problem: %w", err)
	}
	if problem == nil {
		return nil, nil, nil, errors.NewUseCaseError(errors.NotFoundError, "problem not found")
	}
	problemFields, err := i.problemFieldRepository.FindByProblemID(ctx, problemID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to find problem fields: %w", err)
	}
	if len(problemFields) == 0 {
		return nil, nil, nil, errors.NewUseCaseError(errors.NotFoundError, "problem fields not found")
	}
	hearing, err := i.hearingRepository.FindByProblemId(ctx, problemID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to find hearing: %w", err)
	}
	if hearing == nil {
		return nil, nil, nil, errors.NewUseCaseError(errors.NotFoundError, "hearing not found")
	}
	hearingMessages, err := i.hearingMessageRepository.FindByHearingID(ctx, hearing.GetID())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to find hearing messages: %w", err)
	}
	return problem, problemFields, hearingMessages, nil
}

func (i *ExecuteProposalInteractor) saveAction(ctx context.Context, action actionEntity.Action) error {
	err := i.actionRepository.Create(ctx, &action)
	if err != nil {
		return fmt.Errorf("failed to save action: %w", err)
	}
	return nil
}
