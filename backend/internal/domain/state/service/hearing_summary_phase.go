package service

import (
	"context"
	"fmt"
	"strings"

	hearingEntity "github.com/goda6565/ai-consultant/backend/internal/domain/hearing/entity"
	hearingRepository "github.com/goda6565/ai-consultant/backend/internal/domain/hearing/repository"
	hearingMessageEntity "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/entity"
	hearingMessageRepository "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/repository"
	hearingMessageValue "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/value"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
	problemEntity "github.com/goda6565/ai-consultant/backend/internal/domain/problem/entity"
	problemRepository "github.com/goda6565/ai-consultant/backend/internal/domain/problem/repository"
	problemFieldEntity "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/entity"
	problemFieldRepository "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/repository"
	"github.com/goda6565/ai-consultant/backend/internal/domain/state/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/state/value"
	logger "github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
)

type HearingSummaryPhase struct {
	llmClient                llm.LLMClient
	problemRepository        problemRepository.ProblemRepository
	problemFieldRepository   problemFieldRepository.ProblemFieldRepository
	hearingRepository        hearingRepository.HearingRepository
	hearingMessageRepository hearingMessageRepository.HearingMessageRepository
}

func NewHearingSummaryPhase(
	llmClient llm.LLMClient,
	problemRepository problemRepository.ProblemRepository,
	hearingRepository hearingRepository.HearingRepository,
	hearingMessageRepository hearingMessageRepository.HearingMessageRepository,
) PhaseTemplate {
	return &HearingSummaryPhase{
		llmClient:                llmClient,
		problemRepository:        problemRepository,
		hearingRepository:        hearingRepository,
		hearingMessageRepository: hearingMessageRepository,
	}
}

func (h *HearingSummaryPhase) Execute(ctx context.Context, state *entity.State) (*entity.State, error) {
	logger := logger.GetLogger(ctx)
	// pre-fetch
	preFetchOutput, err := h.preFetch(ctx, state.GetProblem())
	if err != nil {
		return nil, fmt.Errorf("failed to pre-fetch: %w", err)
	}

	userPrompt := h.createUserPrompt(preFetchOutput)
	systemPrompt := h.createSystemPrompt()
	logger.Info("user prompt", "userPrompt", userPrompt)
	logger.Info("system prompt", "systemPrompt", systemPrompt)

	// generate summary
	summary, err := h.llmClient.GenerateText(ctx, llm.GenerateTextInput{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		Temperature:  0.0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate summary: %w", err)
	}

	hearingSummary, err := value.NewHearingSummary(summary.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to create hearing summary state: %w", err)
	}

	state.SetHearingSummary(hearingSummary)
	return state, nil
}

func (h *HearingSummaryPhase) createSystemPrompt() string {
	return hearingSummarySystemPrompt
}

func (h *HearingSummaryPhase) createUserPrompt(preFetchOutput *preFetchOutput) string {
	var b strings.Builder
	b.WriteString(preFetchOutput.toPrompt())
	b.WriteString("\n\n")
	b.WriteString("## 指示\n")
	b.WriteString("上記の課題情報と会話履歴をもとに、カテゴリ（項目一覧）ごとに整理された事実リストを作成してください。\n")
	b.WriteString("- 各カテゴリ名を見出しとして表示する（例: ## シフト調整）\n")
	b.WriteString("- その下に得られた事実を箇条書きで列挙する\n")
	b.WriteString("- 感情表現や推測は除外し、事実だけを残す\n")
	b.WriteString("- 頻度や条件などの具体的情報は必ず残す（例: 週1〜2回）\n")
	b.WriteString("- 出力はMarkdown形式、カテゴリ順は入力と同じ\n")
	b.WriteString("- 余計な前置きや解説は不要、結果だけを出力する\n")
	return b.String()
}

type preFetchOutput struct {
	problem         problemEntity.Problem
	problemFields   []problemFieldEntity.ProblemField
	hearing         hearingEntity.Hearing
	hearingMessages []hearingMessageEntity.HearingMessage
}

func (p *preFetchOutput) toPrompt() string {
	var b strings.Builder
	b.WriteString("【課題情報】\n")
	b.WriteString(fmt.Sprintf("タイトル: %s\n", p.problem.GetTitle().Value()))
	b.WriteString(fmt.Sprintf("説明: %s\n\n", p.problem.GetDescription().Value()))

	b.WriteString("【項目一覧】\n")
	for _, f := range p.problemFields {
		fld := f.GetField()
		b.WriteString(fmt.Sprintf("%s\n", (&fld).Value()))
	}

	fieldMap := make(map[string]string)
	for _, f := range p.problemFields {
		fid := f.GetID()
		fld := f.GetField()
		fieldMap[(&fid).Value()] = (&fld).Value()
	}

	b.WriteString("【会話履歴】\n")
	for _, m := range p.hearingMessages {
		role := "ユーザー"
		if m.GetRole().Value() == string(hearingMessageValue.RoleAssistant) {
			role = "アシスタント"
		}
		fid := m.GetProblemFieldID()
		fld := fieldMap[(&fid).Value()]
		msg := m.GetMessage()
		b.WriteString(fmt.Sprintf("%s [対象項目:%s]\n  %s\n", role, fld, (&msg).Value()))
	}
	return b.String()
}

func (h *HearingSummaryPhase) preFetch(ctx context.Context, problem *problemEntity.Problem) (*preFetchOutput, error) {
	problemFields, err := h.problemFieldRepository.FindByProblemID(ctx, problem.GetID())
	if err != nil {
		return nil, fmt.Errorf("failed to find problem fields: %w", err)
	}
	hearing, err := h.hearingRepository.FindByProblemId(ctx, problem.GetID())
	if err != nil {
		return nil, fmt.Errorf("failed to find hearing: %w", err)
	}
	hearingMessages, err := h.hearingMessageRepository.FindByHearingID(ctx, hearing.GetID())
	if err != nil {
		return nil, fmt.Errorf("failed to find hearing messages: %w", err)
	}
	return &preFetchOutput{problem: *problem, problemFields: problemFields, hearing: *hearing, hearingMessages: hearingMessages}, nil
}

var hearingSummarySystemPrompt = `
あなたはプロの業務コンサルタントとして、クライアントから収集したヒアリング結果を
分析フェーズに渡す前に整理・構造化する専門家です。

目的:
- ヒアリング結果をカテゴリ別に整理し、事実ベースの情報として要約する
- 後続の課題定義や分析フェーズで活用できるよう、ノイズを除去し重要情報を残す

重要ルール:
- 情報は一切削除せず、漏れのないよう箇条書きにする
- 感情表現・雑談・推測は含めない（忠実性を担保）
- 頻度・数値・条件は必ず保持する（例: 週1〜2回, 月3回）
- 出力はMarkdownでカテゴリ見出し「##」+ 箇条書きのみ
- 解説や余計な文章は不要
`
