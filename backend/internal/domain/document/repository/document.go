package repository

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/document/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/document/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type DocumentRepository interface {
	FindAll(ctx context.Context) ([]entity.Document, error)
	FindById(ctx context.Context, id sharedValue.ID) (*entity.Document, error)
	FindByTitle(ctx context.Context, title value.Title) (*entity.Document, error)
	Create(ctx context.Context, document *entity.Document) error
	Update(ctx context.Context, document *entity.Document) (numUpdated int64, err error)
	Delete(ctx context.Context, id sharedValue.ID) (numDeleted int64, err error)
}
