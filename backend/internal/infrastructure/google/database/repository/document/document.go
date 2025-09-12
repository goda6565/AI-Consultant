package document

import (
	"context"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/document/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/document/repository"
	"github.com/goda6565/ai-consultant/backend/internal/domain/document/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/internal/gen/app"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/helper"
	"github.com/jackc/pgx/v5/pgtype"
)

type DocumentRepository struct {
	pool *database.AppPool
}

func NewDocumentRepository(pool *database.AppPool) repository.DocumentRepository {
	return &DocumentRepository{pool: pool}
}

func (r *DocumentRepository) FindAll(ctx context.Context) ([]entity.Document, error) {
	q := app.New(r.pool)
	documents, err := q.GetAllDocuments(ctx)
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to get all documents: %v", err))
	}
	entities := make([]entity.Document, len(documents))
	for i, document := range documents {
		entity, err := toEntity(document)
		if err != nil {
			return nil, fmt.Errorf("failed to convert document to entity: %v", err)
		}
		entities[i] = *entity
	}
	return entities, nil
}

func (r *DocumentRepository) FindById(ctx context.Context, id sharedValue.ID) (*entity.Document, error) {
	q := app.New(r.pool)
	var documentID pgtype.UUID
	if err := documentID.Scan(id.Value()); err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan id: %v", err))
	}
	document, err := q.GetDocument(ctx, documentID)
	if helper.IsNoRowsError(err) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to get document: %v", err))
	}
	return toEntity(document)
}

func (r *DocumentRepository) FindByTitle(ctx context.Context, title value.Title) (*entity.Document, error) {
	q := app.New(r.pool)
	document, err := q.GetDocumentByTitle(ctx, title.Value())
	if helper.IsNoRowsError(err) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to get document by title: %v", err))
	}
	return toEntity(document)
}

func (r *DocumentRepository) Create(ctx context.Context, document *entity.Document) error {
	q := app.New(r.pool)
	var id pgtype.UUID
	if err := id.Scan(document.GetID().Value()); err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan id: %v", err))
	}
	err := q.CreateDocument(ctx, app.CreateDocumentParams{
		ID:             id,
		Title:          document.GetTitle().Value(),
		DocumentType:   document.GetDocumentType().Value(),
		BucketName:     document.GetStorageInfo().BucketName(),
		ObjectName:     document.GetStorageInfo().ObjectName(),
		DocumentStatus: document.GetStatus().Value(),
		RetryCount:     int32(document.GetRetryCount().Value()),
	})
	if err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to create document: %v", err))
	}
	return nil
}

func (r *DocumentRepository) Update(ctx context.Context, document *entity.Document) (numUpdated int64, err error) {
	q := app.New(r.pool)
	var id pgtype.UUID
	if err := id.Scan(document.GetID().Value()); err != nil {
		return 0, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan id: %v", err))
	}
	numUpdated, err = q.UpdateDocument(ctx, app.UpdateDocumentParams{
		ID:             id,
		Title:          document.GetTitle().Value(),
		DocumentType:   document.GetDocumentType().Value(),
		BucketName:     document.GetStorageInfo().BucketName(),
		ObjectName:     document.GetStorageInfo().ObjectName(),
		DocumentStatus: document.GetStatus().Value(),
		RetryCount:     int32(document.GetRetryCount().Value()),
	})
	if err != nil {
		return 0, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to update document: %v", err))
	}
	return numUpdated, nil
}

func (r *DocumentRepository) Delete(ctx context.Context, id sharedValue.ID) (numDeleted int64, err error) {
	q := app.New(r.pool)
	var documentID pgtype.UUID
	if err := documentID.Scan(id.Value()); err != nil {
		return 0, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan id: %v", err))
	}
	numDeleted, err = q.DeleteDocument(ctx, documentID)
	if err != nil {
		return 0, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to delete document: %v", err))
	}
	return numDeleted, nil
}

func toEntity(document app.Document) (*entity.Document, error) {
	id, err := sharedValue.NewID(document.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create id: %w", err)
	}
	title, err := value.NewTitle(document.Title)
	if err != nil {
		return nil, fmt.Errorf("failed to create title: %w", err)
	}
	documentType, err := value.NewDocumentType(document.DocumentType)
	if err != nil {
		return nil, fmt.Errorf("failed to create document type: %w", err)
	}
	documentStatus, err := value.NewDocumentStatus(document.DocumentStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to create document status: %w", err)
	}
	retryCount := value.NewRetryCount(int(document.RetryCount))
	storagePath := value.NewStorageInfo(document.BucketName, document.ObjectName)
	createdAt := document.CreatedAt.Time
	updatedAt := document.UpdatedAt.Time

	return entity.NewDocument(
		id,
		title,
		documentType,
		storagePath,
		documentStatus,
		retryCount,
		&createdAt,
		&updatedAt,
	), nil
}
