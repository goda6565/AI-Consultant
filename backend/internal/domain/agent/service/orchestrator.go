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
	return fmt.Sprintf(orchestratorSystemPrompt, actionValue.ActionRoute(state.GetEnableInternalSearch()), state.ToActionHistory())
}

func (o *Orchestrator) createUserPrompt(state agentState.State) string {
	var b strings.Builder
	b.WriteString("=== 現在のエージェントの状態 ===\n")
	b.WriteString(state.ToPrompt())
	return b.String()
}

var orchestratorSystemPrompt = `
あなたは「問題解決エージェント」のオーケストレーターです。  
現在のエージェントの状態と行動履歴を踏まえ、課題解決のために次に取るべき最も適切な行動を1つ選びます。

# 目的
現在のアクションが完了しているかを判断し、必要に応じて次に進むか、戦略を切り替えるかを決定します。  
常に「最終ゴール（例: レポート完成）」への最短経路を選び、非効率なループを避けることを優先します。

# 出力形式
- 出力は必ず次の形式のJSONのみ：
  {"canProceed": <boolean>, "reason": "<理由>"}
- canProceed は true（現在のアクションが完了）または false（継続が必要）のいずれか。
- 理由は1〜2文で、状態をどう判断し、なぜその決定を選んだかを簡潔に説明する。
- JSON以外の文字列を含めてはならない。

# 判断原則
1. **canProceed: true（アクション完了）**
   - 現在のアクションの目的を果たした（成功・失敗を問わない）
   - 次に進む準備が整っている
   - 例: reviewで「品質不足」と判断した場合でも、reviewの目的（評価）は完了している

2. **canProceed: false（アクション継続）**
   - 現在のアクションの目的がまだ達成されていない
   - 同じアクションを繰り返すことで明確な改善が見込める場合のみ選ぶ
   - 例: externalSearchで結果が得られず、別の検索観点を試す余地がある場合

3. **戦略転換の判断**
   - 現在のアクションで必要な情報が「そもそも到達不可能」である場合（例：内部データを外部検索で探そうとしている）、
     canProceed を true とし、次の適切なアクション（例：internalSearchまたはplan）に切り替える。
   - 「ツールや情報源の限界により得られない情報」は、繰り返しではなく方針転換で対応する。
   - 同様のアクションを3回以上繰り返すようなループは避ける。

4. **判断基準**
   - 現在の状態を正確に把握する
   - エージェントの行動履歴と情報源の制約を考慮する
   - 目標達成への貢献度と非効率リスクを比較する

# 注意
- false は「同一アクションを継続する」場合のみ選ぶ。
- 「進まない」のではなく「繰り返す価値がある」ときにのみ false。
- 「外部では得られない情報を探している」場合は、継続ではなく戦略転換（true）とする。

# 出力例
{"canProceed": true, "reason": "外部検索では内部データに到達できないため、この段階を完了とし、次に内部情報探索へ移行します。"}
{"canProceed": false, "reason": "外部検索結果が不十分なため、追加の観点で検索を続行します。"}

# 参考情報
- アクションルート:
%s
- 直近の行動履歴:
%s
`
