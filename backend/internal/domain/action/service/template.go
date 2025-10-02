package service

import (
	"context"
	"fmt"

	actionEntity "github.com/goda6565/ai-consultant/backend/internal/domain/action/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	agentState "github.com/goda6565/ai-consultant/backend/internal/domain/agent/state"
	agentValue "github.com/goda6565/ai-consultant/backend/internal/domain/agent/value"
	"github.com/goda6565/ai-consultant/backend/internal/domain/errors"
)

type ActionTemplateInput struct {
	State agentState.State
}

type ActionTemplateOutput struct {
	Content agentValue.Content
	Action  actionEntity.Action
}

type ActionTemplate interface {
	Execute(ctx context.Context, input ActionTemplateInput) (*ActionTemplateOutput, error)
}

type PlanActionInterface ActionTemplate
type SearchActionInterface ActionTemplate
type AnalyzeActionInterface ActionTemplate
type WriteActionInterface ActionTemplate
type ReviewActionInterface ActionTemplate

type ActionFactory struct {
	planActionTemplate    PlanActionInterface
	searchActionTemplate  SearchActionInterface
	analyzeActionTemplate AnalyzeActionInterface
	writeActionTemplate   WriteActionInterface
	reviewActionTemplate  ReviewActionInterface
}

func NewActionFactory(
	planActionTemplate PlanActionInterface,
	searchActionTemplate SearchActionInterface,
	analyzeActionTemplate AnalyzeActionInterface,
	writeActionTemplate WriteActionInterface,
	reviewActionTemplate ReviewActionInterface,
) *ActionFactory {
	return &ActionFactory{planActionTemplate: planActionTemplate, searchActionTemplate: searchActionTemplate, analyzeActionTemplate: analyzeActionTemplate, writeActionTemplate: writeActionTemplate, reviewActionTemplate: reviewActionTemplate}
}

func (f *ActionFactory) GetActionTemplate(actionType value.ActionType) (ActionTemplate, error) {
	switch actionType {
	case value.ActionTypePlan:
		return f.planActionTemplate, nil
	case value.ActionTypeSearch:
		return f.searchActionTemplate, nil
	case value.ActionTypeAnalyze:
		return f.analyzeActionTemplate, nil
	case value.ActionTypeWrite:
		return f.writeActionTemplate, nil
	case value.ActionTypeReview:
		return f.reviewActionTemplate, nil
	default:
		return nil, errors.NewDomainError(errors.InvalidActionType, fmt.Sprintf("invalid action type: %s", actionType))
	}
}
