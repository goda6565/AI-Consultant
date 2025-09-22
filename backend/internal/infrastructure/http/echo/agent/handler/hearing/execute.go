package hearing

import (
	"context"
	"fmt"

	problemValue "github.com/goda6565/ai-consultant/backend/internal/domain/problem/value"
	infraErrors "github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	gen "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/agent/internal"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/hearing"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/problem"
)

type ExecuteHearingHandler struct {
	executeHearingUseCase hearing.ExecuteHearingInputPort
	getProblemUseCase     problem.GetProblemInputPort
}

func NewExecuteHearingHandler(executeHearingUseCase hearing.ExecuteHearingInputPort, getProblemUseCase problem.GetProblemInputPort) *ExecuteHearingHandler {
	return &ExecuteHearingHandler{
		executeHearingUseCase: executeHearingUseCase,
		getProblemUseCase:     getProblemUseCase,
	}
}

func (h *ExecuteHearingHandler) ExecuteHearing(ctx context.Context, request gen.ExecuteHearingRequestObject) (gen.ExecuteHearingResponseObject, error) {
	logger := logger.GetLogger(ctx)
	problemID := request.ProblemId.String()
	// find problem by problem id
	problem, err := h.getProblemUseCase.Execute(ctx, problem.GetProblemUseCaseInput{
		ProblemID: problemID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get problem: %w", err)
	}

	switch problem.Problem.GetStatus() {
	case problemValue.StatusProcessing:
		return nil, infraErrors.NewInfrastructureError(infraErrors.BadRequestError, "problem already processing")
	case problemValue.StatusDone:
		return nil, infraErrors.NewInfrastructureError(infraErrors.BadRequestError, "problem already done")
	}

	output, err := h.executeHearingUseCase.Execute(ctx, hearing.ExecuteHearingUseCaseInput{
		ProblemID:   problemID,
		UserMessage: request.Body.UserMessage,
	}, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to execute hearing: %w", err)
	}

	// response success
	return gen.ExecuteHearing200JSONResponse{
		ExecuteHearingSuccessJSONResponse: gen.ExecuteHearingSuccessJSONResponse{
			AssistantMessage: output.AssistantMessage,
			IsCompleted:      output.IsCompleted,
		},
	}, nil
}
