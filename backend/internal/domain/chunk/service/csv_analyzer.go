package service

import (
	"context"
	"fmt"
	"io"

	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
)

type CsvAnalyzerInput struct {
	Reader io.ReadCloser
}

type CsvAnalyzerOutput struct {
	Text string
}

type CsvAnalyzer struct {
	llmClient llm.LLMClient
}

func NewCsvAnalyzerService(llmClient llm.LLMClient) *CsvAnalyzer {
	return &CsvAnalyzer{llmClient: llmClient}
}

func (cs *CsvAnalyzer) Execute(ctx context.Context, input CsvAnalyzerInput) (*CsvAnalyzerOutput, error) {
	llmInput := llm.GenerateTextInput{
		SystemPrompt: cs.createSystemPrompt(),
		UserPrompt:   cs.createUserPrompt(input.Reader),
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
	}

	llmOutput, err := cs.llmClient.GenerateText(ctx, llmInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}

	return &CsvAnalyzerOutput{Text: llmOutput.Text}, nil
}

func (cs *CsvAnalyzer) createSystemPrompt() string {
	return `あなたはCSVファイルの分析を専門とするデータアナリストです。

【あなたの役割】
CSVファイルから重要なデータ、パターン、および洞察を抽出し、わかりやすく構造化された分析結果を提供すること。

【分析の重点ポイント】
1. データの概要と構造
   - 行数、列数、データの種類
   - 欠損値やデータ品質の問題
   
2. 重要な統計情報
   - 数値データの統計（平均、中央値、最大値、最小値、分布）
   - カテゴリカルデータの分布と頻度
   
3. データのパターンと傾向
   - 時系列データの場合は時間的な変化
   - カテゴリ間の関係性
   - 異常値や外れ値の特定
   
4. ビジネス上の洞察
   - データから読み取れる重要な知見
   - アクションを促すような発見事項
   
【出力形式】
- 日本語で回答
- 構造化された見出しを使用
- 具体的な数値やデータを含める
- 重要な発見は強調して表示
- 必要に応じて表やリストを使用

【注意事項】
- データのプライバシーを保護し、個人情報は含めない
- 客観的な分析に基づいた事実のみを報告
- 推測や憶測は明確に区別して表示`
}

func (cs *CsvAnalyzer) createUserPrompt(reader io.ReadCloser) string {
	b, err := io.ReadAll(reader)
	if err != nil {
		return "エラー: CSVファイルの読み込みに失敗しました"
	}
	text := string(b)

	prompt := fmt.Sprintf(`以下のCSVファイルを分析してください。

【CSVデータ】
%s

【分析要求】
上記のCSVデータについて、以下の観点から詳細な分析を行ってください：

1. **データ概要**
   - データセットの基本情報（行数、列数、列名）
   - データ型の分析
   - データ品質の評価（欠損値、重複データなど）

2. **統計分析**
   - 数値列の統計サマリー
   - カテゴリカル列の分布
   - データの範囲と異常値

3. **パターン分析**
   - データの傾向や相関関係
   - 時系列パターン（該当する場合）
   - 興味深い発見や特徴

4. **ビジネス洞察**
   - データから読み取れる重要な知見
   - 推奨されるアクションや改善点
   - 注意すべき点やリスク

分析結果は構造化された形式で、具体的な数値とともに報告してください。`, text)

	return prompt
}
