package service

import (
	"context"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
)

type GenerateTitleService struct {
	llmClient llm.LLMClient
}

type GenerateTitleServiceInput struct {
	Description string
}

type GenerateTitleServiceOutput struct {
	Title string
}

func NewGenerateTitleService(llmClient llm.LLMClient) *GenerateTitleService {
	return &GenerateTitleService{llmClient: llmClient}
}

func (s *GenerateTitleService) Execute(ctx context.Context, input GenerateTitleServiceInput) (*GenerateTitleServiceOutput, error) {
	llmInput := llm.GenerateTextInput{
		SystemPrompt: s.createSystemPrompt(),
		UserPrompt:   s.createUserPrompt(input.Description),
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
	}

	llmOutput, err := s.llmClient.GenerateText(ctx, llmInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}

	return &GenerateTitleServiceOutput{Title: llmOutput.Text}, nil
}

func (s *GenerateTitleService) createSystemPrompt() string {
	return `あなたはコンサルティング課題のタイトルを生成するためのヘルパーです。

以下の条件に従って、適切なタイトルを生成してください：

1. タイトルは簡潔で分かりやすく、20文字以内にしてください
2. 課題の核心を表現し、クライアントや関係者がすぐに内容を理解できるようにしてください
3. ビジネス用語は適度に使用し、専門的すぎず、かつ曖昧すぎない表現を心がけてください
4. 疑問形や感嘆符は避け、名詞句の形で表現してください
5. 「〜について」「〜に関して」などの冗長な表現は使わないでください
6. 課題の種類（戦略、オペレーション、組織、システムなど）が分かるようなキーワードを含めてください

出力は生成されたタイトルのみを返してください。余分な説明や前置きは不要です。`
}

func (s *GenerateTitleService) createUserPrompt(description string) string {
	return fmt.Sprintf(`以下の課題内容から適切なタイトルを生成してください。

課題内容：
%s`, description)
}
