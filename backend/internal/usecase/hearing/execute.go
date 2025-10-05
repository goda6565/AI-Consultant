package hearing

import (
	"context"
	"fmt"

	hearingEntity "github.com/goda6565/ai-consultant/backend/internal/domain/hearing/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing/repository"
	hearingService "github.com/goda6565/ai-consultant/backend/internal/domain/hearing/service"
	hearingMessageEntity "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/entity"
	hearingMessageRepository "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/repository"
	hearingMessageService "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/service"
	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/value"
	problemEntity "github.com/goda6565/ai-consultant/backend/internal/domain/problem/entity"
	problemRepository "github.com/goda6565/ai-consultant/backend/internal/domain/problem/repository"
	problemValue "github.com/goda6565/ai-consultant/backend/internal/domain/problem/value"
	problemFieldEntity "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/entity"
	problemFieldRepository "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/repository"
	problemFieldService "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/service"
	problemFieldValue "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/uuid"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/errors"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/ports/job"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/ports/transaction"
)

const IsCompletedMessage = "ヒアリングが完了しました。"

type ExecuteHearingInputPort interface {
	Execute(ctx context.Context, input ExecuteHearingUseCaseInput, logger logger.Logger) (*ExecuteHearingOutput, error)
}

type ExecuteHearingUseCaseInput struct {
	ProblemID   string
	HearingID   string
	UserMessage *string
}

type ExecuteHearingOutput struct {
	AssistantMessage string
	IsCompleted      bool
}

type ExecuteHearingInteractor struct {
	hearingRepository                  repository.HearingRepository
	hearingMessageRepository           hearingMessageRepository.HearingMessageRepository
	problemRepository                  problemRepository.ProblemRepository
	problemFieldRepository             problemFieldRepository.ProblemFieldRepository
	duplicateCheckerService            *hearingService.DuplicateCheckerService
	generateHearingMessageService      *hearingMessageService.GenerateHearingMessageService
	judgeProblemFieldCompletionService *problemFieldService.JudgeProblemFieldCompletionService
	adminUnitOfWork                    transaction.AdminUnitOfWork
	jobClient                          job.Job
	env                                *environment.Environment
}

func NewExecuteHearingUseCase(
	hearingRepository repository.HearingRepository,
	hearingMessageRepository hearingMessageRepository.HearingMessageRepository,
	problemRepository problemRepository.ProblemRepository,
	problemFieldRepository problemFieldRepository.ProblemFieldRepository,
	duplicateCheckerService *hearingService.DuplicateCheckerService,
	generateHearingMessageService *hearingMessageService.GenerateHearingMessageService,
	judgeProblemFieldCompletionService *problemFieldService.JudgeProblemFieldCompletionService,
	adminUnitOfWork transaction.AdminUnitOfWork,
	jobClient job.Job,
	env *environment.Environment,
) ExecuteHearingInputPort {
	return &ExecuteHearingInteractor{
		hearingRepository:                  hearingRepository,
		hearingMessageRepository:           hearingMessageRepository,
		problemRepository:                  problemRepository,
		problemFieldRepository:             problemFieldRepository,
		duplicateCheckerService:            duplicateCheckerService,
		generateHearingMessageService:      generateHearingMessageService,
		judgeProblemFieldCompletionService: judgeProblemFieldCompletionService,
		adminUnitOfWork:                    adminUnitOfWork,
		jobClient:                          jobClient,
		env:                                env,
	}
}

func (i *ExecuteHearingInteractor) Execute(ctx context.Context, input ExecuteHearingUseCaseInput, logger logger.Logger) (*ExecuteHearingOutput, error) {
	// validate and create problem ID
	problemID, err := sharedValue.NewID(input.ProblemID)
	if err != nil {
		return nil, fmt.Errorf("failed to create problem id: %w", err)
	}

	// validate and create hearing ID
	hearingID, err := sharedValue.NewID(input.HearingID)
	if err != nil {
		return nil, fmt.Errorf("failed to create hearing id: %w", err)
	}

	// pre-fetch problem, problem fields, and hearing
	problem, problemFields, hearing, hearingMessages, err := i.preFetch(ctx, problemID, hearingID)
	if err != nil {
		return nil, fmt.Errorf("failed to pre-fetch: %w", err)
	}

	var targetProblemFieldID sharedValue.ID
	if len(hearingMessages) == 0 {
		targetProblemFieldID = problemFields[0].GetID() // random select first problem field
	} else {
		// find hearing messages by hearing ID
		hearingMessages, err = i.hearingMessageRepository.FindByHearingID(ctx, hearing.GetID())
		if err != nil {
			return nil, fmt.Errorf("failed to find hearing messages: %w", err)
		}
		// latest hearing message is target problem field ID
		latestHearingMessage := hearingMessages[len(hearingMessages)-1]
		targetProblemFieldID = latestHearingMessage.GetProblemFieldID()
		// require user message
		if input.UserMessage == nil {
			return nil, errors.NewUseCaseError(errors.InternalError, "user message is required")
		}
		// save hearing message
		saved, err := i.saveHearingMessage(ctx, hearing.GetID(), targetProblemFieldID, value.RoleUser, *input.UserMessage)
		if err != nil {
			return nil, fmt.Errorf("failed to save hearing message: %w", err)
		}
		// reflect new user message in history before generation
		hearingMessages = append(hearingMessages, *saved)
	}

	// judge problem field completion
	judgeProblemFieldCompletionOutput, err := i.judgeProblemFieldCompletionService.Execute(ctx, problemFieldService.JudgeProblemFieldCompletionServiceInput{
		Problem:              *problem,
		HearingMessages:      hearingMessages,
		TargetProblemFieldID: targetProblemFieldID,
		ProblemFields:        problemFields,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to judge problem field completion: %w", err)
	}

	if judgeProblemFieldCompletionOutput.IsTargetProblemFieldAnswered {
		// update problem field if target problem field is answered
		err = i.problemFieldRepository.UpdateAnswered(ctx, targetProblemFieldID, *problemFieldValue.NewAnswered(true))
		if err != nil {
			return nil, fmt.Errorf("failed to update problem field: %w", err)
		}

		// find problem fields by problem ID
		newProblemFields, err := i.problemFieldRepository.FindByProblemID(ctx, problemID)
		if err != nil {
			return nil, fmt.Errorf("failed to find problem fields: %w", err)
		}

		// check if all problem fields are answered
		notAnsweredFields := []problemFieldEntity.ProblemField{}
		for _, newProblemField := range newProblemFields {
			answered := newProblemField.GetAnswered()
			if !answered.Value() {
				notAnsweredFields = append(notAnsweredFields, newProblemField)
			}
		}
		if len(notAnsweredFields) == 0 {
			id, err := sharedValue.NewID(uuid.NewUUID())
			if err != nil {
				return nil, fmt.Errorf("failed to create id: %w", err)
			}
			messageValue, err := value.NewMessage(IsCompletedMessage)
			if err != nil {
				return nil, fmt.Errorf("failed to create message: %w", err)
			}
			err = i.jobClient.CallJob(ctx, job.JobInput{
				JobName:   i.env.JobName,
				ProblemID: problemID.Value(),
			})
			if err != nil {
				return nil, fmt.Errorf("failed to call job: %w", err)
			}
			hearingMessage := hearingMessageEntity.NewHearingMessage(id, hearing.GetID(), targetProblemFieldID, value.RoleAssistant, *messageValue, nil)
			// update problem status
			err = i.adminUnitOfWork.WithTx(ctx, func(txCtx context.Context) error {
				// save hearing message
				err = i.adminUnitOfWork.HearingMessageRepository(txCtx).Create(txCtx, hearingMessage)
				if err != nil {
					return fmt.Errorf("failed to save hearing message: %w", err)
				}
				// update problem status
				err = i.adminUnitOfWork.ProblemRepository(txCtx).UpdateStatus(txCtx, problemID, problemValue.StatusProcessing)
				if err != nil {
					return fmt.Errorf("failed to update problem status: %w", err)
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("failed to update problem status: %w", err)
			}
			return &ExecuteHearingOutput{
				IsCompleted:      true,
				AssistantMessage: IsCompletedMessage,
			}, nil
		} else {
			// find next problem field
			nextProblemFieldID := notAnsweredFields[0].GetID()
			targetProblemFieldID = nextProblemFieldID
		}
	}

	// create hearing message
	generateHearingMessageOutput, err := i.generateHearingMessageService.Execute(ctx, hearingMessageService.GenerateHearingMessageInput{
		Problem:              problem,
		HearingMessages:      hearingMessages,
		TargetProblemFieldID: targetProblemFieldID,
		ProblemFields:        problemFields,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate hearing message: %w", err)
	}

	// save assistant message
	_, err = i.saveHearingMessage(ctx, hearing.GetID(), targetProblemFieldID, value.RoleAssistant, generateHearingMessageOutput.AssistantMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to save assistant message: %w", err)
	}

	return &ExecuteHearingOutput{
		IsCompleted:      false,
		AssistantMessage: generateHearingMessageOutput.AssistantMessage,
	}, nil
}

func (i *ExecuteHearingInteractor) preFetch(ctx context.Context, problemID sharedValue.ID, hearingID sharedValue.ID) (*problemEntity.Problem, []problemFieldEntity.ProblemField, *hearingEntity.Hearing, []hearingMessageEntity.HearingMessage, error) {
	// find problem by problem ID
	problem, err := i.problemRepository.FindById(ctx, problemID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to find problem: %w", err)
	}
	if problem == nil {
		return nil, nil, nil, nil, errors.NewUseCaseError(errors.NotFoundError, "problem not found")
	}

	// find problem fields by problem ID
	problemFields, err := i.problemFieldRepository.FindByProblemID(ctx, problemID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to find problem fields: %w", err)
	}
	if len(problemFields) == 0 {
		return nil, nil, nil, nil, errors.NewUseCaseError(errors.NotFoundError, "problem fields not found")
	}

	// find hearing by problem ID
	hearing, err := i.hearingRepository.FindById(ctx, hearingID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to find hearing: %w", err)
	}
	if hearing == nil {
		return nil, nil, nil, nil, errors.NewUseCaseError(errors.NotFoundError, "hearing not found")
	}

	// find hearing messages by hearing ID
	hearingMessages, err := i.hearingMessageRepository.FindByHearingID(ctx, hearing.GetID())
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to find hearing messages: %w", err)
	}

	return problem, problemFields, hearing, hearingMessages, nil
}

func (i *ExecuteHearingInteractor) saveHearingMessage(ctx context.Context, hearingID sharedValue.ID, problemFieldID sharedValue.ID, role value.Role, message string) (*hearingMessageEntity.HearingMessage, error) {
	id, err := sharedValue.NewID(uuid.NewUUID())
	if err != nil {
		return nil, fmt.Errorf("failed to create id: %w", err)
	}
	messageValue, err := value.NewMessage(message)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}
	hearingMessage := hearingMessageEntity.NewHearingMessage(id, hearingID, problemFieldID, role, *messageValue, nil)
	err = i.hearingMessageRepository.Create(ctx, hearingMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to save hearing message: %w", err)
	}
	return hearingMessage, nil
}
