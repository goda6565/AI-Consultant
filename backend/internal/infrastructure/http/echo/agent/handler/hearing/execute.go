package hearing

import (
	"errors"

	domainErrors "github.com/goda6565/ai-consultant/backend/internal/domain/errors"
	problemValue "github.com/goda6565/ai-consultant/backend/internal/domain/problem/value"
	infraErrors "github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	logger "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/shared/middleware"
	usecaseErrors "github.com/goda6565/ai-consultant/backend/internal/usecase/errors"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/hearing"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/problem"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{}

type ExecuteHearingHandler echo.HandlerFunc

type executeHearingHandler struct {
	executeHearingUseCase hearing.ExecuteHearingInputPort
	getProblemUseCase     problem.GetProblemInputPort
}

func NewExecuteHearingHandler(executeHearingUseCase hearing.ExecuteHearingInputPort, getProblemUseCase problem.GetProblemInputPort) ExecuteHearingHandler {
	h := &executeHearingHandler{
		executeHearingUseCase: executeHearingUseCase,
		getProblemUseCase:     getProblemUseCase,
	}
	return h.Handle
}

type WsRequest struct {
	UserMessage string `json:"user_message"`
}

type WsResponseType string

const (
	WsResponseTypeHearingResponse  WsResponseType = "hearing_response"
	WsResponseTypeHearingCompleted WsResponseType = "hearing_completed"
	WsResponseTypeError            WsResponseType = "error"
)

type WsResponse struct {
	Type             string `json:"type"`
	AssistantMessage string `json:"assistant_message"`
}

func (h *executeHearingHandler) Handle(c echo.Context) (err error) {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return infraErrors.NewInfrastructureError(infraErrors.ExternalServiceError, "failed to upgrade websocket")
	}
	ctx := c.Request().Context()
	logger := logger.GetLogger(ctx)
	defer func() {
		if err := ws.Close(); err != nil {
			logger.Error("failed to close websocket", "error", err)
		}
	}()

	// Handle panics and errors in one place
	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic occurred in websocket handler", "panic", r)
			writeError(ws, infraErrors.NewInfrastructureError(infraErrors.ExternalServiceError, "internal server error"))
		}
		if err != nil {
			logger.Error("error occurred in websocket handler", "error", err)
			writeError(ws, err)
		}
	}()

	// validate problem id
	problemID := c.Param("problemId")

	// find problem by problem id
	problem, err := h.getProblemUseCase.Execute(ctx, problem.GetProblemUseCaseInput{
		ProblemID: problemID,
	})
	if err != nil {
		return infraErrors.NewInfrastructureError(infraErrors.BadRequestError, "failed to find problem")
	}
	if problem == nil {
		return infraErrors.NewInfrastructureError(infraErrors.BadRequestError, "problem not found")
	}

	var isFirst bool
	switch problem.Problem.GetStatus() {
	case problemValue.StatusPending:
		isFirst = true
	case problemValue.StatusHearing:
		isFirst = false
	case problemValue.StatusProcessing:
		return infraErrors.NewInfrastructureError(infraErrors.BadRequestError, "problem already processing")
	case problemValue.StatusDone:
		return infraErrors.NewInfrastructureError(infraErrors.BadRequestError, "problem already done")
	}

	for {
		var req WsRequest
		if !isFirst {
			err := ws.ReadJSON(&req)
			if err != nil {
				return err
			}
		}

		output, err := h.executeHearingUseCase.Execute(ctx, hearing.ExecuteHearingUseCaseInput{
			ProblemID:   problemID,
			UserMessage: &req.UserMessage,
		}, logger)
		if err != nil {
			return err
		}

		// response success
		h.writeSuccess(ws, output.AssistantMessage)

		// if hearing is completed, close the connection
		if output.IsCompleted {
			h.writeCompleted(ws, output.AssistantMessage)
			break
		}
		// set isFirst to false
		isFirst = false
	}

	return nil
}

func (h *executeHearingHandler) writeSuccess(ws *websocket.Conn, assistantMessage string) {
	successMsg := WsResponse{
		Type:             string(WsResponseTypeHearingResponse),
		AssistantMessage: assistantMessage,
	}
	_ = ws.WriteJSON(successMsg)
}

func (h *executeHearingHandler) writeCompleted(ws *websocket.Conn, assistantMessage string) {
	completedMsg := WsResponse{
		Type:             string(WsResponseTypeHearingCompleted),
		AssistantMessage: assistantMessage,
	}
	_ = ws.WriteJSON(completedMsg)
}

func writeError(ws *websocket.Conn, err error) {
	var message string
	var infraErr *infraErrors.InfrastructureError
	var domainErr *domainErrors.DomainError
	var usecaseErr *usecaseErrors.UseCaseError

	if errors.As(err, &infraErr) {
		message = infraErrToWsMessage(infraErr)
	} else if errors.As(err, &domainErr) {
		message = domainErrToWsMessage(domainErr)
	} else if errors.As(err, &usecaseErr) {
		message = usecaseErrToWsMessage(usecaseErr)
	} else {
		message = InternalErrorMessage
	}

	errorMsg := WsResponse{
		Type:             string(WsResponseTypeError),
		AssistantMessage: message,
	}
	_ = ws.WriteJSON(errorMsg)
}

const (
	// Common error messages
	InternalErrorMessage   = "Internal server error"
	ExternalServiceMessage = "External service error"
	ValidationErrorMessage = "Invalid input"
	UnauthorizedMessage    = "Authentication required"
	ForbiddenMessage       = "Access denied"
	NotFoundMessage        = "Resource not found"
	DuplicateErrorMessage  = "Resource already exists"
	ProcessingErrorMessage = "Processing error occurred"
	UnexpectedErrorMessage = "Unexpected error occurred"
)

func infraErrToWsMessage(err *infraErrors.InfrastructureError) string {
	switch err.ErrorType {
	case infraErrors.ExternalServiceError:
		return ExternalServiceMessage
	case infraErrors.UnauthorizedError:
		return UnauthorizedMessage
	case infraErrors.ForbiddenError:
		return ForbiddenMessage
	case infraErrors.BadRequestError:
		return ValidationErrorMessage
	default:
		return InternalErrorMessage
	}
}

func domainErrToWsMessage(err *domainErrors.DomainError) string {
	switch err.ErrorType {
	case domainErrors.ValidationError:
		return ValidationErrorMessage
	default:
		return InternalErrorMessage
	}
}

func usecaseErrToWsMessage(err *usecaseErrors.UseCaseError) string {
	switch err.ErrorType {
	case usecaseErrors.DuplicateError:
		return DuplicateErrorMessage
	case usecaseErrors.NotFoundError:
		return NotFoundMessage
	case usecaseErrors.InternalError:
		return ProcessingErrorMessage
	default:
		return UnexpectedErrorMessage
	}
}
