package chunk

import (
	"fmt"
	"net/http"

	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	logger "github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/chunk"
	queuePorts "github.com/goda6565/ai-consultant/backend/internal/usecase/ports/queue"
	"github.com/labstack/echo/v4"
)

type CreateHandler echo.HandlerFunc

type createChunkHandler struct {
	createChunkUseCase chunk.CreateChunkInputPort
}

func NewCreateChunkHandler(createChunkUseCase chunk.CreateChunkInputPort) CreateHandler {
	h := &createChunkHandler{createChunkUseCase: createChunkUseCase}
	return CreateHandler(h.Handle)
}

func (h *createChunkHandler) Handle(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logger.GetLogger(ctx)
	var req queuePorts.SyncQueueMessage
	if err := c.Bind(&req); err != nil {
		logger.Error("failed to bind body", "error", err)
		return errors.NewInfrastructureError(errors.BadRequestError, "failed to bind body")
	}

	if _, err := h.createChunkUseCase.Execute(ctx, chunk.CreateChunkUseCaseInput{DocumentID: req.DocumentID}); err != nil {
		logger.Error("failed to create chunk", "error", err)
		return fmt.Errorf("failed to create chunk: %w", err)
	}

	return c.NoContent(http.StatusNoContent)
}
