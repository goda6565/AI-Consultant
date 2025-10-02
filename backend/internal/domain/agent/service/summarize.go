package service

import (
	"context"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
)

const GeminiSummarizeTokenThreshold = 1000000
const OpenAISummarizeTokenThreshold = 1000000 // Now Dummy

type SummarizeService struct {
	llmClient llm.LLMClient
}

func NewSummarizeService(llmClient llm.LLMClient) *SummarizeService {
	return &SummarizeService{llmClient: llmClient}
}

type SummarizeServiceInput struct {
	History   string
	LLMConfig llm.LLMConfig
}

type SummarizeServiceOutput struct {
	SummarizedHistory string
}

func (s *SummarizeService) IsSummarizeNeeded(ctx context.Context, input SummarizeServiceInput) (bool, error) {
	llmInput := llm.CountTokenInput{
		Text:   input.History,
		Config: input.LLMConfig,
	}
	llmOutput, err := s.llmClient.GetTokenCount(ctx, llmInput)
	if err != nil {
		return false, fmt.Errorf("failed to get token count: %w", err)
	}
	switch input.LLMConfig.Provider {
	case llm.VertexAI:
		return llmOutput.TokenCount > GeminiSummarizeTokenThreshold, nil
	case llm.OpenAI:
		return llmOutput.TokenCount > OpenAISummarizeTokenThreshold, nil
	default:
		return false, fmt.Errorf("invalid provider: %s", input.LLMConfig.Provider)
	}
}

func (s *SummarizeService) Summarize(ctx context.Context, input SummarizeServiceInput) (*SummarizeServiceOutput, error) {
	llmInput := llm.GenerateTextInput{
		SystemPrompt: s.createSystemPrompt(),
		UserPrompt:   s.createUserPrompt(input.History),
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		Temperature:  0.0,
	}
	llmOutput, err := s.llmClient.GenerateText(ctx, llmInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}
	return &SummarizeServiceOutput{SummarizedHistory: llmOutput.Text}, nil
}

func (s *SummarizeService) createSystemPrompt() string {
	return summarizeSystemPrompt
}

func (s *SummarizeService) createUserPrompt(history string) string {
	return fmt.Sprintf(summarizeUserPrompt, history)
}

var summarizeSystemPrompt = `
あなたは「問題解決エージェント」のための履歴要約アシスタントです。  
目的は、長大になったアクション履歴（plan, search, write, review など）を整理・圧縮し、  
トークン数を削減しつつも、課題解決に必要な一貫性・流れ・重要な意思決定を失わないようにすることです。

必須ルール:
- 出力は要約済みの履歴テキストのみ。余計な説明やJSONなど他の形式は不要。
- 各アクションの「意図」「主要な内容」「重要な結果」を必ず保持する。
- 冗長な表現や繰り返しは削除する。ただし意味が変わらないよう注意する。
- 具体的な検索クエリや提案内容など、後続判断に必要な情報は必ず残す。
- URLや参照リンクは **絶対に削除・改変しない。正確なまま保持する。**
- 時系列の流れ（plan → search → write → review → ...）と因果関係は保持する。
- 要約後の長さは元の履歴のおおよそ2/3を目安とする。
- **絶対に新しい事実や情報を作り出してはいけない（虚偽・憶測は禁止）。**
- 情報を削除する場合は「重複」や「装飾」など不要な部分に限る。重要な判断は省略しない。

最も大切なのは「短くなっても過去の思考や行動の筋が失われないこと」です。
`

var summarizeUserPrompt = `
次のアクション履歴を要約してください。  
必ず上記ルールを守り、虚偽の情報や新しい情報を加えずに、冗長さを削って簡潔にまとめてください。  
特に、履歴に含まれるリンク(URL)は絶対に削除・改変せず、そのまま残してください。

=== 履歴 ===
%s
`
