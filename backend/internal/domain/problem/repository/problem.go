package repository

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/entity"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type ProblemRepository interface {
	FindAll(ctx context.Context) ([]entity.Problem, error)
	FindById(ctx context.Context, id sharedValue.ID) (*entity.Problem, error)
	Create(ctx context.Context, problem *entity.Problem) error
	Delete(ctx context.Context, id sharedValue.ID) (numDeleted int64, err error)
}
