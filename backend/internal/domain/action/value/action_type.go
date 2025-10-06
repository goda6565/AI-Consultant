package value

import (
	"strings"

	"github.com/goda6565/ai-consultant/backend/internal/domain/errors"
)

type ActionType string

const (
	ActionTypePlan           ActionType = "plan"
	ActionTypeExternalSearch ActionType = "externalSearch"
	ActionTypeInternalSearch ActionType = "internalSearch"
	ActionTypeAnalyze        ActionType = "analyze"
	ActionTypeWrite          ActionType = "write"
	ActionTypeReview         ActionType = "review"
	ActionTypeDone           ActionType = "done"

	// self action type
	SelfActionTypeOrchestrator ActionType = "orchestrator"
	SelfActionTypeTerminator   ActionType = "terminator"
	SelfActionTypeSkipper      ActionType = "skipper"
	SelfActionTypeReflection   ActionType = "reflection"
	SelfActionTypeSummarize    ActionType = "summarize"
)

func ActionRoute(enableInternalSearch bool) string {
	var b strings.Builder
	b.WriteString("【アクションルート】\n")

	if enableInternalSearch {
		b.WriteString(`
┌──────────────┐
│     plan     │──▶ externalSearch ──▶ internalSearch ──▶ analyze ──▶ write ──▶ review
└──────────────┘                                  │
        ▲                                         │
        └─────────────────────────────────────────┘（内容が不十分な場合は再ループ）

reviewで十分な品質に達した場合のみ「完了」と判断する。
`)
	} else {
		b.WriteString(`
┌──────────────┐
│     plan     │──▶ externalSearch ──▶ analyze ──▶ write ──▶ review
└──────────────┘                                  │
        ▲                                         │
        └─────────────────────────────────────────┘（内容が不十分な場合は再ループ）

reviewで十分な品質に達した場合のみ「完了」と判断する。
`)
	}

	b.WriteString("\n【各アクションの役割と限界】\n")
	b.WriteString("- plan: 課題を分析し、次に取るべき行動を計画する\n")
	b.WriteString("    得意: 問題構造の把握・優先度設定\n")
	b.WriteString("    限界: 新しい情報の取得や数値データの生成は行わない\n\n")

	b.WriteString("- externalSearch: 業界全体や他社動向など、外部公開情報を検索・収集する\n")
	b.WriteString("    得意: 公開レポート、ニュース、統計、他社事例の抽出\n")
	b.WriteString("    限界: 自社内部データ（顧客満足度スコア、支店別指標など）は取得できない\n\n")

	if enableInternalSearch {
		b.WriteString("- internalSearch: 支店レポートや顧客アンケートなど、内部ナレッジベースから情報を検索・抽出する\n")
		b.WriteString("    得意: 社内文書・顧客アンケート・支店レポートの参照（RAGなど）\n")
		b.WriteString("    限界: 社内で未収集の定量データや、暗黙知（個人経験・感情など）は取得できない\n\n")
	}

	b.WriteString("- analyze: 外部・内部の情報を統合し、要約・構造化・比較評価を行う\n")
	b.WriteString("    得意: 情報の整理・仮説構築・優先度分析\n")
	b.WriteString("    限界: 情報源が不足している場合、新しい事実を補完することはできない\n\n")

	b.WriteString("- write: 分析結果をもとに、レポートや提案文を生成する\n")
	b.WriteString("    得意: 明確で論理的なレポート構成、改善提案の文章化\n")
	b.WriteString("    限界: 不完全な情報や曖昧な分析を自動的に補正することはできない\n\n")

	b.WriteString("- review: 生成物の品質を検証し、改善点を特定する\n")
	b.WriteString("    得意: 内容の一貫性・論理性・KPIの妥当性の評価\n")
	b.WriteString("    限界: 情報不足や根本的なデータ欠如は指摘できるが、自動補完はできない\n\n")

	b.WriteString("【注意】\n")
	b.WriteString("reviewは最終アクションではなく、品質が基準を満たさない場合はplanに戻って再実行します。\n")
	b.WriteString("このサイクルにより、情報収集→分析→生成→検証の反復によってレポート精度を高めます。\n")

	return b.String()
}

func (a ActionType) Equals(other ActionType) bool {
	return a == other
}

func (a ActionType) Value() string {
	return string(a)
}

func (a ActionType) Proceed(enableInternalSearch bool) ActionType {
	switch a {
	case ActionTypePlan:
		return ActionTypeExternalSearch
	case ActionTypeExternalSearch:
		if enableInternalSearch {
			return ActionTypeInternalSearch
		}
		return ActionTypeAnalyze
	case ActionTypeInternalSearch:
		return ActionTypeAnalyze
	case ActionTypeAnalyze:
		return ActionTypeWrite
	case ActionTypeWrite:
		return ActionTypeReview
	case ActionTypeReview:
		return ActionTypePlan
	default:
		return a
	}
}

func NewActionType(value string) (ActionType, error) {
	switch value {
	case string(ActionTypePlan):
		return ActionTypePlan, nil
	case string(ActionTypeExternalSearch):
		return ActionTypeExternalSearch, nil
	case string(ActionTypeInternalSearch):
		return ActionTypeInternalSearch, nil
	case string(ActionTypeAnalyze):
		return ActionTypeAnalyze, nil
	case string(ActionTypeWrite):
		return ActionTypeWrite, nil
	case string(ActionTypeReview):
		return ActionTypeReview, nil
	case string(ActionTypeDone):
		return ActionTypeDone, nil
	case string(SelfActionTypeOrchestrator):
		return SelfActionTypeOrchestrator, nil
	case string(SelfActionTypeTerminator):
		return SelfActionTypeTerminator, nil
	case string(SelfActionTypeSkipper):
		return SelfActionTypeSkipper, nil
	case string(SelfActionTypeReflection):
		return SelfActionTypeReflection, nil
	case string(SelfActionTypeSummarize):
		return SelfActionTypeSummarize, nil
	default:
		return "", errors.NewDomainError(errors.ValidationError, "invalid action type")
	}
}
