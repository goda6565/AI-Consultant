package entity

import (
	"fmt"
	"time"

	"github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type Action struct {
	id         sharedValue.ID
	problemID  sharedValue.ID
	actionType value.ActionType
	input      value.ActionInput
	output     value.ActionOutput
	createdAt  *time.Time
}

func NewAction(id sharedValue.ID, problemID sharedValue.ID, actionType value.ActionType, input value.ActionInput, output value.ActionOutput, createdAt *time.Time) *Action {
	return &Action{id: id, problemID: problemID, actionType: actionType, input: input, output: output, createdAt: createdAt}
}

func (a *Action) GetID() sharedValue.ID {
	return a.id
}

func (a *Action) GetProblemID() sharedValue.ID {
	return a.problemID
}

func (a *Action) GetActionType() value.ActionType {
	return a.actionType
}

func (a *Action) GetInput() value.ActionInput {
	return a.input
}

func (a *Action) GetOutput() value.ActionOutput {
	return a.output
}

func (a *Action) GetCreatedAt() *time.Time {
	return a.createdAt
}

func (a *Action) SetCreatedAt(createdAt *time.Time) {
	a.createdAt = createdAt
}

func (a *Action) ToHistory() string {
	return fmt.Sprintf("%s: %s", a.actionType.Value(), a.output.Value())
}
