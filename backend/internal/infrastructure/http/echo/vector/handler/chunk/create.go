package chunk

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/goda6565/ai-consultant/backend/internal/pkg/logger"

	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/chunk"
	"github.com/labstack/echo/v4"
)

type CreateHandler echo.HandlerFunc

type pushRequest struct {
	Message struct {
		Data string `json:"data"`
	} `json:"message"`
}

type createChunkHandler struct {
	logger             logger.Logger
	createChunkUseCase chunk.CreateChunkInputPort
}

func NewCreateChunkHandler(createChunkUseCase chunk.CreateChunkInputPort, logger logger.Logger) CreateHandler {
	h := &createChunkHandler{createChunkUseCase: createChunkUseCase, logger: logger}
	return CreateHandler(h.Handle)
}

func (h *createChunkHandler) Handle(c echo.Context) error {
	var req pushRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("failed to bind body", "error", err)
		return errors.NewInfrastructureError(errors.BadRequestError, "failed to bind body")
	}

	decoded, err := base64.StdEncoding.DecodeString(req.Message.Data)
	if err != nil {
		h.logger.Error("failed to decode base64", "error", err)
		return errors.NewInfrastructureError(errors.BadRequestError, "failed to decode base64")
	}

	if _, err = h.createChunkUseCase.Execute(c.Request().Context(), chunk.CreateChunkUseCaseInput{DocumentID: string(decoded)}); err != nil {
		h.logger.Error("failed to create chunk", "error", err)
		return fmt.Errorf("failed to create chunk: %w", err)
	}

	return c.NoContent(http.StatusNoContent)
}
