package repository

import (
	"github.com/goda6565/ai-consultant/backend/internal/domain/document/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/document/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type DocumentRepository interface {
	FindAll() ([]entity.Document, error)
	FindById(id sharedValue.ID) (*entity.Document, error)
	FindByTitle(title value.Title) (*entity.Document, error)
	FindByPath(path value.StoragePath) (*entity.Document, error)
	FindByStatus(status value.DocumentStatus) ([]entity.Document, error)
	Create(document *entity.Document) error
	Delete(id sharedValue.ID) error
}
