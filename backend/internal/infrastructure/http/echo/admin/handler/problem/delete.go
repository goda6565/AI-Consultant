package problem

import (
	"context"

	gen "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/problem"
)

type DeleteProblemHandler struct {
	deleteProblemUseCase problem.DeleteProblemInputPort
}

func NewDeleteProblemHandler(deleteProblemUseCase problem.DeleteProblemInputPort) *DeleteProblemHandler {
	return &DeleteProblemHandler{deleteProblemUseCase: deleteProblemUseCase}
}

func (h *DeleteProblemHandler) DeleteProblem(ctx context.Context, request gen.DeleteProblemRequestObject) (gen.DeleteProblemResponseObject, error) {
	err := h.deleteProblemUseCase.Execute(ctx, problem.DeleteProblemUseCaseInput{ProblemID: request.ProblemId.String()})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
