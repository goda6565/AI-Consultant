package service

import (
	"fmt"

	actionEntity "github.com/goda6565/ai-consultant/backend/internal/domain/action/entity"
	actionValue "github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	agentState "github.com/goda6565/ai-consultant/backend/internal/domain/agent/state"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/uuid"
)

func CreateAction(state agentState.State, actionType actionValue.ActionType, input, output string) (*actionEntity.Action, error) {
	id, err := sharedValue.NewID(uuid.NewUUID())
	if err != nil {
		return nil, fmt.Errorf("failed to create id: %w", err)
	}
	inputValue, err := actionValue.NewActionInput(input)
	if err != nil {
		return nil, fmt.Errorf("failed to create action input: %w", err)
	}
	outputValue, err := actionValue.NewActionOutput(output)
	if err != nil {
		return nil, fmt.Errorf("failed to create action output: %w", err)
	}
	problem := state.GetProblem()
	action := actionEntity.NewAction(id, problem.GetID(), actionType, *inputValue, *outputValue, nil)
	return action, nil
}
