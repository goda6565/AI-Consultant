package entity

import (
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	value "github.com/goda6565/ai-consultant/backend/internal/domain/state/value"
)

type ProblemDefinition struct {
	id                 sharedValue.ID
	summary            value.DefinedProblemSummary
	hearingInformation value.DefinedProblemHearingInformation
	goal               value.DefinedProblemGoal
}

func NewProblemDefinition(id sharedValue.ID, summary value.DefinedProblemSummary, hearingInformation value.DefinedProblemHearingInformation, goal value.DefinedProblemGoal) *ProblemDefinition {
	return &ProblemDefinition{id: id, summary: summary, hearingInformation: hearingInformation, goal: goal}
}

func (p *ProblemDefinition) GetSummary() value.DefinedProblemSummary {
	return p.summary
}

func (p *ProblemDefinition) GetHearingInformation() value.DefinedProblemHearingInformation {
	return p.hearingInformation
}

func (p *ProblemDefinition) GetGoal() value.DefinedProblemGoal {
	return p.goal
}
