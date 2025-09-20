package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	hearingMessageEntity "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
	problemEntity "github.com/goda6565/ai-consultant/backend/internal/domain/problem/entity"
	problemFieldEntity "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/entity"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
)

type JudgeProblemFieldCompletionService struct {
	llmClient llm.LLMClient
}

type JudgeProblemFieldCompletionServiceInput struct {
	Problem              problemEntity.Problem
	HearingMessages      []hearingMessageEntity.HearingMessage
	TargetProblemFieldID sharedValue.ID
	ProblemFields        []problemFieldEntity.ProblemField
}

type JudgeProblemFieldCompletionServiceOutput struct {
	IsTargetProblemFieldAnswered bool
}

type judgeProblemFieldCompletionLLMOutputStruct struct {
	IsTargetProblemFieldAnswered bool `json:"is_target_problem_field_answered"`
}

func NewJudgeProblemFieldCompletionService(llmClient llm.LLMClient) *JudgeProblemFieldCompletionService {
	return &JudgeProblemFieldCompletionService{llmClient: llmClient}
}

func (s *JudgeProblemFieldCompletionService) Execute(ctx context.Context, input JudgeProblemFieldCompletionServiceInput, logger logger.Logger) (*JudgeProblemFieldCompletionServiceOutput, error) {
	userPrompt := s.createUserPrompt(input.Problem, input.HearingMessages, input.TargetProblemFieldID, input.ProblemFields)
	logger.Info("user prompt", "userPrompt", userPrompt)
	llmInput := llm.GenerateStructuredTextInput{
		SystemPrompt: s.createSystemPrompt(),
		UserPrompt:   userPrompt,
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		Schema: json.RawMessage(`
			{
				"type": "object",
				"properties": {
					"is_target_problem_field_answered": {"type": "boolean"}
				}
			}
		`),
	}

	llmOutput, err := s.llmClient.GenerateStructuredText(ctx, llmInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}

	var parsed judgeProblemFieldCompletionLLMOutputStruct
	if err := json.Unmarshal([]byte(llmOutput.Text), &parsed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal output: %w", err)
	}

	return &JudgeProblemFieldCompletionServiceOutput{IsTargetProblemFieldAnswered: parsed.IsTargetProblemFieldAnswered}, nil
}

func (s *JudgeProblemFieldCompletionService) createSystemPrompt() string {
	return judgeProblemFieldCompletionSystemPrompt
}

func (s *JudgeProblemFieldCompletionService) createUserPrompt(
	problem problemEntity.Problem,
	hearingMessages []hearingMessageEntity.HearingMessage,
	targetProblemFieldID sharedValue.ID,
	problemFields []problemFieldEntity.ProblemField,
) string {
	var targetProblemField problemFieldEntity.ProblemField
	for _, f := range problemFields {
		if f.GetID() == targetProblemFieldID {
			targetProblemField = f
			break
		}
	}
	var b strings.Builder

	// --- 課題情報 ---
	b.WriteString("【課題情報】\n")
	b.WriteString(fmt.Sprintf("タイトル: %s\n", problem.GetTitle().Value()))
	b.WriteString(fmt.Sprintf("説明: %s\n\n", problem.GetDescription().Value()))

	// --- 項目一覧 ---
	b.WriteString("【収集すべき項目一覧（インデックス: 名称 / 回答済フラグ）】\n")
	for idx, f := range problemFields {
		flag := "未回答"
		af := f.GetAnswered()
		if af.Value() {
			flag = "回答済"
		}
		fld := f.GetField()
		b.WriteString(fmt.Sprintf("%d: %s (%s)\n", idx, (&fld).Value(), flag))
	}

	// --- ターゲット項目 ---
	b.WriteString("\n【判定対象フィールドID】\n")
	fld := targetProblemField.GetField()
	b.WriteString((&fld).Value())
	b.WriteString("\n")

	// --- 会話履歴 ---
	b.WriteString("\n【これまでの会話履歴（古い→新しい）】\n")
	if len(hearingMessages) == 0 {
		b.WriteString("（まだ会話履歴はありません）\n")
	} else {
		for _, m := range hearingMessages {
			role := "user"
			if m.GetRole().Value() == "assistant" {
				role = "assistant"
			}
			msg := m.GetMessage()
			b.WriteString(fmt.Sprintf("%s: %s\n", role, (&msg).Value()))
		}
	}

	return b.String()
}

var judgeProblemFieldCompletionSystemPrompt = `
あなたは経験豊富な戦略・業務コンサルタントです。あなたの役割は、クライアントとの会話履歴を読み取り、特定の課題情報スロット（ProblemField）が十分に埋まっているかを判定することです。

## 判定基準
1. **明確さ**: このフィールドを埋めるために必要な情報が具体的に得られているか
2. **網羅性**: フィールドの趣旨を満たす最小限の情報が揃っているか（多少の不足は許容するが、追加質問が明らかに必要なら未回答とする）
3. **一貫性**: 回答内容に矛盾がなく、他のフィールドで補完されるべき情報でないか

## 出力要件
- JSON形式で is_target_problem_field_answeredを返す
- true = 追加質問せずに次のフィールドに進める
- false = まだ追加で質問が必要
`
