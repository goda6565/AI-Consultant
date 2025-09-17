package repository

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/entity"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type HearingMessageRepository interface {
	FindByHearingID(ctx context.Context, hearingID sharedValue.ID) ([]entity.HearingMessage, error)
	Create(ctx context.Context, hearingMessage *entity.HearingMessage) error
	DeleteByHearingID(ctx context.Context, hearingID sharedValue.ID) (numDeleted int64, err error)
}
