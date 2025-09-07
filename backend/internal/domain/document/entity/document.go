package entity

import (
	"time"

	"github.com/goda6565/ai-consultant/backend/internal/domain/document/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type Document struct {
	id                sharedValue.ID
	title             value.Title
	documentExtension value.DocumentExtension
	storagePath       value.StorageInfo
	status            value.DocumentStatus
	syncStep          value.SyncStep
	createdAt         *time.Time
	updatedAt         *time.Time
}

func (d *Document) MarkAsSyncVectorDone() {
	d.syncStep = value.SyncStepDone
}

func (d *Document) MarkAsSyncDone() {
	d.status = value.DocumentStatusDone
}

func (d *Document) MarkAsFailed() {
	d.status = value.DocumentStatusFailed
}

func (d *Document) SetUpdatedAt(updatedAt *time.Time) {
	d.updatedAt = updatedAt
}

func (d *Document) GetID() sharedValue.ID {
	return d.id
}

func (d *Document) GetTitle() value.Title {
	return d.title
}

func (d *Document) GetDocumentExtension() value.DocumentExtension {
	return d.documentExtension
}

func (d *Document) GetStoragePath() value.StorageInfo {
	return d.storagePath
}

func (d *Document) GetStatus() value.DocumentStatus {
	return d.status
}

func (d *Document) GetSyncStep() value.SyncStep {
	return d.syncStep
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
	documentExtension value.DocumentExtension,
	storageInfo value.StorageInfo,
	status value.DocumentStatus,
	syncStep value.SyncStep,
	createdAt *time.Time,
	updatedAt *time.Time,
) *Document {
	return &Document{
		id:                id,
		title:             title,
		documentExtension: documentExtension,
		storagePath:       storageInfo,
		status:            status,
		syncStep:          syncStep,
		createdAt:         createdAt,
		updatedAt:         updatedAt,
	}
}
