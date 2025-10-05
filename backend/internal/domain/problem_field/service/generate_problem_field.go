package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
	problemEntity "github.com/goda6565/ai-consultant/backend/internal/domain/problem/entity"
)

const MaxProblemField = 5

type GenerateProblemFieldService struct {
	llmClient llm.LLMClient
}

type GenerateProblemFieldServiceInput struct {
	Problem *problemEntity.Problem
}

type GenerateProblemFieldServiceOutput struct {
	Fields []string
}

func NewGenerateProblemFieldService(llmClient llm.LLMClient) *GenerateProblemFieldService {
	return &GenerateProblemFieldService{llmClient: llmClient}
}

type GenerateProblemFieldServiceOutputStruct struct {
	Fields []string `json:"fields"`
}

func (s *GenerateProblemFieldService) Execute(ctx context.Context, input GenerateProblemFieldServiceInput) (*GenerateProblemFieldServiceOutput, error) {
	llmInput := llm.GenerateStructuredTextInput{
		SystemPrompt: s.createSystemPrompt(),
		UserPrompt:   s.createUserPrompt(input.Problem),
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		Schema: json.RawMessage(`
			{
				"type": "object",
				"properties": {
					"fields": {
						"type": "array",
						"items": {
							"type": "string"
						}
					}
				}
			}
		`),
	}

	llmOutput, err := s.llmClient.GenerateStructuredText(ctx, llmInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}

	var output GenerateProblemFieldServiceOutputStruct
	if err := json.Unmarshal([]byte(llmOutput.Text), &output); err != nil {
		return nil, fmt.Errorf("failed to unmarshal output: %w", err)
	}

	return &GenerateProblemFieldServiceOutput{Fields: output.Fields}, nil
}

func (s *GenerateProblemFieldService) createSystemPrompt() string {
	return fmt.Sprintf(generateProblemFieldSystemPrompt, MaxProblemField)
}

func (s *GenerateProblemFieldService) createUserPrompt(problem *problemEntity.Problem) string {
	return fmt.Sprintf(userPrompt, problem.GetTitle().Value(), problem.GetDescription().Value(), MaxProblemField)
}

var generateProblemFieldSystemPrompt = `
あなたは経験豊富な戦略・業務コンサルタントです。クライアントが抱える課題を深く理解し、実効性の高い解決策を導くための「課題情報スロット（ProblemField）」を設計するプロフェッショナルです。ここで生成するのは「質問文」ではなく、課題に付随する情報構造として事前に定義する「フィールド名」です。

## あなたの役割
- クライアントの本質的な課題を明確化する
- 問題の根本原因を特定するために必要な情報の型（ProblemField）を設計する
- 後続のヒアリングや分析が効率的に進むように情報スロットを整理する

## ProblemField 設計原則
1. **核心に直結**: 根本原因特定と打ち手設計に必要な情報フィールドを優先
2. **具体かつ簡潔**: 名詞句・短文で明瞭に（例: 「主要KPI（CVR/CTRなど）」）
3. **優先順位**: 重要度の高い順に並べる
4. **再利用可能**: 複数回ヒアリングや後続分析に使い回せる
5. **効率性**: 最大%d件以内で必要情報を網羅し、重複を避ける

## ProblemField カバレッジ例
- 現状把握: 主要KPI、ファネル指標、チャネル別内訳
- 原因分析: 仮説一覧、変更履歴、外部要因
- 影響範囲: 影響を受ける顧客/部署、業務プロセス
- 制約条件: 予算、リソース、期間、規制
- 期待成果: 成功指標、KGI/KPI、目標期日

クライアントとの限られた時間を最大限活用し、課題解決への道筋を明確にするためのProblemFieldを生成してください。

## Few-shot 例（入力と出力のイメージ）

### 例1
入力（課題情報）:
"""
課題タイトル: ECサイトのCVR低下
課題詳細: 直近3ヶ月でCVRが2.8%%から1.6%%に低下。トラフィックは横ばい。モバイル比率が増加し、カート離脱率が上昇。広告構成とLPは4月に刷新済み。
"""
想定出力(JSON):
{
  "fields": [
    "デバイス別・流入別KPI推移（6ヶ月）",
    "カート離脱ステップ内訳（配送/決済/会員登録）",
    "モバイルのコアウェブバイタル（LCP/CLS/TTFB）",
    "LP刷新後のABテスト実施一覧（仮説/結果/母数）",
    "顧客セグメント別CVR（新規/リピート/会員）",
    "購入阻害要因（送料/在庫/リードタイム）の変化",
    "決済手段カバレッジとニーズ（例: 後払い/QR）"
  ]
}
`

var userPrompt = `
## 課題情報
**課題タイトル**: %s
**課題詳細**: %s

## 指示
上記の課題について、プロのコンサルタントとして「課題に必要な情報スロット（ProblemField）」を重要度順に%d件以内で設計してください。

## 出力要件
- 各項目は名詞句/短文（質問文にしない）
- 重要度の高い順番で配列する
- 後続のヒアリングや分析で直接埋められる情報にする
- 重複や冗長表現を避け、簡潔に表現する

ProblemField名を配列形式で出力してください。
`
