package service

import (
	"context"
	"fmt"

	actionValue "github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	agentState "github.com/goda6565/ai-consultant/backend/internal/domain/agent/state"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
)

type AnalyzeAction struct {
	llmClient llm.LLMClient
}

func NewAnalyzeAction(llmClient llm.LLMClient) AnalyzeActionInterface {
	return &AnalyzeAction{llmClient: llmClient}
}

func (a *AnalyzeAction) Execute(ctx context.Context, input ActionTemplateInput) (*ActionTemplateOutput, error) {
	systemPrompt := a.createSystemPrompt()
	userPrompt := a.createUserPrompt(input.State)

	llmInput := llm.GenerateTextInput{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
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

func (a *AnalyzeAction) createSystemPrompt() string {
	return analyzeSystemPrompt
}

func (a *AnalyzeAction) createUserPrompt(state agentState.State) string {
	return fmt.Sprintf(analyzeUserPrompt, state.ToPrompt())
}

var analyzeSystemPrompt = `
あなたは「問題解決エージェント」のアナリストです。
ユーザーが提供する「現在の状態」には、最終ゴールが明示されています。
この最終ゴールを常に参照しながら、収集済み情報を整理・比較・解釈してください。

出力フォーマット:
1. 要旨（150字以内）
2. 情報整理（カテゴリ別に簡潔に）
3. 観察・比較（特徴や共通点・相違点）
4. 解釈・含意（背景要因や意味合い）
5. 不足情報（追加で必要な論点やデータ）

# ポリシー
- 「現在の状態」に含まれる最終ゴールを必ず意識する
- 事実と解釈を分ける
- 外部への依頼（ヒアリング等）は提案のみ
- 冗長な説明や逐語的な思考過程は不要
`

var analyzeUserPrompt = `
以下は現在のエージェントの状態です。
この情報をもとに、収集情報の整理・比較・解釈を行い、不足情報や次に必要な調査・分析の方向性を明確化してください。

=== 現在の状態 ===
%s
`
