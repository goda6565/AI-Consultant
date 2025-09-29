package service

import (
	"context"
	"encoding/json"
	"fmt"

	actionValue "github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	agentState "github.com/goda6565/ai-consultant/backend/internal/domain/agent/state"
	agentValue "github.com/goda6565/ai-consultant/backend/internal/domain/agent/value"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
)

type WriteAction struct {
	llmClient llm.LLMClient
}

func NewWriteAction(llmClient llm.LLMClient) WriteActionInterface {
	return &WriteAction{llmClient: llmClient}
}

type WriteActionOutputStruct struct {
	Content      string `json:"content"`
	ChangeReason string `json:"change_reason"`
}

func (w *WriteAction) Execute(ctx context.Context, input ActionTemplateInput) (*ActionTemplateOutput, error) {
	systemPrompt := w.createSystemPrompt()
	userPrompt := w.createUserPrompt(input.State)
	llmInput := llm.GenerateStructuredTextInput{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
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

func (w *WriteAction) createSystemPrompt() string {
	return writeSystemPrompt
}

func (w *WriteAction) createUserPrompt(state agentState.State) string {
	return fmt.Sprintf(writeUserPrompt, state.ToPrompt())
}

var writeSystemPrompt = `
あなたは「問題解決エージェント」のライターです。
目的は、ユーザーに提出する「提案書（Proposal）」の本文を、現時点の情報に基づいて日本語で作成することです。

必須ルール:
- 出力は必ず JSON {"content": "...", "change_reason": "..."} のみとする。
- "content" には **提案書の本文を必ず全文** Markdown 形式で生成すること。
- "change_reason" には、今回どの部分を更新／修正したか、または変更の理由を1〜2文で記述すること。
- 変更が微小であっても、"content" は全文を再生成すること。

Markdown の書式制約:
- 大見出しは "##"、小見出しは "###" を使う。
- 箇条書きは "-" を使い、重要点を簡潔にまとめる。
- 段落は1〜3文以内で区切る。
- 表を使う場合は Markdown の表記法を守る。
- 不要な空行や装飾は使わない。

内容に関するルール:
- 冗長にならず、ビジネス文書として簡潔かつ明確に書く。
- 事実関係が不確かな点は断定せず「追加調査が必要」と明示する。
- 数値・根拠・参照がある場合は本文中に簡潔に示す。
`
var writeUserPrompt = `
以下は現在のエージェント状態です。
この内容を踏まえて、新しい提案書本文を生成してください。

=== 現在の状態 ===
%s
`
