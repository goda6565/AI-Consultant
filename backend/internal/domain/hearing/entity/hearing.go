package entity

import (
	"time"

	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type Hearing struct {
	id        sharedValue.ID
	problemID sharedValue.ID
	createdAt *time.Time
}

func NewHearing(id sharedValue.ID, problemID sharedValue.ID, createdAt *time.Time) *Hearing {
	return &Hearing{id: id, problemID: problemID, createdAt: createdAt}
}

func (h *Hearing) GetID() sharedValue.ID {
	return h.id
}

func (h *Hearing) GetProblemID() sharedValue.ID {
	return h.problemID
}

func (h *Hearing) GetCreatedAt() *time.Time {
	return h.createdAt
}
