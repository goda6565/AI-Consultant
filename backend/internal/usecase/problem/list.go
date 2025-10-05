package problem

import (
	"context"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/repository"
)

type ListProblemInputPort interface {
	Execute(ctx context.Context) (*ListProblemOutput, error)
}

type ListProblemOutput struct {
	Problems []entity.Problem
}

type ListProblemInteractor struct {
	problemRepository repository.ProblemRepository
}

func NewListProblemUseCase(problemRepository repository.ProblemRepository) ListProblemInputPort {
	return &ListProblemInteractor{problemRepository: problemRepository}
}

func (i *ListProblemInteractor) Execute(ctx context.Context) (*ListProblemOutput, error) {
	problems, err := i.problemRepository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find problems: %w", err)
	}
	return &ListProblemOutput{Problems: problems}, nil
}
