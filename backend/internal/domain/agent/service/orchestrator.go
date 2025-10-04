package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	actionValue "github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	agentState "github.com/goda6565/ai-consultant/backend/internal/domain/agent/state"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
)

const MaxSameActionExecutionCount = 3

type Orchestrator struct {
	llmClient llm.LLMClient
}

func NewOrchestrator(llmClient llm.LLMClient) *Orchestrator {
	return &Orchestrator{llmClient: llmClient}
}

type OrchestratorInput struct {
	State agentState.State
}

type OrchestratorOutput struct {
	CanProceed bool
	Reason     string
}

type OrchestratorOutputStruct struct {
	CanProceed bool   `json:"canProceed"`
	Reason     string `json:"reason"`
}

func (o *Orchestrator) Execute(ctx context.Context, input OrchestratorInput) (*OrchestratorOutput, error) {
	state := input.State

	if state.IsInitialAction() {
		return &OrchestratorOutput{CanProceed: false, Reason: "初回実行"}, nil
	}

	if state.GetCurrentActionCount() > MaxSameActionExecutionCount {
		return &OrchestratorOutput{CanProceed: true, Reason: "同じアクションの実行回数が上限に達し停滞しているため、一度次のアクションに進む"}, nil
	}

	llmInput := llm.GenerateStructuredTextInput{
		SystemPrompt: o.createSystemPrompt(state),
		UserPrompt:   o.createUserPrompt(state),
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		Schema: json.RawMessage(`
			{
				"type": "object",
				"properties": {
					"canProceed": {
						"type": "boolean"
					},
					"reason": {
						"type": "string"
					}
				},
				"required": ["canProceed", "reason"]
			}
		`),
		Temperature: 0.0,
	}
	llmOutput, err := o.llmClient.GenerateStructuredText(ctx, llmInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}

	var output OrchestratorOutputStruct
	if err := json.Unmarshal([]byte(llmOutput.Text), &output); err != nil {
		return nil, fmt.Errorf("failed to unmarshal output: %w", err)
	}

	return &OrchestratorOutput{
		CanProceed: output.CanProceed,
		Reason:     output.Reason,
	}, nil
}

func (o *Orchestrator) createSystemPrompt(state agentState.State) string {
	return fmt.Sprintf(orchestratorSystemPrompt, actionValue.ActionRoute(), state.ToActionHistory())
}

func (o *Orchestrator) createUserPrompt(state agentState.State) string {
	var b strings.Builder
	b.WriteString("=== 現在のエージェントの状態 ===\n")
	b.WriteString(state.ToPrompt())
	return b.String()
}

var orchestratorSystemPrompt = `
あなたは「問題解決エージェント」のオーケストレーターです。  
現在のエージェントの状態を踏まえ、課題解決のために次に取るべき最も適切な行動を1つ選びます。

# 絶対的な目的
常に「最終ゴール（例: レポート完成）」を参照し、そこに確実に近づく行動を選びます。

# 出力形式
- 出力は必ず次の形式のJSONのみ：
  {"canProceed": <boolean>, "reason": "<理由>"}
- canProceedはtrue（現在のアクション完了）またはfalse（現在のアクション継続）のいずれか。
- 理由は1〜2文で、現在の状態をどう判断し、その決定を選んだか説明する。
- JSON以外の文字列を含めてはならない。

# 決定の原則
1. **canProceed: true を選ぶ場合**  
   - 現在のアクションが十分に完了した
   - 次のアクションに進む準備ができている
   - 現在のタスクが適切に終了している
   - エージェントが次のステップに移れる状態
   - 次のアクションを実行する場合

2. **canProceed: false を選ぶ場合**  
   - 現在のアクションがまだ不十分
   - 現在のタスクが完了していない
   - 追加の作業や情報収集が必要
   - 現在のアクションを継続する必要がある
   - 現在のアクションを繰り返す場合

3. **判断基準**  
   - 現在の状態を正確に把握する
   - エージェントの能力と制約を考慮する
   - 目標達成の可能性を評価する
   - 継続の価値とコストを比較する

# 注意
- falseは情報やレポートが不十分の場合ではなく、現在のアクションを繰り返す場合である。

# 出力例（形式の参考のみ）
{"canProceed": true, "reason": "（現在のアクション完了の理由）"}
{"canProceed": false, "reason": "（現在のアクション継続の理由）"}

# アクションルート
%s

# 直近の行動履歴
%s
`
