package repository

import (
	"context"

	problemFieldEntity "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/entity"
	problemFieldValue "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type ProblemFieldRepository interface {
	FindByProblemID(ctx context.Context, problemID sharedValue.ID) ([]problemFieldEntity.ProblemField, error)
	Create(ctx context.Context, problemField *problemFieldEntity.ProblemField) error
	UpdateAnswered(ctx context.Context, id sharedValue.ID, answered problemFieldValue.Answered) error
	DeleteByProblemID(ctx context.Context, problemID sharedValue.ID) (numDeleted int64, err error)
}
