package service

import (
	"context"
	"fmt"

	actionValue "github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	"github.com/goda6565/ai-consultant/backend/internal/domain/agent/state"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
)

type PlanAction struct {
	llmClient llm.LLMClient
}

func NewPlanAction(llmClient llm.LLMClient) PlanActionInterface {
	return &PlanAction{llmClient: llmClient}
}

func (p *PlanAction) Execute(ctx context.Context, input ActionTemplateInput) (*ActionTemplateOutput, error) {
	systemPrompt := p.createSystemPrompt()
	userPrompt := p.createUserPrompt(input.State)
	llmInput := llm.GenerateTextInput{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		Temperature:  0.0,
	}
	llmOutput, err := p.llmClient.GenerateText(ctx, llmInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}
	action, err := CreateAction(input.State, actionValue.ActionTypePlan, "", llmOutput.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to create action: %w", err)
	}
	return &ActionTemplateOutput{Action: *action, Content: input.State.GetContent()}, nil // plan action does not change content
}

func (p *PlanAction) createSystemPrompt() string {
	return planSystemPrompt
}

func (p *PlanAction) createUserPrompt(state state.State) string {
	return fmt.Sprintf(planUserPrompt, state.ToPrompt())
}

var planSystemPrompt = `
あなたは「問題解決エージェント」のプランナーです。  
目的は、ユーザー課題の解決に向けて、今後の進め方を簡潔に整理し、  
次に取るべきアクションをTODOリスト形式で提案することです。

必須ルール:
- 出力は日本語のTODOリストのみ。余計な説明やJSONは禁止。
- TODOリストは最大3項目まで。必ず順序を明示すること。
- 各TODOは「目的」と「手段」がわかる簡潔な文にする。
- 曖昧な指示は禁止。後続の search / write / review が実行可能になるレベルで具体化する。
- 不足情報がある場合は「調査する」「確認する」という形でTODOに含める。
- 直近で取り組むべきことから順に記述する。
`

var planUserPrompt = `
以下は現在のエージェント状態です。  
これを踏まえ、次に行うべきTODOリストを最大3つまで提案してください。

=== 現在の状態 ===
%s
`
