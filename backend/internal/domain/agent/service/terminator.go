package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	agentState "github.com/goda6565/ai-consultant/backend/internal/domain/agent/state"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
)

const MaxActionLoopCount = 3

type Terminator struct {
	llmClient llm.LLMClient
}

func NewTerminator(llmClient llm.LLMClient) *Terminator {
	return &Terminator{llmClient: llmClient}
}

type TerminatorInput struct {
	State agentState.State
}

type TerminatorOutput struct {
	ShouldTerminate bool
	Reason          string
}

type TerminatorOutputStruct struct {
	ShouldTerminate bool   `json:"shouldTerminate"`
	Reason          string `json:"reason"`
}

func (t *Terminator) Execute(ctx context.Context, input TerminatorInput) (*TerminatorOutput, error) {
	state := input.State

	if state.GetActionLoopCount() >= MaxActionLoopCount {
		return &TerminatorOutput{
			ShouldTerminate: true,
			Reason:          fmt.Sprintf("%d回以上のループが発生しました。", MaxActionLoopCount),
		}, nil
	}

	llmInput := llm.GenerateStructuredTextInput{
		SystemPrompt: t.createSystemPrompt(state),
		UserPrompt:   t.createUserPrompt(state),
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		Schema: json.RawMessage(`
			{
				"type": "object",
				"properties": {
					"shouldTerminate": {
						"type": "boolean"
					},
					"reason": {
						"type": "string"
					}
				},
				"required": ["shouldTerminate", "reason"]
			}
		`),
		Temperature: 0.0,
	}
	llmOutput, err := t.llmClient.GenerateStructuredText(ctx, llmInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}

	var output TerminatorOutputStruct
	if err := json.Unmarshal([]byte(llmOutput.Text), &output); err != nil {
		return nil, fmt.Errorf("failed to unmarshal output: %w", err)
	}

	return &TerminatorOutput{
		ShouldTerminate: output.ShouldTerminate,
		Reason:          output.Reason,
	}, nil
}

func (t *Terminator) createSystemPrompt(state agentState.State) string {
	return fmt.Sprintf(terminatorSystemPrompt, state.ToActionHistory())
}

func (t *Terminator) createUserPrompt(state agentState.State) string {
	var b strings.Builder
	b.WriteString("=== 現在のエージェントの状態 ===\n")
	b.WriteString(state.ToPrompt())
	return b.String()
}

var terminatorSystemPrompt = `
あなたは「問題解決エージェント」の終了判定者です。
現在のエージェントの状態を踏まえ、作業を終了すべきかどうかを判定します。

# 絶対的な目的
最終ゴール（例: レポート完成）が達成されたかどうかを正確に判定し、終了のタイミングを決定します。

# 出力形式
- 出力は必ず次の形式のJSONのみ：
  {"shouldTerminate": <boolean>, "reason": "<理由>"}
- shouldTerminateはtrue（終了）またはfalse（継続）のいずれか。
- 理由は1〜2文で、現在の状態をどう判断し、その判定を選んだか説明する。
- JSON以外の文字列を含めてはならない。

# 判定の原則
1. **shouldTerminate: true を選ぶ場合**  
   - 最終ゴールが明確に達成された
   - 必要な成果物が完成している
   - 追加の作業が不要である
   - 目標の品質基準を満たしている

2. **shouldTerminate: false を選ぶ場合**  
   - 最終ゴールが未達成
   - 成果物が不完全
   - 追加の作業や改善が必要
   - 目標の品質基準を満たしていない

3. **判断基準**  
   - 最終ゴールの定義を正確に理解する
   - 現在の成果物の完成度を評価する
   - 品質基準や要件を確認する
   - 追加作業の必要性を判断する

# 出力例（形式の参考のみ）
{"shouldTerminate": true, "reason": "（終了理由）"}
{"shouldTerminate": false, "reason": "（継続理由）"}
`
