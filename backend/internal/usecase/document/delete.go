package document

import (
	"context"
	"fmt"

	chunkRepository "github.com/goda6565/ai-consultant/backend/internal/domain/chunk/repository"
	documentRepository "github.com/goda6565/ai-consultant/backend/internal/domain/document/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/errors"
	storagePort "github.com/goda6565/ai-consultant/backend/internal/usecase/ports/storage"
)

type DeleteDocumentInputPort interface {
	Execute(ctx context.Context, input DeleteDocumentUseCaseInput) error
}

type DeleteDocumentUseCaseInput struct {
	DocumentID string
}

type DeleteDocumentInteractor struct {
	documentRepository documentRepository.DocumentRepository
	chunkRepository    chunkRepository.ChunkRepository
	storagePort        storagePort.StoragePort
}

func NewDeleteDocumentUseCase(documentRepository documentRepository.DocumentRepository, chunkRepository chunkRepository.ChunkRepository) DeleteDocumentInputPort {
	return &DeleteDocumentInteractor{documentRepository: documentRepository}
}

func (i *DeleteDocumentInteractor) Execute(ctx context.Context, input DeleteDocumentUseCaseInput) error {
	// find document
	documentID, err := sharedValue.NewID(input.DocumentID)
	if err != nil {
		return fmt.Errorf("failed to create document id: %w", err)
	}
	document, err := i.documentRepository.FindById(ctx, documentID)
	if err != nil {
		return fmt.Errorf("failed to find document: %w", err)
	}
	if document == nil {
		return errors.NewUseCaseError(errors.NotFoundError, "document not found")
	}

	// delete chunks
	_, err = i.chunkRepository.Delete(ctx, documentID)
	if err != nil {
		return fmt.Errorf("failed to delete chunks: %w", err)
	}

	// delete document from storage
	err = i.storagePort.Delete(ctx, document.GetStorageInfo())
	if err != nil {
		return fmt.Errorf("failed to delete document from storage: %w", err)
	}

	// delete document
	_, err = i.documentRepository.Delete(ctx, documentID)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	return nil
}
