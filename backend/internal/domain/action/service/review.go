package service

import (
	"context"
	"fmt"

	actionValue "github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	agentState "github.com/goda6565/ai-consultant/backend/internal/domain/agent/state"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
)

type ReviewAction struct {
	llmClient llm.LLMClient
}

func NewReviewAction(llmClient llm.LLMClient) ReviewActionInterface {
	return &ReviewAction{llmClient: llmClient}
}

func (r *ReviewAction) Execute(ctx context.Context, input ActionTemplateInput) (*ActionTemplateOutput, error) {
	systemPrompt := r.createSystemPrompt()
	userPrompt := r.createUserPrompt(input.State)
	llmInput := llm.GenerateTextInput{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
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

func (r *ReviewAction) createSystemPrompt() string {
	return reviewSystemPrompt
}

func (r *ReviewAction) createUserPrompt(state agentState.State) string {
	return fmt.Sprintf(reviewUserPrompt, state.ToPrompt())
}

var reviewSystemPrompt = `
あなたは「問題解決エージェント」のレビュアーです。  
最終目的は、ユーザーに提出する「提案書（Proposal）」の品質を高めることです。  
そのために、与えられた現在の内容をレビューし、改善点を指摘してください。  

必須ルール:
- 出力は日本語。箇条書き形式で簡潔にまとめる。
- 各項目は「問題点: ... / 改善提案: ...」の形式にする。
- 重要度の高いものから順に並べ、最大で5項目までに絞る。
- 指摘内容は論理の飛躍、不整合、根拠不足、曖昧表現、網羅性不足を中心に行う。
- 改善提案は実行可能な方針（例: 「信頼できる統計データを追加」など）を短文で示す。
- 自ら新しい文章や提案書本文を作成してはいけない。レビューと改善方針のみを返す。
`

var reviewUserPrompt = `
以下は現在のエージェント状態です。
この内容をレビューし、改善点と次の打ち手を提示してください。

=== 現在の状態 ===
%s
`
