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
	errors "github.com/goda6565/ai-consultant/backend/internal/usecase/error"
	storagePort "github.com/goda6565/ai-consultant/backend/internal/usecase/ports/storage"
)

type CreateDocumentInputPort interface {
	Execute(ctx context.Context, input CreateDocumentUseCaseInput) (*CreateDocumentOutput, error)
}

type CreateDocumentUseCaseInput struct {
	Title             string
	DocumentExtension string
	File              io.Reader
}

type CreateDocumentOutput struct {
	DocumentID string
}

type CreateDocumentInteractor struct {
	env                *environment.Environment
	storagePort        storagePort.StoragePort
	documentRepository repository.DocumentRepository
	duplicateChecker   *documentService.DuplicateChecker
}

func NewCreateDocumentUseCase(env *environment.Environment, documentRepository repository.DocumentRepository, storagePort storagePort.StoragePort, duplicateChecker *documentService.DuplicateChecker) CreateDocumentInputPort {
	return &CreateDocumentInteractor{
		env:                env,
		documentRepository: documentRepository,
		storagePort:        storagePort,
		duplicateChecker:   duplicateChecker,
	}
}

func (i *CreateDocumentInteractor) Execute(ctx context.Context, input CreateDocumentUseCaseInput) (*CreateDocumentOutput, error) {
	title, err := value.NewTitle(input.Title)
	if err != nil {
		return nil, fmt.Errorf("failed to create title: %w", err)
	}

	// check duplicate
	isDuplicate, err := i.duplicateChecker.CheckDuplicateByTitle(ctx, title)
	if err != nil {
		return nil, fmt.Errorf("failed to check duplicate: %w", err)
	}
	if isDuplicate {
		return nil, errors.NewUseCaseError(errors.DuplicateError, "same title document already exists")
	}

	// upload file to storage
	bucketName := i.env.BucketName
	storagePath := value.NewStorageInfo(bucketName, input.Title)
	err = i.storagePort.Upload(ctx, bucketName, input.Title, input.File)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to storage: %w", err)
	}

	// create value objects
	uuid := uuid.NewUUID()
	id, err := sharedValue.NewID(uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to create id: %w", err)
	}
	documentExtension, err := value.NewDocumentExtension(input.DocumentExtension)
	if err != nil {
		return nil, fmt.Errorf("failed to create document extension: %w", err)
	}

	// create document
	document := entity.NewDocument(
		id,
		title,
		documentExtension,
		storagePath,
		value.DocumentStatusProcessing,
		value.SyncStepPending,
		nil,
		nil,
	)

	// save document
	err = i.documentRepository.Create(ctx, document)
	if err != nil {
		return nil, fmt.Errorf("failed to save document: %w", err)
	}

	//TODO: create vector by async

	return &CreateDocumentOutput{DocumentID: document.GetID().Value()}, nil
}
