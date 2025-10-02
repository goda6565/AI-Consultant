package service

import (
	"context"
	"fmt"
	"strings"

	agentState "github.com/goda6565/ai-consultant/backend/internal/domain/agent/state"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
)

type ReflectionService struct {
	llmClient llm.LLMClient
}

func NewReflectionService(llmClient llm.LLMClient) *ReflectionService {
	return &ReflectionService{llmClient: llmClient}
}

type ReflectionServiceInput struct {
	State agentState.State
}

type ReflectionServiceOutput struct {
	Reflection string
}

func (r *ReflectionService) Execute(ctx context.Context, input ReflectionServiceInput) (*ReflectionServiceOutput, error) {
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
	return &ReflectionServiceOutput{Reflection: llmOutput.Text}, nil
}

func (r *ReflectionService) createSystemPrompt() string {
	return reflectionSystemPrompt
}

func (r *ReflectionService) createUserPrompt(state agentState.State) string {
	var b strings.Builder
	b.WriteString("=== 現在のエージェントの状態 ===\n")
	b.WriteString(state.ToPrompt())
	return b.String()
}

var reflectionSystemPrompt = `
あなたは「問題解決エージェント」の**内省モジュール**です。  
直近の行動履歴および状態を振り返り、次の行動選択を改善できるような自己反省を生成してください。

# 出力ルール
- 出力は日本語テキストのみ。3〜5文程度。
- 具体的な気づき、問題点、改善方針を含めること。
- 単なる感想ではなく、反省を次の行動に反映可能な洞察を含むこと。

# 内省観点（反省に含むべき要素）
1. **行動評価**：直前の action は目的に即していたか？選択ミスないか？  
2. **停滞・繰り返し**：同じパターンが続いていないか？なぜ繰り返したか？  
3. **過不足の要因**：何が足りない／過剰だったか？  
4. **代替案の検討**：次に試すべき別のアプローチや改善案  
5. **最終目標との距離**：この反省を通じて目標（レポート完成）へどう近づくか  

# 複数視点を取り入れるヒント（optional）  
- 異なる観点からの反省を生む（例：技術観点、ユーザー観点、時間効率観点）  
- 自己反省が過信や偏りに陥らないよう、複数視点を対比できるように  

# 出力例
「直前の検索は広く網羅的だったが、主要課題に絞り込みが甘く、情報の重複も多かった。特に“支店内情報共有のデジタル化”についての具体的事例が不足していた。他方、“行員スキル差”への検索は浅く終わっており、不足情報を特定できていない。次回はまず判明事項と不足事項を明確化し、それに基づくキーワードを複数案出して検索範囲を狭めてから掘り下げていく。こうすることで、次の plan／search の質を高め、最終レポート完成に近づきたい。」  
`
