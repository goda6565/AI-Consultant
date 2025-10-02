package service

import (
	"context"
	"fmt"

	actionValue "github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
	"github.com/goda6565/ai-consultant/backend/internal/domain/prompt/service"
)

type ReviewAction struct {
	llmClient     llm.LLMClient
	promptBuilder *service.PromptBuilder
}

func NewReviewAction(llmClient llm.LLMClient, promptBuilder *service.PromptBuilder) ReviewActionInterface {
	return &ReviewAction{llmClient: llmClient, promptBuilder: promptBuilder}
}

func (r *ReviewAction) Execute(ctx context.Context, input ActionTemplateInput) (*ActionTemplateOutput, error) {
	prompt := r.promptBuilder.Build(service.PromptBuilderInput{
		ActionType: actionValue.ActionTypeReview,
		State:      input.State,
	})
	llmInput := llm.GenerateTextInput{
		SystemPrompt: prompt.SystemPrompt,
		UserPrompt:   prompt.UserPrompt,
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		Temperature:  0.0,
	}
	llmOutput, err := r.llmClient.GenerateText(ctx, llmInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}
	action, err := CreateAction(input.State, actionValue.ActionTypeReview, "", llmOutput.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to create action: %w", err)
	}
	return &ActionTemplateOutput{Action: *action, Content: input.State.GetContent()}, nil // review action does not change content
}
