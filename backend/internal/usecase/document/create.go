package document

import (
	"context"
	"fmt"
	"io"

	"github.com/goda6565/ai-consultant/backend/internal/domain/document/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/document/repository"
	documentService "github.com/goda6565/ai-consultant/backend/internal/domain/document/service"
	"github.com/goda6565/ai-consultant/backend/internal/domain/document/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/uuid"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/errors"
	syncQueuePort "github.com/goda6565/ai-consultant/backend/internal/usecase/ports/queue"
	storagePort "github.com/goda6565/ai-consultant/backend/internal/usecase/ports/storage"
)

type CreateDocumentInputPort interface {
	Execute(ctx context.Context, input CreateDocumentUseCaseInput) (*CreateDocumentOutput, error)
}

type CreateDocumentUseCaseInput struct {
	Title        string
	DocumentType string
	File         io.Reader
}

type CreateDocumentOutput struct {
	Document *entity.Document
}

type CreateDocumentInteractor struct {
	env                *environment.Environment
	storagePort        storagePort.StoragePort
	documentRepository repository.DocumentRepository
	duplicateChecker   *documentService.DuplicateChecker
	syncQueue          syncQueuePort.SyncQueue
}

func NewCreateDocumentUseCase(env *environment.Environment, documentRepository repository.DocumentRepository, storagePort storagePort.StoragePort, duplicateChecker *documentService.DuplicateChecker, syncQueue syncQueuePort.SyncQueue) CreateDocumentInputPort {
	return &CreateDocumentInteractor{
		env:                env,
		documentRepository: documentRepository,
		storagePort:        storagePort,
		duplicateChecker:   duplicateChecker,
		syncQueue:          syncQueue,
	}
}

func (i *CreateDocumentInteractor) Execute(ctx context.Context, input CreateDocumentUseCaseInput) (*CreateDocumentOutput, error) {
	title, err := value.NewTitle(input.Title)
	if err != nil {
		return nil, fmt.Errorf("failed to create title: %w", err)
	}
	documentType, err := value.NewDocumentType(input.DocumentType)
	if err != nil {
		return nil, fmt.Errorf("failed to create document type: %w", err)
	}

	// check duplicate
	isDuplicate, err := i.duplicateChecker.Execute(ctx, title)
	if err != nil {
		return nil, fmt.Errorf("failed to check duplicate: %w", err)
	}
	if isDuplicate {
		return nil, errors.NewUseCaseError(errors.DuplicateError, "same title document already exists")
	}

	// upload file to storage
	bucketName := i.env.BucketName
	objectName := fmt.Sprintf("%s.%s", input.Title, documentType.GetExtension())
	storagePath := value.NewStorageInfo(bucketName, objectName)
	err = i.storagePort.Upload(ctx, storagePath, input.File)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to storage: %w", err)
	}

	// create value objects
	uuid := uuid.NewUUID()
	id, err := sharedValue.NewID(uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to create id: %w", err)
	}

	// create document
	document := entity.NewDocument(
		id,
		title,
		documentType,
		storagePath,
		value.DocumentStatusPending, // initial document status is pending
		value.NewRetryCount(0),      // initial retry count is 0
		nil,
		nil,
	)

	// save document
	err = i.documentRepository.Create(ctx, document)
	if err != nil {
		// delete document from storage
		err = i.storagePort.Delete(ctx, storagePath)
		if err != nil {
			return nil, fmt.Errorf("failed to delete document from storage after failed to save document: %w", err)
		}
		return nil, fmt.Errorf("failed to save document: %w", err)
	}

	// publish sync queue message to vector service
	if err := i.syncQueue.Enqueue(ctx, syncQueuePort.SyncQueueMessage{DocumentID: document.GetID().Value()}); err != nil {
		return nil, fmt.Errorf("failed to publish sync queue message to vector service: %w", err)
	}

	return &CreateDocumentOutput{Document: document}, nil
}
