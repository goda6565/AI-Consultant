package document

import (
	"context"
	"fmt"
	"github.com/goda6565/ai-consultant/backend/internal/domain/document/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/document/repository"
)

type ListDocumentInputPort interface {
	Execute(ctx context.Context) (*ListDocumentOutput, error)
}

type ListDocumentOutput struct {
	Documents []entity.Document
}

type ListDocumentInteractor struct {
	documentRepository repository.DocumentRepository
}

func NewListDocumentUseCase(documentRepository repository.DocumentRepository) ListDocumentInputPort {
	return &ListDocumentInteractor{documentRepository: documentRepository}
}

func (i *ListDocumentInteractor) Execute(ctx context.Context) (*ListDocumentOutput, error) {
	documents, err := i.documentRepository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	return &ListDocumentOutput{Documents: documents}, nil
}
