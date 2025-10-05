package entity

import (
	"time"

	"github.com/goda6565/ai-consultant/backend/internal/domain/document/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type Document struct {
	id           sharedValue.ID
	title        value.Title
	documentType value.DocumentType
	storageInfo  value.StorageInfo
	status       value.DocumentStatus
	retryCount   value.RetryCount
	createdAt    *time.Time
	updatedAt    *time.Time
}

func (d *Document) MarkAsSyncStart() {
	d.status = value.DocumentStatusProcessing
}

func (d *Document) MarkAsSyncDone() {
	d.status = value.DocumentStatusDone
}

func (d *Document) MarkAsSyncFailed() {
	d.status = value.DocumentStatusFailed
}

func (d *Document) SetUpdatedAt(updatedAt *time.Time) {
	d.updatedAt = updatedAt
}

func (d *Document) IncrementRetryCount() {
	newRetryCount := value.NewRetryCount(d.retryCount.Value() + 1)
	d.retryCount = newRetryCount
}

func (d *Document) GetID() sharedValue.ID {
	return d.id
}

func (d *Document) GetTitle() value.Title {
	return d.title
}

func (d *Document) GetDocumentType() value.DocumentType {
	return d.documentType
}

func (d *Document) GetStorageInfo() value.StorageInfo {
	return d.storageInfo
}

func (d *Document) GetStatus() value.DocumentStatus {
	return d.status
}

func (d *Document) GetRetryCount() value.RetryCount {
	return d.retryCount
}

func (d *Document) GetCreatedAt() *time.Time {
	return d.createdAt
}

func (d *Document) GetUpdatedAt() *time.Time {
	return d.updatedAt
}

func NewDocument(
	id sharedValue.ID,
	title value.Title,
	documentType value.DocumentType,
	storageInfo value.StorageInfo,
	status value.DocumentStatus,
	retryCount value.RetryCount,
	createdAt *time.Time,
	updatedAt *time.Time,
) *Document {
	return &Document{
		id:           id,
		title:        title,
		documentType: documentType,
		storageInfo:  storageInfo,
		status:       status,
		retryCount:   retryCount,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}
