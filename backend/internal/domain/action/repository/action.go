package repository

import (
	"context"

	actionEntity "github.com/goda6565/ai-consultant/backend/internal/domain/action/entity"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type ActionRepository interface {
	FindByProblemID(ctx context.Context, problemID sharedValue.ID) ([]actionEntity.Action, error)
	Create(ctx context.Context, action *actionEntity.Action) error
	DeleteByProblemID(ctx context.Context, problemID sharedValue.ID) (numDeleted int64, err error)
}
