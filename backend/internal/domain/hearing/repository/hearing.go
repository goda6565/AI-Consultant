package repository

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing/entity"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type HearingRepository interface {
	FindById(ctx context.Context, id sharedValue.ID) (*entity.Hearing, error)
	FindByProblemId(ctx context.Context, problemID sharedValue.ID) (*entity.Hearing, error)
	FindAllByProblemId(ctx context.Context, problemID sharedValue.ID) ([]entity.Hearing, error)
	Create(ctx context.Context, hearing *entity.Hearing) error
	DeleteByProblemID(ctx context.Context, problemID sharedValue.ID) (numDeleted int64, err error)
}
