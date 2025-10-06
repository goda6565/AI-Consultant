package repository

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing_map/entity"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type HearingMapRepository interface {
	Create(ctx context.Context, hearingMap *entity.HearingMap) error
	DeleteByHearingID(ctx context.Context, hearingID sharedValue.ID) (numDeleted int64, err error)
	FindByHearingID(ctx context.Context, hearingID sharedValue.ID) (*entity.HearingMap, error)
}
