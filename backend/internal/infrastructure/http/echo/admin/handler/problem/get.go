package problem

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/entity"
	gen "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/problem"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type GetProblemHandler struct {
	getProblemUseCase problem.GetProblemInputPort
}

func NewGetProblemHandler(getProblemUseCase problem.GetProblemInputPort) *GetProblemHandler {
	return &GetProblemHandler{getProblemUseCase: getProblemUseCase}
}

func (h *GetProblemHandler) GetProblem(ctx context.Context, request gen.GetProblemRequestObject) (gen.GetProblemResponseObject, error) {
	getProblemOutput, err := h.getProblemUseCase.Execute(ctx, problem.GetProblemUseCaseInput{ProblemID: request.ProblemId.String()})
	if err != nil {
		return nil, err
	}
	return toProblemJSONResponse(getProblemOutput.Problem), nil
}

func toProblemJSONResponse(problem *entity.Problem) gen.GetProblemResponseObject {
	return gen.GetProblem200JSONResponse{
		GetProblemSuccessJSONResponse: gen.GetProblemSuccessJSONResponse{
			Id:          openapi_types.UUID(uuid.MustParse(problem.GetID().Value())),
			Title:       problem.GetTitle().Value(),
			Description: problem.GetDescription().Value(),
			Status:      gen.ProblemStatus(problem.GetStatus().Value()),
			CreatedAt:   *problem.GetCreatedAt(),
		},
	}
}
