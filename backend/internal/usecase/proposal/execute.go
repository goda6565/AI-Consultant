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
	eventEntity "github.com/goda6565/ai-consultant/backend/internal/domain/event/entity"
	eventRepository "github.com/goda6565/ai-consultant/backend/internal/domain/event/repository"
	eventValue "github.com/goda6565/ai-consultant/backend/internal/domain/event/value"
	hearingRepository "github.com/goda6565/ai-consultant/backend/internal/domain/hearing/repository"
	hearingMessageEntity "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/entity"
	hearingMessageRepository "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/repository"
	jobConfigEntity "github.com/goda6565/ai-consultant/backend/internal/domain/job_config/entity"
	jobConfigRepository "github.com/goda6565/ai-consultant/backend/internal/domain/job_config/repository"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
	problemEntity "github.com/goda6565/ai-consultant/backend/internal/domain/problem/entity"
	problemRepository "github.com/goda6565/ai-consultant/backend/internal/domain/problem/repository"
	problemValue "github.com/goda6565/ai-consultant/backend/internal/domain/problem/value"
	problemFieldEntity "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/entity"
	problemFieldRepository "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/repository"
	reportEntity "github.com/goda6565/ai-consultant/backend/internal/domain/report/entity"
	reportRepository "github.com/goda6565/ai-consultant/backend/internal/domain/report/repository"
	reportValue "github.com/goda6565/ai-consultant/backend/internal/domain/report/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/uuid"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/errors"
)

var ActionMessage = "%sを実行開始"

const MaxAllowedFailures = 3

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
	eventRepository          eventRepository.EventRepository
	orchestrator             *agentService.Orchestrator
	summarizeService         *agentService.SummarizeService
	goalService              *agentService.GoalService
	terminator               *agentService.Terminator
	actionFactory            *actionService.ActionFactory
	reportRepository         reportRepository.ReportRepository
	jobConfigRepository      jobConfigRepository.JobConfigRepository
}

func NewExecuteProposalUseCase(
	problemRepository problemRepository.ProblemRepository,
	problemFieldRepository problemFieldRepository.ProblemFieldRepository,
	hearingRepository hearingRepository.HearingRepository,
	hearingMessageRepository hearingMessageRepository.HearingMessageRepository,
	actionRepository actionRepository.ActionRepository,
	eventRepository eventRepository.EventRepository,
	orchestrator *agentService.Orchestrator,
	summarizeService *agentService.SummarizeService,
	goalService *agentService.GoalService,
	terminator *agentService.Terminator,
	actionFactory *actionService.ActionFactory,
	reportRepository reportRepository.ReportRepository,
	jobConfigRepository jobConfigRepository.JobConfigRepository,
) ExecuteProposalInputPort {
	return &ExecuteProposalInteractor{
		problemRepository:        problemRepository,
		problemFieldRepository:   problemFieldRepository,
		hearingRepository:        hearingRepository,
		hearingMessageRepository: hearingMessageRepository,
		actionRepository:         actionRepository,
		eventRepository:          eventRepository,
		orchestrator:             orchestrator,
		summarizeService:         summarizeService,
		goalService:              goalService,
		terminator:               terminator,
		actionFactory:            actionFactory,
		reportRepository:         reportRepository,
		jobConfigRepository:      jobConfigRepository,
	}
}

func (i *ExecuteProposalInteractor) Execute(ctx context.Context, input ExecuteProposalUseCaseInput) error {
	logger := logger.GetLogger(ctx)
	problemID, err := sharedValue.NewID(input.ProblemID)
	if err != nil {
		return fmt.Errorf("failed to create problem id: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case error:
				logger.Error("failed to execute proposal", "error", r)
			default:
				logger.Error("failed to execute proposal", "panic", r)
			}
			// update problem status
			err := i.problemRepository.UpdateStatus(ctx, problemID, problemValue.StatusFailed)
			if err != nil {
				logger.Error("failed to update problem status", "error", err)
			}
		}
	}()
	preFetchOutput, err := i.preFetch(ctx, problemID)
	if err != nil {
		return fmt.Errorf("failed to pre-fetch: %w", err)
	}
	problem := preFetchOutput.Problem
	problemFields := preFetchOutput.ProblemFields
	hearingMessages := preFetchOutput.HearingMessages
	jobConfig := preFetchOutput.JobConfig
	state := state.NewState(*problem, *value.NewContent(""), problemFields, hearingMessages, *value.NewHistory(""), []actionValue.ActionType{}, jobConfig.GetEnableInternalSearch())
	// goal
	goal, err := i.goalService.Execute(ctx, agentService.GoalServiceInput{State: *state})
	if err != nil {
		return fmt.Errorf("failed to execute goal: %w", err)
	}
	logger.Debug("goal", "goal", goal.Goal.Value())
	state.SetGoal(goal.Goal)

	var failCount int

	for {
		// orchestrator
		decision, err := i.orchestrator.Execute(ctx, agentService.OrchestratorInput{State: *state})
		logger.Debug("nextAction", "canProceed", decision.CanProceed)
		logger.Debug("nextAction", "reason", decision.Reason)
		if err != nil {
			return fmt.Errorf("failed to execute orchestrator: %w", err)
		}

		if decision.CanProceed && state.GetCurrentAction() == actionValue.ActionTypeReview {
			state.IncrementActionLoopCount()
			terminatorOutput, err := i.terminator.Execute(ctx, agentService.TerminatorInput{State: *state})
			if err != nil {
				return fmt.Errorf("failed to execute terminator: %w", err)
			}
			logger.Debug("terminatorOutput", "shouldTerminate", terminatorOutput.ShouldTerminate)
			logger.Debug("terminatorOutput", "reason", terminatorOutput.Reason)
			if terminatorOutput.ShouldTerminate {
				state.Done()
				err = i.createEvent(ctx, problemID, eventValue.EventTypeAction, state.GetCurrentAction(), "提案作成が完了しました。")
				if err != nil {
					return fmt.Errorf("failed to create event: %w", err)
				}
				break
			}
		}

		// action
		state.ToNextAction(decision.CanProceed)
		err = i.createEvent(ctx, problemID, eventValue.EventTypeAction, state.GetCurrentAction(), fmt.Sprintf(ActionMessage, state.GetCurrentAction().Value()))
		if err != nil {
			return fmt.Errorf("failed to create event: %w", err)
		}
		tmpl, err := i.actionFactory.GetActionTemplate(state.GetCurrentAction())
		if err != nil {
			return fmt.Errorf("failed to get action: %w", err)
		}
		output, err := tmpl.Execute(ctx, actionService.ActionTemplateInput{
			State: *state,
		})
		if err != nil {
			// agent must be alive even if some actions fail
			logger.Error("failed to execute action", "error", err)
			failCount++
			if failCount >= MaxAllowedFailures {
				logger.Error("failed to execute action", "error", err)
				return fmt.Errorf("failed to execute action: %w", err)
			}
			continue
		}
		// save action
		err = i.saveAction(ctx, output.Action)
		if err != nil {
			return fmt.Errorf("failed to save action: %w", err)
		}
		// save to state
		inputValue := output.Action.GetInput()
		outputValue := output.Action.GetOutput()
		logger.Debug("action", "action", output.Action.GetActionType().Value())
		logger.Debug("inputValue", "inputValue", inputValue.Value())
		logger.Debug("outputValue", "outputValue", outputValue.Value())
		if inputValue.Value() != "" {
			err = i.createEvent(ctx, problemID, eventValue.EventTypeInput, state.GetCurrentAction(), inputValue.Value())
			if err != nil {
				return fmt.Errorf("failed to create event: %w", err)
			}
		}
		if outputValue.Value() != "" {
			err = i.createEvent(ctx, problemID, eventValue.EventTypeOutput, state.GetCurrentAction(), outputValue.Value())
			if err != nil {
				return fmt.Errorf("failed to create event: %w", err)
			}
		}
		state.SetContent(output.Content)
		state.AddHistory(state.GetCurrentAction(), output.Action.ToHistory())
		// summarize
		history := state.GetHistory()
		summarizeNeeded, err := i.summarizeService.IsSummarizeNeeded(ctx, agentService.SummarizeServiceInput{
			History:   history.GetValue(),
			LLMConfig: llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		})
		if err != nil {
			return fmt.Errorf("failed to check if summarize is needed: %w", err)
		}
		if summarizeNeeded {
			// summarize
			summarizedHistory, err := i.summarizeService.Summarize(ctx, agentService.SummarizeServiceInput{
				History:   history.GetValue(),
				LLMConfig: llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
			})
			logger.Debug("summarizedHistory", "summarizedHistory", summarizedHistory.SummarizedHistory)
			if err != nil {
				return fmt.Errorf("failed to summarize history: %w", err)
			}
			state.SetHistory(*value.NewHistory(summarizedHistory.SummarizedHistory))
		}
	}
	// debug
	logger.Info("state", "state", state)

	// save report
	reportID, err := sharedValue.NewID(uuid.NewUUID())
	if err != nil {
		return fmt.Errorf("failed to create report id: %w", err)
	}
	content := state.GetContent()
	reportContent := reportValue.NewContent(content.Value())
	reportEntity := reportEntity.NewReport(reportID, problemID, *reportContent, nil)
	err = i.reportRepository.Create(ctx, reportEntity)
	if err != nil {
		return fmt.Errorf("failed to save report: %w", err)
	}

	// update problem status
	err = i.problemRepository.UpdateStatus(ctx, problemID, problemValue.StatusDone)
	if err != nil {
		return fmt.Errorf("failed to update problem status: %w", err)
	}
	return nil
}

type PreFetchOutput struct {
	Problem         *problemEntity.Problem
	ProblemFields   []problemFieldEntity.ProblemField
	HearingMessages []hearingMessageEntity.HearingMessage
	JobConfig       *jobConfigEntity.JobConfig
}

func (i *ExecuteProposalInteractor) preFetch(ctx context.Context, problemID sharedValue.ID) (*PreFetchOutput, error) {
	problem, err := i.problemRepository.FindById(ctx, problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to find problem: %w", err)
	}
	if problem == nil {
		return nil, errors.NewUseCaseError(errors.NotFoundError, "problem not found")
	}
	problemFields, err := i.problemFieldRepository.FindByProblemID(ctx, problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to find problem fields: %w", err)
	}
	if len(problemFields) == 0 {
		return nil, errors.NewUseCaseError(errors.NotFoundError, "problem fields not found")
	}
	hearing, err := i.hearingRepository.FindByProblemId(ctx, problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to find hearing: %w", err)
	}
	if hearing == nil {
		return nil, errors.NewUseCaseError(errors.NotFoundError, "hearing not found")
	}
	hearingMessages, err := i.hearingMessageRepository.FindByHearingID(ctx, hearing.GetID())
	if err != nil {
		return nil, fmt.Errorf("failed to find hearing messages: %w", err)
	}
	jobConfig, err := i.jobConfigRepository.FindByProblemID(ctx, problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to find job config: %w", err)
	}
	if jobConfig == nil {
		return nil, errors.NewUseCaseError(errors.NotFoundError, "job config not found")
	}
	return &PreFetchOutput{
		Problem:         problem,
		ProblemFields:   problemFields,
		HearingMessages: hearingMessages,
		JobConfig:       jobConfig,
	}, nil
}

func (i *ExecuteProposalInteractor) saveAction(ctx context.Context, action actionEntity.Action) error {
	err := i.actionRepository.Create(ctx, &action)
	if err != nil {
		return fmt.Errorf("failed to save action: %w", err)
	}
	return nil
}

func (i *ExecuteProposalInteractor) createEvent(ctx context.Context, problemID sharedValue.ID, eventType eventValue.EventType, actionType actionValue.ActionType, message string) error {
	logger := logger.GetLogger(ctx)
	id, err := sharedValue.NewID(uuid.NewUUID())
	if err != nil {
		return fmt.Errorf("failed to create id: %w", err)
	}
	messageValue, err := eventValue.NewMessage(message)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}
	err = i.eventRepository.Create(ctx, &eventEntity.Event{
		ID:         id,
		ProblemID:  problemID,
		EventType:  eventType,
		ActionType: actionType,
		Message:    *messageValue,
	})
	if err != nil {
		// ignore error
		logger.Error("failed to create event", "error", err)
		return nil
	}
	return nil
}
