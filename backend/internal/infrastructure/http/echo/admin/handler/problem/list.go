package problem

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/entity"
	gen "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/problem"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type ListProblemHandler struct {
	listProblemUseCase problem.ListProblemInputPort
}

func NewListProblemHandler(listProblemUseCase problem.ListProblemInputPort) *ListProblemHandler {
	return &ListProblemHandler{listProblemUseCase: listProblemUseCase}
}

func (h *ListProblemHandler) ListProblems(ctx context.Context, request gen.ListProblemsRequestObject) (gen.ListProblemsResponseObject, error) {
	listProblemOutput, err := h.listProblemUseCase.Execute(ctx)
	if err != nil {
		return nil, err
	}
	return toProblemsJSONResponse(listProblemOutput.Problems), nil
}

func toProblemsJSONResponse(problems []entity.Problem) gen.ListProblemsResponseObject {
	problemsJSON := make([]gen.Problem, len(problems))
	for i, problem := range problems {
		problemsJSON[i] = toSingleProblemJSON(&problem)
	}
	return gen.ListProblems200JSONResponse{
		ListProblemsSuccessJSONResponse: gen.ListProblemsSuccessJSONResponse{
			Problems: problemsJSON,
		},
	}
}

func toSingleProblemJSON(problem *entity.Problem) gen.Problem {
	return gen.Problem{
		Id:          openapi_types.UUID(uuid.MustParse(problem.GetID().Value())),
		Title:       problem.GetTitle().Value(),
		Description: problem.GetDescription().Value(),
		Status:      gen.ProblemStatus(problem.GetStatus().Value()),
		CreatedAt:   *problem.GetCreatedAt(),
	}
}
