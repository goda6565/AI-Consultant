package service

import (
	"context"
	"fmt"
	"strings"

	agentState "github.com/goda6565/ai-consultant/backend/internal/domain/agent/state"
	"github.com/goda6565/ai-consultant/backend/internal/domain/agent/value"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
)

type GoalService struct {
	llmClient llm.LLMClient
}

func NewGoalService(llmClient llm.LLMClient) *GoalService {
	return &GoalService{llmClient: llmClient}
}

type GoalServiceInput struct {
	State agentState.State
}

type GoalServiceOutput struct {
	Goal value.Goal
}

func (g *GoalService) Execute(ctx context.Context, input GoalServiceInput) (*GoalServiceOutput, error) {
	systemPrompt := g.createSystemPrompt()
	userPrompt := g.createUserPrompt(input.State)
	llmInput := llm.GenerateTextInput{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		Temperature:  0.0,
	}
	llmOutput, err := g.llmClient.GenerateText(ctx, llmInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}
	goal := value.NewGoal(llmOutput.Text)
	return &GoalServiceOutput{Goal: *goal}, nil
}

func (g *GoalService) createSystemPrompt() string {
	return goalSystemPrompt
}

func (g *GoalService) createUserPrompt(state agentState.State) string {
	var b strings.Builder
	b.WriteString("=== 現在のエージェントの状態 ===\n")
	b.WriteString(state.ToGoalPrompt())
	return b.String()
}

var goalSystemPrompt = `
あなたは「問題解決エージェント」のゴール定義モジュールです。  
最終的なゴールは常に **「ユーザー課題に関するレポートを完成させること」** です。  
このレポートは、課題の背景・現状分析・改善策・期待効果を含み、経営層や現場に実用可能な形で提示されることを目的とします。

# 出力ルール
- 出力は日本語のテキストのみ
- 「絶対ゴール」と「達成基準（成功条件）」を必ず含める
- 達成基準は SMART 基準に従って具体的に示すこと

# SMART 基準とは
- **Specific（具体的）**: ゴールが明確で曖昧さがない  
- **Measurable（測定可能）**: 定量または定性データで達成度を評価できる  
- **Achievable（達成可能）**: 現実的かつ実行可能な水準である  
- **Relevant（関連性）**: ユーザー課題の解決やレポート完成に直接つながっている  
- **Time-bound（期限がある）**: 達成の時期や範囲が定義されている  

# 出力例
絶対ゴール: ユーザー課題に関するレポートを完成させる。レポートは課題の背景、現状分析、改善策、期待効果を含む。  

達成基準 (SMART):
- **Specific**: レポートには課題の背景、現状データ、改善策、期待効果が含まれている  
- **Measurable**: 改善策ごとに効果指標（例：業務効率の向上率、コスト削減額、満足度スコア）が数値で示されている  
- **Achievable**: 提案された改善策は利用可能なリソースや実行可能な範囲に基づいている  
- **Relevant**: 改善策はユーザーが提示した課題解決に直接貢献している  
- **Time-bound**: レポートは定められた期限内に完成し、短中期で評価できる指標を含んでいる  
`
