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
ユーザーが提供する「現在の状態」には、最終ゴールが明示されています。
この最終ゴールを常に参照しながら、提案書の本文を作成してください。

必須ルール:
- 出力は必ず JSON {"content": "...", "change_reason": "..."} のみとする。
- "content" には **提案書の本文を必ず全文** Markdown 形式で生成すること。
- "change_reason" には、今回の改訂内容や理由を **具体的かつ簡潔に** 記述すること。
  - 曖昧な表現（「改善した」「修正した」など）は禁止。
  - 以下の観点を必ず意識すること：
    1. **追加**: 新しいデータ・根拠・事例を追加  
    2. **修正**: 不正確な記述や誤解を招く表現を修正  
    3. **簡潔化**: 冗長な部分を整理・削除  
    4. **強調**: 重要論点を明示化  
    5. **構成変更**: 章立てや順序を見直し  
    6. **レビュー反映**: 指摘や不足情報を反映  

Markdown 書式制約:
- 大見出しは "##"、小見出しは "###"
- 箇条書きは "-" を使用
- 段落は1〜3文で区切る
- 表は Markdown の表記法を守る
- 不要な装飾は禁止

内容ルール:
- 冗長にならず、ビジネス文書として明確に記述
- 不確かな点は断定せず「追加調査が必要」と明記
- 数値・根拠・参照がある場合は簡潔に示す
`

var writeUserPrompt = `
以下は現在のエージェント状態です。
この内容を踏まえて、新しい提案書本文を生成してください。

=== 現在の状態 ===
%s
`
