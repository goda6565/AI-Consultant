package document

import (
	"context"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/document/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/document/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/errors"
)

type GetDocumentInputPort interface {
	Execute(ctx context.Context, input GetDocumentUseCaseInput) (*GetDocumentOutput, error)
}

type GetDocumentUseCaseInput struct {
	DocumentID string
}

type GetDocumentOutput struct {
	Document *entity.Document
}

type GetDocumentInteractor struct {
	documentRepository repository.DocumentRepository
}

func NewGetDocumentUseCase(documentRepository repository.DocumentRepository) GetDocumentInputPort {
	return &GetDocumentInteractor{documentRepository: documentRepository}
}

func (i *GetDocumentInteractor) Execute(ctx context.Context, input GetDocumentUseCaseInput) (*GetDocumentOutput, error) {
	// find document
	documentID, err := sharedValue.NewID(input.DocumentID)
	if err != nil {
		return nil, fmt.Errorf("failed to create document id: %w", err)
	}
	document, err := i.documentRepository.FindById(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("failed to find document: %w", err)
	}
	if document == nil {
		return nil, errors.NewUseCaseError(errors.NotFoundError, "document not found")
	}
	return &GetDocumentOutput{Document: document}, nil
}
