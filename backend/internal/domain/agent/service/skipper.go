package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	agentState "github.com/goda6565/ai-consultant/backend/internal/domain/agent/state"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
)

type Skipper struct {
	llmClient llm.LLMClient
}

func NewSkipper(llmClient llm.LLMClient) *Skipper {
	return &Skipper{llmClient: llmClient}
}

type SkipperInput struct {
	State agentState.State
}

type SkipperOutput struct {
	ShouldSkip bool
	Reason     string
}

type SkipperOutputStruct struct {
	ShouldSkip bool   `json:"shouldSkip"`
	Reason     string `json:"reason"`
}

func (t *Skipper) Execute(ctx context.Context, input SkipperInput) (*SkipperOutput, error) {
	state := input.State

	llmInput := llm.GenerateStructuredTextInput{
		SystemPrompt: t.createSystemPrompt(state),
		UserPrompt:   t.createUserPrompt(state),
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		Schema: json.RawMessage(`
			{
				"type": "object",
				"properties": {
					"shouldSkip": {
						"type": "boolean"
					},
					"reason": {
						"type": "string"
					}
				},
				"required": ["shouldSkip", "reason"]
			}
		`),
		Temperature: 0.0,
	}
	llmOutput, err := t.llmClient.GenerateStructuredText(ctx, llmInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}

	var output SkipperOutputStruct
	if err := json.Unmarshal([]byte(llmOutput.Text), &output); err != nil {
		return nil, fmt.Errorf("failed to unmarshal output: %w", err)
	}

	return &SkipperOutput{
		ShouldSkip: output.ShouldSkip,
		Reason:     output.Reason,
	}, nil
}

func (t *Skipper) createSystemPrompt(state agentState.State) string {
	return fmt.Sprintf(skipperSystemPrompt, state.ToActionHistory())
}

func (t *Skipper) createUserPrompt(state agentState.State) string {
	var b strings.Builder
	b.WriteString("=== 現在のエージェントの状態 ===\n")
	b.WriteString(state.ToPrompt())
	return b.String()
}

var skipperSystemPrompt = `
あなたはエージェントの「行動一貫性モニタ（Skipper）」です。  
目的は、現在のアクションが過去の計画（plan）や方針と矛盾していないかを判定し、  
無駄または不適切な行動をスキップさせることです。

# 判定ルール

1. **方針の再違反**
   - 直近のplanまたはreviewで、「externalSearchを行わない」「検索方針を切り替える」「データ収集を終了する」と明示されている場合、  
     externalSearchを再実行してはいけません。  
     この場合は shouldSkip = true とし、理由に「以前の方針に反するため」と明記します。

2. **利用不可能な手段**
   - internalSearchが無効な環境で、内部データ取得（internalSearch）を実行しようとした場合は常にスキップします。  
     「無効な手段のためスキップ」と明記してください。

3. **重複行動の抑制**
   - 同一目的で同じアクション（externalSearchなど）が直前に実行されており、  
     新しい条件・計画・入力が提示されていない場合はスキップします。  
     「同一目的での再実行のためスキップ」と明記してください。

4. **方向転換の促進**
   - 必要な情報が得られない場合や、同一手段を繰り返しても成果がないと判断できる場合は、  
     代わりに plan / write / review などのアクションを推奨します。  
     「同一手段の繰り返しを避け、方針転換を推奨」と理由に書きます。

# 出力形式
必ず次のJSON形式で出力してください。

{
  "shouldSkip": true or false,
  "reason": "スキップまたは実行の理由を簡潔に述べる"
}

# 出力例

## 方針違反の検出
{
  "shouldSkip": true,
  "reason": "直前のplanでexternalSearchの中止が指示されており、再実行は方針に反します。"
}

## 行動の妥当
{
  "shouldSkip": false,
  "reason": "今回のsearchは新しい目的（外部比較データの取得）に基づくため、一貫性があります。"
}
`
