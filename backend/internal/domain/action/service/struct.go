package service

import (
	"context"
	"fmt"

	actionValue "github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	agentState "github.com/goda6565/ai-consultant/backend/internal/domain/agent/state"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
)

type StructAction struct {
	llmClient llm.LLMClient
}

func NewStructAction(llmClient llm.LLMClient) StructActionInterface {
	return &StructAction{llmClient: llmClient}
}

func (s *StructAction) Execute(ctx context.Context, input ActionTemplateInput) (*ActionTemplateOutput, error) {
	systemPrompt := s.createSystemPrompt()
	userPrompt := s.createUserPrompt(input.State)

	llmInput := llm.GenerateTextInput{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		Temperature:  0.0,
	}

	llmOutput, err := s.llmClient.GenerateText(ctx, llmInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}

	action, err := CreateAction(input.State, actionValue.ActionTypeStruct, "", llmOutput.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to create action: %w", err)
	}

	return &ActionTemplateOutput{
		Action:  *action,
		Content: input.State.GetContent(), // struct action does not change content
	}, nil
}

func (s *StructAction) createSystemPrompt() string {
	return structSystemPrompt
}

func (s *StructAction) createUserPrompt(state agentState.State) string {
	return fmt.Sprintf(structUserPrompt, state.ToPrompt())
}

var structSystemPrompt = `
あなたはコンサルタントとして、提案書の「構成案」を作成する専門家です。

目的:
- 提案書やレポートの章立て・見出し・論点を整理し、全体の流れを明確にする。

出力ルール:
- 日本語で出力する
- 箇条書きまたは番号付きリストで章構成を示す
- 最大で 5〜7 章程度
- 簡潔なタイトル + 必要なら補足説明
- 曖昧表現は避け、論理的な流れに従う
`

var structUserPrompt = `
以下は現在のエージェントの状態です。
この情報をもとに、提案書の章立て・目次案を作成してください。

=== 現在の状態 ===
%s
`
