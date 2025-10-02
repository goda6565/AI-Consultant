package service

import (
	"context"
	"encoding/json"
	"fmt"

	actionValue "github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	agentValue "github.com/goda6565/ai-consultant/backend/internal/domain/agent/value"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
	"github.com/goda6565/ai-consultant/backend/internal/domain/prompt/service"
)

type WriteAction struct {
	llmClient     llm.LLMClient
	promptBuilder *service.PromptBuilder
}

func NewWriteAction(llmClient llm.LLMClient, promptBuilder *service.PromptBuilder) WriteActionInterface {
	return &WriteAction{llmClient: llmClient, promptBuilder: promptBuilder}
}

type WriteActionOutputStruct struct {
	Content      string `json:"content"`
	ChangeReason string `json:"change_reason"`
}

func (w *WriteAction) Execute(ctx context.Context, input ActionTemplateInput) (*ActionTemplateOutput, error) {
	prompt := w.promptBuilder.Build(service.PromptBuilderInput{
		ActionType: actionValue.ActionTypeWrite,
		State:      input.State,
	})
	llmInput := llm.GenerateStructuredTextInput{
		SystemPrompt: prompt.SystemPrompt,
		UserPrompt:   prompt.UserPrompt,
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		Temperature:  0.0,
		Schema: json.RawMessage(`
			{
				"type": "object",
				"properties": {
					"content": {
						"type": "string"
					},
					"change_reason": {
						"type": "string"
					}
				},
				"required": ["content", "change_reason"]
			}
		`),
	}
	llmOutput, err := w.llmClient.GenerateStructuredText(ctx, llmInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}

	var output WriteActionOutputStruct
	if err := json.Unmarshal([]byte(llmOutput.Text), &output); err != nil {
		return nil, fmt.Errorf("failed to unmarshal output: %w", err)
	}

	newContent := agentValue.NewContent(output.Content)
	action, err := CreateAction(input.State, actionValue.ActionTypeWrite, "", output.ChangeReason)
	if err != nil {
		return nil, fmt.Errorf("failed to create action: %w", err)
	}
	return &ActionTemplateOutput{Action: *action, Content: *newContent}, nil
}
