package entity

import (
	"time"

	"github.com/goda6565/ai-consultant/backend/internal/domain/report/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type Report struct {
	id        sharedValue.ID
	problemID sharedValue.ID
	content   value.Content
	createdAt *time.Time
}

func NewReport(id sharedValue.ID, problemID sharedValue.ID, content value.Content, createdAt *time.Time) *Report {
	return &Report{id: id, problemID: problemID, content: content, createdAt: createdAt}
}

func (r *Report) GetID() sharedValue.ID {
	return r.id
}

func (r *Report) GetProblemID() sharedValue.ID {
	return r.problemID
}

func (r *Report) GetContent() value.Content {
	return r.content
}

func (r *Report) GetCreatedAt() *time.Time {
	return r.createdAt
}
