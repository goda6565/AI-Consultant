package entity

import (
	"time"

	"github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type ProblemField struct {
	id        sharedValue.ID
	problemID sharedValue.ID
	field     value.Field
	answered  value.Answered
	createdAt *time.Time
}

func NewProblemField(id sharedValue.ID, problemID sharedValue.ID, field value.Field, answered value.Answered, createdAt *time.Time) *ProblemField {
	return &ProblemField{id: id, problemID: problemID, field: field, answered: answered, createdAt: createdAt}
}

func (h *ProblemField) GetID() sharedValue.ID {
	return h.id
}

func (h *ProblemField) GetProblemID() sharedValue.ID {
	return h.problemID
}

func (h *ProblemField) GetField() value.Field {
	return h.field
}

func (h *ProblemField) GetAnswered() value.Answered {
	return h.answered
}

func (h *ProblemField) GetCreatedAt() *time.Time {
	return h.createdAt
}
