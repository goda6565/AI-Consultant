package problem

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/entity"
	gen "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/problem"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type CreateProblemHandler struct {
	createProblemUseCase problem.CreateProblemInputPort
}

func NewCreateProblemHandler(createProblemUseCase problem.CreateProblemInputPort) *CreateProblemHandler {
	return &CreateProblemHandler{createProblemUseCase: createProblemUseCase}
}

func (h *CreateProblemHandler) CreateProblem(ctx context.Context, request gen.CreateProblemRequestObject) (gen.CreateProblemResponseObject, error) {
	createProblemOutput, err := h.createProblemUseCase.Execute(ctx, problem.CreateProblemUseCaseInput{
		Description: request.Body.Description,
	})
	if err != nil {
		return nil, err
	}
	return toCreateProblemJSONResponse(createProblemOutput.Problem), nil
}

func toCreateProblemJSONResponse(problem *entity.Problem) gen.CreateProblemResponseObject {
	return gen.CreateProblem201JSONResponse{
		CreateProblemSuccessJSONResponse: gen.CreateProblemSuccessJSONResponse{
			Id: openapi_types.UUID(uuid.MustParse(problem.GetID().Value())),
		},
	}
}
