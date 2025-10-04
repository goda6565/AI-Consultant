package llmasjudge

import (
	"context"
	"encoding/json"
	"fmt"

	actionEntity "github.com/goda6565/ai-consultant/backend/internal/domain/action/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
	reportEntity "github.com/goda6565/ai-consultant/backend/internal/domain/report/entity"
)

type Judge struct {
	llmClient llm.LLMClient
}

func NewJudge(llmClient llm.LLMClient) *Judge {
	return &Judge{
		llmClient: llmClient,
	}
}

// JudgmentResult represents the evaluation result from LLM
type JudgmentResult struct {
	OverallScore          int    `json:"overall_score"`          // 1-10
	ReliabilityCheck      int    `json:"reliability_check"`      // 1-10 信頼性・参考文献チェック
	LogicalConsistency    int    `json:"logical_consistency"`    // 1-10 論理的整合性
	Practicality          int    `json:"practicality"`           // 1-10 実用性・実行可能性
	Completeness          int    `json:"completeness"`           // 1-10 網羅性・完成度
	ActionAppropriateness int    `json:"action_appropriateness"` // 1-10 アクションの適切さ
	GoalAchievement       int    `json:"goal_achievement"`       // 1-10 目標達成度
	Summary               string `json:"summary"`                // 総合評価
	Strengths             string `json:"strengths"`              // 良かった点
	Weaknesses            string `json:"weaknesses"`             // 改善点
	GoalAchieved          bool   `json:"goal_achieved"`          // 目標達成したか
	Reasoning             string `json:"reasoning"`              // 判定理由
}

// Judge evaluates the actions and report using LLM with multiple rounds and aggregation
func (j *Judge) Judge(ctx context.Context, problemDescription string, actions []actionEntity.Action, report *reportEntity.Report, numEvaluations int) (*JudgmentResult, error) {
	systemPrompt := buildSystemPrompt()
	userPrompt := buildUserPrompt(problemDescription, actions, report)
	schema := buildJudgmentSchema()

	// 複数回評価を実行（出力後の最適化）
	results := make([]JudgmentResult, 0, numEvaluations)

	var temperature float32
	temperature = 0.0
	for i := 0; i < numEvaluations; i++ {
		input := llm.GenerateStructuredTextInput{
			SystemPrompt: systemPrompt,
			UserPrompt:   userPrompt,
			Temperature:  temperature,
			Schema:       schema,
			Config: llm.LLMConfig{
				Provider: llm.VertexAI,
				Model:    llm.Gemini25Flash,
			},
		}
		temperature += 0.1

		output, err := j.llmClient.GenerateStructuredText(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to generate judgment (round %d): %w", i+1, err)
		}

		var result JudgmentResult
		if err := json.Unmarshal([]byte(output.Text), &result); err != nil {
			return nil, fmt.Errorf("failed to parse judgment result (round %d): %w", i+1, err)
		}

		results = append(results, result)
	}

	// 複数の評価結果を統合
	aggregatedResult := aggregateResults(results)

	return aggregatedResult, nil
}

// aggregateResults aggregates multiple judgment results into one
func aggregateResults(results []JudgmentResult) *JudgmentResult {
	if len(results) == 0 {
		return nil
	}

	if len(results) == 1 {
		return &results[0]
	}

	// スコアの平均を計算
	var totalOverall, totalReliability, totalLogical, totalPracticality, totalCompleteness, totalAction, totalGoal int
	goalAchievedCount := 0

	for _, r := range results {
		totalOverall += r.OverallScore
		totalReliability += r.ReliabilityCheck
		totalLogical += r.LogicalConsistency
		totalPracticality += r.Practicality
		totalCompleteness += r.Completeness
		totalAction += r.ActionAppropriateness
		totalGoal += r.GoalAchievement
		if r.GoalAchieved {
			goalAchievedCount++
		}
	}

	n := len(results)

	// 最も詳細なフィードバックを選択（最長の文章を持つもの）
	bestResult := results[0]
	maxLength := len(bestResult.Summary) + len(bestResult.Strengths) + len(bestResult.Weaknesses) + len(bestResult.Reasoning)

	for _, r := range results[1:] {
		length := len(r.Summary) + len(r.Strengths) + len(r.Weaknesses) + len(r.Reasoning)
		if length > maxLength {
			bestResult = r
			maxLength = length
		}
	}

	return &JudgmentResult{
		OverallScore:          (totalOverall + n/2) / n, // 四捨五入
		ReliabilityCheck:      (totalReliability + n/2) / n,
		LogicalConsistency:    (totalLogical + n/2) / n,
		Practicality:          (totalPracticality + n/2) / n,
		Completeness:          (totalCompleteness + n/2) / n,
		ActionAppropriateness: (totalAction + n/2) / n,
		GoalAchievement:       (totalGoal + n/2) / n,
		GoalAchieved:          goalAchievedCount > n/2, // 過半数
		Summary:               bestResult.Summary,
		Strengths:             bestResult.Strengths,
		Weaknesses:            bestResult.Weaknesses,
		Reasoning:             bestResult.Reasoning,
	}
}

// buildJudgmentSchema returns JSON schema for structured output
func buildJudgmentSchema() json.RawMessage {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"overall_score": map[string]interface{}{
				"type":        "integer",
				"description": "Overall score from 1 to 10",
				"minimum":     1,
				"maximum":     10,
			},
			"reliability_check": map[string]interface{}{
				"type":        "integer",
				"description": "信頼性・参考文献チェック (1-10)",
				"minimum":     1,
				"maximum":     10,
			},
			"logical_consistency": map[string]interface{}{
				"type":        "integer",
				"description": "論理的整合性 (1-10)",
				"minimum":     1,
				"maximum":     10,
			},
			"practicality": map[string]interface{}{
				"type":        "integer",
				"description": "実用性・実行可能性 (1-10)",
				"minimum":     1,
				"maximum":     10,
			},
			"completeness": map[string]interface{}{
				"type":        "integer",
				"description": "網羅性・完成度 (1-10)",
				"minimum":     1,
				"maximum":     10,
			},
			"action_appropriateness": map[string]interface{}{
				"type":        "integer",
				"description": "アクションの適切さ (1-10)",
				"minimum":     1,
				"maximum":     10,
			},
			"goal_achievement": map[string]interface{}{
				"type":        "integer",
				"description": "目標達成度 (1-10)",
				"minimum":     1,
				"maximum":     10,
			},
			"summary": map[string]interface{}{
				"type":        "string",
				"description": "Overall evaluation summary in Japanese",
			},
			"strengths": map[string]interface{}{
				"type":        "string",
				"description": "What was done well in Japanese",
			},
			"weaknesses": map[string]interface{}{
				"type":        "string",
				"description": "What could be improved in Japanese",
			},
			"goal_achieved": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether the goal was achieved",
			},
			"reasoning": map[string]interface{}{
				"type":        "string",
				"description": "Reasoning for the judgment in Japanese",
			},
		},
		"required": []string{
			"overall_score",
			"reliability_check",
			"logical_consistency",
			"practicality",
			"completeness",
			"action_appropriateness",
			"goal_achievement",
			"summary",
			"strengths",
			"weaknesses",
			"goal_achieved",
			"reasoning",
		},
	}

	schemaBytes, _ := json.Marshal(schema)
	return schemaBytes
}
