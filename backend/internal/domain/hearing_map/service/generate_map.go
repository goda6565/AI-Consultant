package service

import (
	"context"
	"fmt"
	"strings"

	hearingMessageEntity "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
)

type GenerateHearingMapService struct {
	llmClient llm.LLMClient
}

type GenerateHearingMapServiceInput struct {
	HearingMessages []hearingMessageEntity.HearingMessage
}

type GenerateHearingMapServiceOutput struct {
	Content string
}

func NewGenerateHearingMapService(llmClient llm.LLMClient) *GenerateHearingMapService {
	return &GenerateHearingMapService{llmClient: llmClient}
}

func (s *GenerateHearingMapService) Execute(ctx context.Context, input GenerateHearingMapServiceInput) (*GenerateHearingMapServiceOutput, error) {
	llmInput := llm.GenerateTextInput{
		SystemPrompt: s.createSystemPrompt(),
		UserPrompt:   s.createUserPrompt(input),
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
	}

	llmOutput, err := s.llmClient.GenerateText(ctx, llmInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}

	return &GenerateHearingMapServiceOutput{Content: llmOutput.Text}, nil
}

func (s *GenerateHearingMapService) createSystemPrompt() string {
	return generateHearingMapSystemPrompt
}

func (s *GenerateHearingMapService) createUserPrompt(input GenerateHearingMapServiceInput) string {
	var b strings.Builder
	b.WriteString("【会話履歴】\n")
	for _, m := range input.HearingMessages {
		msg := m.GetMessage()
		b.WriteString(fmt.Sprintf("%s: %s\n", m.GetRole().Value(), (&msg).Value()))
	}
	return b.String()
}

var generateHearingMapSystemPrompt = `
あなたは企業コンサルティング支援AIの構造化エンジンです。

# 目的
ユーザーとAIのヒアリング対話内容を分析し、
課題・要因・対応策の関係をツリー構造として整理し、
「mindmap」形式で可視化可能な形に変換してください。

# 出力形式
出力は必ず **Mermaid mindmap** の形式で出力してください。
例：
mindmap
  起点ノード
    サブノード1
      サブノード1-1
    サブノード2

# 出力ルール
- mindmap のルートには、全体テーマ（課題の主題）を配置する。
- 会話の中から「現状」「課題」「要因」「改善案」などを抽出し、意味的に整理する。
- 質問と回答の関係は階層構造（親子関係）で表現する。
- 各ノードは短いフレーズ（名詞または名詞句）で記述し、1行1要素とする。
- 同義語や重複する内容は統合し、簡潔にまとめる。
- mindmap ブロック以外の文章は出力しない。

# Few-shot例

【会話例1】
AI: 現在、若手職員が独り立ちするまでにどのくらいの期間を要していますか？
User: だいたい6〜9か月です。特に窓口業務や融資事務はもう少しかかります。
AI: 理想的にはどのくらいを目標にされていますか？
User: 4か月以内を目指しています。半年で基本業務を自走できる状態が理想です。
AI: 現在のOJTではどのような指導を行っていますか？
User: 先輩職員が横について実務を教えていますが、担当者によってばらつきがあります。
AI: どのような課題を感じていますか？
User: 指導者が通常業務と並行して行うため、時間的に余裕がなく、その場しのぎになることがあります。

【出力例1】
mindmap
  若手職員の育成
    現状
      独り立ちまで6〜9か月
      OJT中心の教育体制
      指導内容にばらつき
    課題
      育成スピードの個人差
      指導者の業務負担
    目標
      4か月以内の独り立ち
      半年で自走可能
    改善案
      OJTプロセスの標準化
      指導者負担の軽減

【会話例2】
AI: 顧客対応の中で難しいと感じる場面はありますか？
User: 特に融資関連の書類や商品の説明が難しいです。
AI: なぜ難しいと感じますか？
User: 手続きが複雑で、お客様にわかりやすく伝えるのが大変です。

【出力例2】
mindmap
  顧客対応スキル
    現状
      商品説明が難しい
      手続きが複雑
    課題
      分かりやすい説明力の不足
    改善案
      実践型ロールプレイ研修
      FAQナレッジ共有化
`
