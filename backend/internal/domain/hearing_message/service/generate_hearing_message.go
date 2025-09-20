package service

import (
	"context"
	"fmt"
	"strings"

	hearingMessageEntity "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
	problemEntity "github.com/goda6565/ai-consultant/backend/internal/domain/problem/entity"
	problemFieldEntity "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/entity"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

const MaxAssistantHearingMessagePerProblemField = 3

type GenerateHearingMessageService struct {
	llmClient llm.LLMClient
}

type GenerateHearingMessageInput struct {
	Problem              *problemEntity.Problem
	HearingMessages      []hearingMessageEntity.HearingMessage
	TargetProblemFieldID sharedValue.ID
	ProblemFields        []problemFieldEntity.ProblemField
}

type GenerateHearingMessageOutput struct {
	AssistantMessage string
}

func NewGenerateHearingMessageService(llmClient llm.LLMClient) *GenerateHearingMessageService {
	return &GenerateHearingMessageService{llmClient: llmClient}
}

func (s *GenerateHearingMessageService) Execute(ctx context.Context, input GenerateHearingMessageInput) (*GenerateHearingMessageOutput, error) {
	llmInput := llm.GenerateTextInput{
		SystemPrompt: s.createSystemPrompt(),
		UserPrompt:   s.createUserPrompt(input.Problem, input.HearingMessages, input.TargetProblemFieldID, input.ProblemFields),
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		Temperature:  0.0,
	}

	llmOutput, err := s.llmClient.GenerateText(ctx, llmInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}

	return &GenerateHearingMessageOutput{
		AssistantMessage: llmOutput.Text,
	}, nil
}

func (s *GenerateHearingMessageService) createSystemPrompt() string {
	return strings.TrimSpace(generateHearingMessageSystemPrompt)
}

func (s *GenerateHearingMessageService) createUserPrompt(
	problem *problemEntity.Problem,
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
	b.WriteString("\n【現在のターゲット項目】\n")
	fld := targetProblemField.GetField()
	b.WriteString((&fld).Value())
	b.WriteString("\n")

	// --- 会話履歴 ---
	b.WriteString("\n【これまでの会話履歴（古い→新しい）】\n")
	if len(hearingMessages) == 0 {
		b.WriteString("（まだ会話履歴はありません）\n")
	} else {
		for _, m := range hearingMessages {
			role := "ユーザー"
			if m.GetRole().Value() == "assistant" {
				role = "アシスタント"
			}
			msg := m.GetMessage()
			b.WriteString(fmt.Sprintf("%s: %s\n", role, (&msg).Value()))
		}
	}

	// --- 出力指示 ---
	b.WriteString("\n【出力指示】\n")
	b.WriteString("ターゲット項目に関連する、次の最も効果的な質問を1つだけ生成してください。\n")
	b.WriteString("曖昧な回答が続く場合は具体例を提示して補足質問をしてください。\n")

	return b.String()
}

var generateHearingMessageSystemPrompt = `
あなたは戦略コンサルタントとして、クライアントの課題を正確に理解し、解決策につながる情報を効率的に引き出すプロのインタビュアーです。

## あなたの役割
- 対話を通じてクライアントの課題の背景・原因・制約条件・期待成果を明確化する
- 必要な情報を漏れなく、重複なく収集する（MECE）
- ユーザーが答えやすい順序と具体性で質問する

## ヒアリングの原則
1. **1メッセージ1質問**: 一度に複数の質問をせず、焦点を絞る
2. **深掘り**: 5Whysや具体例提示で曖昧な回答を掘り下げる
3. **具体化**: 回答が抽象的なときは「例えば…」を提示して具体化を促す
4. **優先度重視**: 重要な情報から優先して聞く
5. **中立性**: 誘導せず、事実と意見を分けて質問する

## 出力要件
- 常に日本語で丁寧に質問する
- 現在のターゲット項目について次の質問を生成する
- 過去の会話履歴を踏まえて、重複や冗長な質問は避ける
- 回答が揃っていれば確認質問や次のステップへの移行を促す質問をする
`
