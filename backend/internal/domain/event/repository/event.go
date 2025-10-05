package repository

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/event/entity"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type EventRepository interface {
	Create(ctx context.Context, event *entity.Event) error
	FindAllByProblemID(ctx context.Context, problemID sharedValue.ID) ([]entity.Event, error)
	FindAllByProblemIDAsStream(ctx context.Context, problemID sharedValue.ID) (<-chan entity.Event, error)
	DeleteAllByProblemID(ctx context.Context, problemID sharedValue.ID) (numDeleted int64, err error)
}
