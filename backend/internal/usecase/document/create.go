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
	pubsubPort "github.com/goda6565/ai-consultant/backend/internal/usecase/ports/pubsub"
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
	Document *entity.Document
}

type CreateDocumentInteractor struct {
	env                *environment.Environment
	storagePort        storagePort.StoragePort
	documentRepository repository.DocumentRepository
	duplicateChecker   *documentService.DuplicateChecker
	publisher          pubsubPort.Publisher
}

func NewCreateDocumentUseCase(env *environment.Environment, documentRepository repository.DocumentRepository, storagePort storagePort.StoragePort, duplicateChecker *documentService.DuplicateChecker, publisher pubsubPort.Publisher) CreateDocumentInputPort {
	return &CreateDocumentInteractor{
		env:                env,
		documentRepository: documentRepository,
		storagePort:        storagePort,
		duplicateChecker:   duplicateChecker,
		publisher:          publisher,
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
	rawStoragePath, err := i.createStoragePath(input.Title, input.DocumentExtension)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage path: %w", err)
	}
	storagePath := value.NewStorageInfo(bucketName, rawStoragePath)
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

	if err := i.publisher.Publish(ctx, pubsubPort.PubsubMessage{DocumentID: document.GetID().Value()}); err != nil {
		return nil, fmt.Errorf("failed to publish pubsub message: %w", err)
	}

	return &CreateDocumentOutput{Document: document}, nil
}

func (i *CreateDocumentInteractor) createStoragePath(title string, documentExtension string) (string, error) {
	switch documentExtension {
	case "pdf":
		return fmt.Sprintf("%s.pdf", title), nil
	case "markdown":
		return fmt.Sprintf("%s.md", title), nil
	case "csv":
		return fmt.Sprintf("%s.csv", title), nil
	default:
		return "", errors.NewUseCaseError(errors.InternalError, "invalid document extension")
	}
}
