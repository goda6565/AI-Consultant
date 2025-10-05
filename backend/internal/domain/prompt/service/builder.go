package service

import (
	actionValue "github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	agentState "github.com/goda6565/ai-consultant/backend/internal/domain/agent/state"
	"github.com/goda6565/ai-consultant/backend/internal/domain/prompt/prompts"
)

type PromptBuilderInput struct {
	Name       string
	ActionType actionValue.ActionType
	State      agentState.State
	Input      string
}

type PromptBuilderOutput struct {
	SystemPrompt string
	UserPrompt   string
}

type PromptBuilder struct {
}

func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{}
}

func (b *PromptBuilder) Build(input PromptBuilderInput) *PromptBuilderOutput {
	switch input.ActionType {
	case actionValue.ActionTypePlan:
		return &PromptBuilderOutput{
			SystemPrompt: prompts.PlanSystemPrompt(input.State),
			UserPrompt:   prompts.PlanUserPrompt(input.State),
		}
	case actionValue.ActionTypeExternalSearch:
		switch input.Name {
		case "decompose":
			return &PromptBuilderOutput{
				SystemPrompt: prompts.ExternalDecomposeSystemPrompt(),
				UserPrompt:   prompts.ExternalDecomposeUserPrompt(input.Input, input.State),
			}
		case "explore":
			return &PromptBuilderOutput{
				SystemPrompt: prompts.SearchExplorePrompt(),
				UserPrompt:   prompts.SearchExploreUserPrompt(input.Input),
			}
		case "synthesize":
			return &PromptBuilderOutput{
				SystemPrompt: prompts.SearchSynthesizePrompt(),
				UserPrompt:   prompts.SearchSynthesizeUserPrompt(input.Input),
			}
		}
	case actionValue.ActionTypeInternalSearch:
		switch input.Name {
		case "decompose":
			return &PromptBuilderOutput{
				SystemPrompt: prompts.InternalDecomposeSystemPrompt(),
				UserPrompt:   prompts.InternalDecomposeUserPrompt(input.Input, input.State),
			}
		case "explore":
			return &PromptBuilderOutput{
				SystemPrompt: prompts.SearchExplorePrompt(),
				UserPrompt:   prompts.SearchExploreUserPrompt(input.Input),
			}
		case "synthesize":
			return &PromptBuilderOutput{
				SystemPrompt: prompts.SearchSynthesizePrompt(),
				UserPrompt:   prompts.SearchSynthesizeUserPrompt(input.Input),
			}
		}
	case actionValue.ActionTypeAnalyze:
		return &PromptBuilderOutput{
			SystemPrompt: prompts.AnalyzeSystemPrompt(input.State),
			UserPrompt:   prompts.AnalyzeUserPrompt(input.State),
		}
	case actionValue.ActionTypeWrite:
		return &PromptBuilderOutput{
			SystemPrompt: prompts.WriteSystemPrompt(input.State),
			UserPrompt:   prompts.WriteUserPrompt(input.State),
		}
	case actionValue.ActionTypeReview:
		return &PromptBuilderOutput{
			SystemPrompt: prompts.ReviewSystemPrompt(input.State),
			UserPrompt:   prompts.ReviewUserPrompt(input.State),
		}
	}
	return nil
}
