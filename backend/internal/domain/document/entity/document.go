package entity

import (
	"github.com/goda6565/ai-consultant/backend/internal/domain/document/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type Document struct {
	id                sharedValue.ID
	title             value.Title
	documentType      value.DocumentType
	documentExtension value.DocumentExtension
	storagePath       *value.StoragePath
	status            value.DocumentStatus
	syncStep          value.SyncStep
	createdAt         sharedValue.DateTime
	updatedAt         sharedValue.DateTime
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

func (d *Document) SetUpdatedAt(updatedAt sharedValue.DateTime) {
	d.updatedAt = updatedAt
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

func (d *Document) GetDocumentExtension() value.DocumentExtension {
	return d.documentExtension
}

func (d *Document) GetStoragePath() *value.StoragePath {
	return d.storagePath
}

func (d *Document) GetStatus() value.DocumentStatus {
	return d.status
}

func (d *Document) GetSyncStep() value.SyncStep {
	return d.syncStep
}

func (d *Document) GetCreatedAt() sharedValue.DateTime {
	return d.createdAt
}

func (d *Document) GetUpdatedAt() sharedValue.DateTime {
	return d.updatedAt
}

func NewDocument(
	id sharedValue.ID,
	title value.Title,
	documentType value.DocumentType,
	documentExtension value.DocumentExtension,
	storagePath *value.StoragePath,
	createdAt sharedValue.DateTime,
) *Document {
	return &Document{
		id:                id,
		title:             title,
		documentType:      documentType,
		documentExtension: documentExtension,
		storagePath:       storagePath,
		status:            value.DocumentStatusProcessing,
		syncStep:          value.SyncStepPending,
		createdAt:         createdAt,
		updatedAt:         createdAt,
	}
}
