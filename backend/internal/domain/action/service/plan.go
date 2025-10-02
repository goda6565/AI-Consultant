package service

import (
	"context"
	"fmt"

	actionValue "github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	"github.com/goda6565/ai-consultant/backend/internal/domain/agent/state"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
)

type PlanAction struct {
	llmClient llm.LLMClient
}

func NewPlanAction(llmClient llm.LLMClient) PlanActionInterface {
	return &PlanAction{llmClient: llmClient}
}

func (p *PlanAction) Execute(ctx context.Context, input ActionTemplateInput) (*ActionTemplateOutput, error) {
	systemPrompt := p.createSystemPrompt()
	userPrompt := p.createUserPrompt(input.State)
	llmInput := llm.GenerateTextInput{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		Temperature:  0.0,
	}
	llmOutput, err := p.llmClient.GenerateText(ctx, llmInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}
	action, err := CreateAction(input.State, actionValue.ActionTypePlan, "", llmOutput.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to create action: %w", err)
	}
	return &ActionTemplateOutput{Action: *action, Content: input.State.GetContent()}, nil // plan action does not change content
}

func (p *PlanAction) createSystemPrompt() string {
	return planSystemPrompt
}

func (p *PlanAction) createUserPrompt(state state.State) string {
	return fmt.Sprintf(planUserPrompt, state.ToPrompt())
}

var planSystemPrompt = `
あなたは「問題解決エージェント」のプランナーです。  
目的は **現在の状態と直近の気づき（反省）を踏まえ、次に取るべき実行可能な計画** を立てることです。  
ユーザーが提供する「現在の状態」には、最終ゴールが明示されています。
この最終ゴールを常に参照しながら、次に取るべき実行可能な計画を立ててください。

# 出力ルール
- 出力は日本語で、「ステップ形式の TODO リスト」または「分岐付きリスト形式」とする。  
- 余計な説明や理由は書かず、実行可能な TODO のみを書く。  
- TODO は最大3個まで。ただし、条件分岐が必要な場合は “もし〜なら” 形式を含めてもよい。  
- 各 TODO は「目的」と「手段」がわかる簡潔な文にする。  
- 不足情報がある場合は「調査する」「確認する」を TODO に含める。  
- TODO の順序は実行順序を示す。  
- 代替案がある場合は分岐として明記してよい（例: 「もしAなら〜、もしBなら〜」）。  

# 計画の観点
1. **反省の活用**：直前までの自己反省を踏まえて改善点を反映させる  
2. **停滞回避**：同じ行動を繰り返さず、新しい進展を促す計画にする  
3. **代替案の明示**：複数の可能性があるときは分岐で示す  
4. **最終目標との整合性**：各 TODO が最終目標（レポート完成）にどう貢献するか明確に意識する  

# 出力例（few-shot）

例1:
=== 状態 ===  
情報収集が十分ではない  
TODO:
1. 調査する：主要課題に関する具体的な事例を検索する  
2. 整理する：収集した情報を課題ごとに分類する  
3. 計画する：不足部分に焦点を当てた次の検索キーワードを立案する  

例2:
=== 状態 ===  
検索済みだが複数選択肢がある  
TODO:
1. 比較する：各選択肢の長所と短所を整理する  
2. 選択する：最適案を決定する  
3. 分岐する：もし決定に不確実性が残るなら追加調査を行う  

例3:
=== 状態 ===  
プランが決まっている  
TODO:
1. 実行準備する：必要なリソースや前提条件を確認する  
2. 実行する：計画に基づき具体的な出力を進める  
`

var planUserPrompt = `
以下は現在のエージェント状態です。  
これを踏まえ、ステップ形式の TODO リストを最大3つまで提案してください（分岐付きも可）。

=== 現在の状態 ===  
%s
`
