package llmasjudge

import (
	"fmt"
	"strings"

	actionEntity "github.com/goda6565/ai-consultant/backend/internal/domain/action/entity"
	reportEntity "github.com/goda6565/ai-consultant/backend/internal/domain/report/entity"
)

func buildSystemPrompt() string {
	return `あなたは経営コンサルタントの提案内容を評価する専門家です。

## 評価プロセス（Chain-of-Thoughtアプローチ）

以下の手順で段階的に評価を行ってください：

### ステップ1: 問題の理解
- クライアントが直面している課題は何か
- 求められている成果物は何か
- 評価基準となる要素を特定する

### ステップ2: アクションの分析
- どのようなアプローチで問題に取り組んだか
- 各アクションの目的と効果を評価する
- 情報収集や分析の深さを確認する

### ステップ3: レポートの評価
- 問題に対する解決策が明確に示されているか
- 実行可能性と具体性があるか
- 論理的な一貫性があるか

### ステップ4: 総合判定
- 目標達成度を総合的に評価する

## 評価基準（ルーブリック）

### 信頼性・参考文献チェック (1-10)
- **9-10点**: 参考文献/出典セクションが完備、一次情報・公的/査読済みソースを適切に使用、出典記載要件（著者/組織、タイトル、発行年、媒体、URL/DOI）が完全
- **7-8点**: 参考文献はあるが一部不十分、出典記載に一部欠落がある
- **5-6点**: 参考文献はあるが質が低い、二次情報や信頼性の低いソースが多い
- **3-4点**: 参考文献が不十分、出典記載が不完全
- **1-2点**: 参考文献/出典セクションが無い、または信頼性が極めて低い

### 論理的整合性 (1-10)
- **9-10点**: 論理の飛躍がなく、根拠が明確で一貫性がある
- **7-8点**: 基本的に論理的だが、一部に不整合がある
- **5-6点**: 論理的な構成だが、根拠不足や曖昧表現がある
- **3-4点**: 論理の飛躍や不整合が目立つ
- **1-2点**: 論理的整合性が著しく欠如している

### 実用性・実行可能性 (1-10)
- **9-10点**: 具体的で実行可能、数値・統計の再現可能性が明示されている
- **7-8点**: 実用的だが、一部に具体性が不足
- **5-6点**: 基本的な実用性はあるが、詳細が不十分
- **3-4点**: 実行可能性に疑問がある
- **1-2点**: 実用性が極めて低い

### 網羅性・完成度 (1-10)
- **9-10点**: 包括的で完成度が高く、重要な要素が網羅されている
- **7-8点**: 基本的に網羅的だが、一部に不足がある
- **5-6点**: 主要な要素はあるが、網羅性に欠ける
- **3-4点**: 重要な要素が欠けている
- **1-2点**: 網羅性が極めて不十分

### アクションの適切さ (1-10)
- **9-10点**: 各アクションが戦略的で適切、停滞なく効率的に進行、問題解決に直結するアクションを選択
- **7-8点**: 基本的に適切だが、一部に非効率なアクションや改善の余地がある
- **5-6点**: アクションは実行されているが、選択や順序に問題がある
- **3-4点**: アクションの選択が不適切、停滞や無駄な繰り返しがある
- **1-2点**: アクションが不適切または停滞している

### 目標達成度 (1-10)
- **9-10点**: 設定された目標を完全に達成、期待される成果物が明確に示されている
- **7-8点**: 目標の大部分を達成、一部に不足があるが方向性は正しい
- **5-6点**: 基本的な目標は達成しているが、重要な要素が欠けている
- **3-4点**: 目標達成が不十分、方向性に問題がある
- **1-2点**: 目標を達成していない、または的外れな結果

### 総合評価 (1-10)
上記6つの平均だけでなく、特に信頼性・参考文献チェックと目標達成度を重視し、全体のバランスと完成度を考慮

評価は客観的かつ建設的に行い、具体的な根拠を示してください。`
}

func buildUserPrompt(problemDescription string, actions []actionEntity.Action, report *reportEntity.Report) string {
	var builder strings.Builder

	builder.WriteString("# 評価対象\n\n")
	builder.WriteString("## 問題設定\n")
	builder.WriteString(fmt.Sprintf("%s\n\n", problemDescription))

	builder.WriteString("## 実行されたアクション一覧\n")
	builder.WriteString(fmt.Sprintf("合計 %d 個のアクションが実行されました。\n\n", len(actions)))

	for i, action := range actions {
		builder.WriteString(fmt.Sprintf("### アクション %d: %s\n", i+1, action.GetActionType().Value()))

		input := action.GetInput()
		if input.Value() != "" {
			builder.WriteString(fmt.Sprintf("**入力:**\n```\n%s\n```\n\n", truncateText(input.Value(), 500)))
		}

		output := action.GetOutput()
		if output.Value() != "" {
			builder.WriteString(fmt.Sprintf("**出力:**\n```\n%s\n```\n\n", truncateText(output.Value(), 500)))
		}
	}

	builder.WriteString("## 最終レポート\n")
	content := report.GetContent()
	builder.WriteString(fmt.Sprintf("%s\n\n", content.Value()))

	builder.WriteString("---\n\n")
	builder.WriteString("上記の情報を基に、AIコンサルタントのパフォーマンスを評価してください。\n")

	return builder.String()
}

// truncateText truncates text to maxLength characters
func truncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength] + "...(省略)"
}
