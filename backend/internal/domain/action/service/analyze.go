package service

import (
	"context"
	"fmt"

	actionValue "github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
	"github.com/goda6565/ai-consultant/backend/internal/domain/prompt/service"
)

type AnalyzeAction struct {
	llmClient     llm.LLMClient
	promptBuilder *service.PromptBuilder
}

func NewAnalyzeAction(llmClient llm.LLMClient, promptBuilder *service.PromptBuilder) AnalyzeActionInterface {
	return &AnalyzeAction{llmClient: llmClient, promptBuilder: promptBuilder}
}

func (a *AnalyzeAction) Execute(ctx context.Context, input ActionTemplateInput) (*ActionTemplateOutput, error) {
	prompt := a.promptBuilder.Build(service.PromptBuilderInput{
		ActionType: actionValue.ActionTypeAnalyze,
		State:      input.State,
	})

	llmInput := llm.GenerateTextInput{
		SystemPrompt: prompt.SystemPrompt,
		UserPrompt:   prompt.UserPrompt,
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		Temperature:  0.0,
	}

	llmOutput, err := a.llmClient.GenerateText(ctx, llmInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}

	action, err := CreateAction(input.State, actionValue.ActionTypeAnalyze, "", llmOutput.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to create action: %w", err)
	}

	return &ActionTemplateOutput{
		Action:  *action,
		Content: input.State.GetContent(), // analyze も content 自体は変えずに補足情報を出すだけ
	}, nil
}
