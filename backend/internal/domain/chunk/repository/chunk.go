package repository

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/chunk/entity"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type ChunkRepository interface {
	Create(ctx context.Context, chunk *entity.Chunk) error
	Delete(ctx context.Context, documentID sharedValue.ID) (numDeleted int64, err error)
}
