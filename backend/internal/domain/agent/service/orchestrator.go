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
}

type OrchestratorOutputStruct struct {
	Action string `json:"action"`
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
					}
				},
				"required": ["action"]
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
	return &OrchestratorOutput{NextAction: nextAction}, nil
}

func (o *Orchestrator) createSystemPrompt() string {
	return fmt.Sprintf(orchestratorSystemPrompt, actionValue.AvailableActionTypes())
}

func (o *Orchestrator) createUserPrompt(state agentState.State) string {
	var b strings.Builder
	b.WriteString("=== 現在のエージェントの状態 ===\n")
	b.WriteString(state.ToPrompt())
	b.WriteString("\n\n次に取るべきアクションを1つ選んでください。")
	b.WriteString("\n\n利用可能なアクション一覧:")
	b.WriteString(actionValue.AvailableActionTypes())
	return b.String()
}

var orchestratorSystemPrompt = `
あなたは「問題解決エージェント」のオーケストレーターです。  
ユーザーの課題解決に向けて、現在のエージェント状態を読み取り、次に実行すべき最も適切なアクションを1つだけ選びます。  

# 出力ルール
- 出力は必ず次の形式のJSONのみ：
{"action": "<アクション名>"}
- 理由・説明・補足は一切書かない
- JSON以外の文字は含めない

# アクション選択の基準
- 直近のエージェント状態を確認し、その状況に最も適したアクションを1つだけ選ぶこと  
- 利用可能なアクション一覧以外の値は絶対に出力しないこと  
- 同じアクションを繰り返すのではなく、状態に応じて適切に遷移させること  

例えば直近でplanを実行していたら、planで生成された指示に従うようなアクションを選ぶ。

利用可能なアクション一覧:
%s
`
