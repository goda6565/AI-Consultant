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

const FinishMessage = "アクション回数が最大値を超えたため、完了しました。"

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
	NextAction actionValue.ActionType
	Reason     string
}

type OrchestratorOutputStruct struct {
	Action string `json:"action"`
	Reason string `json:"reason"`
}

func (o *Orchestrator) Execute(ctx context.Context, input OrchestratorInput) (*OrchestratorOutput, error) {

	state := input.State
	llmInput := llm.GenerateStructuredTextInput{
		SystemPrompt: o.createSystemPrompt(),
		UserPrompt:   o.createUserPrompt(state),
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		Schema: json.RawMessage(`
			{
				"type": "object",
				"properties": {
					"action": {
						"type": "string"
					},
					"reason": {
						"type": "string"
					}
				},
				"required": ["action", "reason"]
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
	nextAction, err := actionValue.NewActionType(output.Action)
	if err != nil {
		return nil, fmt.Errorf("failed to create next action: %w", err)
	}
	return &OrchestratorOutput{NextAction: nextAction, Reason: output.Reason}, nil
}

func (o *Orchestrator) createSystemPrompt() string {
	return fmt.Sprintf(orchestratorSystemPrompt, actionValue.AvailableActionTypes())
}

func (o *Orchestrator) createUserPrompt(state agentState.State) string {
	var b strings.Builder
	b.WriteString("=== 現在のエージェントの状態 ===\n")
	b.WriteString(state.ToPrompt())
	return b.String()
}

var orchestratorSystemPrompt = `
あなたは「問題解決エージェント」のオーケストレーターです。  
役割は、現在のエージェント状態や直近の行動履歴を踏まえて、課題解決のために次に実行すべき最も適切なアクションを1つ選ぶことです。  

# 絶対的な指標
ユーザーが提供する「現在の状態」には、最終ゴールが明示されています。
この最終ゴールを常に参照しながら、次に実行すべき最も適切なアクションを選んでください。

# 出力ルール
- 出力は必ず次の形式のJSONのみ：
  {"action": "<アクション名>", "reason": "<選択理由>"}
- action は利用可能なアクション一覧から必ず1つだけ選ぶ  
- reason は1〜2文で、直前のアクションと現在の状態を振り返りつつ、次の行動を選んだ理由を簡潔に説明する  
- JSON以外の文字や余計な文章は絶対に含めない  

# アクション選択の基準
1. **前進性**：レポート完成に向けて確実に進展する行動を選ぶ  
2. **停滞回避**：同じアクションを連続して繰り返さず、停滞を避ける  
3. **反省の活用**：過去の行動の不足や失敗を踏まえ、改善につながる選択をする  
4. **段階的遷移**：plan → search → analyze → write → review → done など自然な流れを意識する  
5. **代替検討**：もし search が連続した場合は、plan や review など他の行動も検討する  
6. **ゴール整合性**：選んだ行動が最終ゴール（レポート完成）に確実に寄与していることを確認する  

# 出力例
{"action": "search", "reason": "直前のplanで不明点が残ったため、情報を補う必要がある"}  
{"action": "plan", "reason": "直前のsearchで情報を得たが整理不足があるため、方針を固める必要がある"}  
{"action": "write", "reason": "これまでのplanとsearch結果が揃ったので、下書きを進められる段階にある"}  
{"action": "review", "reason": "直近のwriteで提案書が生成されたが、改善点を抽出するために見直す"}  

# 実際の利用可能アクション一覧
%s
`
