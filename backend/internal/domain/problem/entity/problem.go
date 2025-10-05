package entity

import (
	"time"

	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type Problem struct {
	id          sharedValue.ID
	title       value.Title
	description value.Description
	status      value.Status
	createdAt   *time.Time
}

func NewProblem(id sharedValue.ID, title value.Title, description value.Description, status value.Status, createdAt *time.Time) *Problem {
	return &Problem{id: id, title: title, description: description, status: status, createdAt: createdAt}
}

func (p *Problem) GetID() sharedValue.ID {
	return p.id
}

func (p *Problem) GetTitle() value.Title {
	return p.title
}

func (p *Problem) GetDescription() value.Description {
	return p.description
}

func (p *Problem) GetStatus() value.Status {
	return p.status
}

func (p *Problem) GetCreatedAt() *time.Time {
	return p.createdAt
}
