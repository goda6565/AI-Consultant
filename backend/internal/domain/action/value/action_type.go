package value

import (
	"strings"

	"github.com/goda6565/ai-consultant/backend/internal/domain/errors"
)

type ActionType string

const (
	ActionTypePlan    ActionType = "plan"
	ActionTypeSearch  ActionType = "search"
	ActionTypeAnalyze ActionType = "analyze"
	ActionTypeWrite   ActionType = "write"
	ActionTypeReview  ActionType = "review"
	ActionTypeDone    ActionType = "done"

	// self action type
	SelfActionTypeOrchestrator ActionType = "orchestrator"
	SelfActionTypeReflection   ActionType = "reflection"
	SelfActionTypeSummarize    ActionType = "summarize"
)

func ActionRoute() string {
	var b strings.Builder
	b.WriteString("【アクションルート】\n")
	b.WriteString(`
┌──────────────┐
│     plan     │──▶ search ──▶ analyze ──▶ write ──▶ review
└──────────────┘                                  │
        ▲                                         │
        └─────────────────────────────────────────┘（内容が不十分な場合は再ループ）

reviewで十分な品質に達した場合のみ「完了」と判断する。
`)

	b.WriteString("\n【各アクションの役割】\n")
	b.WriteString("- plan: 課題を分析し、次に取るべき行動を計画する\n")
	b.WriteString("- search: 計画に基づいて必要な情報を収集・検索する\n")
	b.WriteString("- analyze: 収集した情報を要約・構造化・評価する\n")
	b.WriteString("- write: 分析結果をもとにレポートや提案文を作成する\n")
	b.WriteString("- review: 生成物の内容を検証し、改善点を特定する（品質が不十分な場合はplanへ戻る）\n")

	b.WriteString("\n【注意】\n")
	b.WriteString("reviewは最終アクションではありません。品質が基準を満たさない場合、planに戻って再ループします。\n")

	return b.String()
}

func AvailableActionTypesListWithoutDone() []ActionType {
	return []ActionType{
		ActionTypePlan,
		ActionTypeSearch,
		ActionTypeAnalyze,
		ActionTypeWrite,
		ActionTypeReview,
	}
}

func (a ActionType) Equals(other ActionType) bool {
	return a == other
}

func (a ActionType) Value() string {
	return string(a)
}

func (a ActionType) Proceed() ActionType {
	switch a {
	case ActionTypePlan:
		return ActionTypeSearch
	case ActionTypeSearch:
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
	case string(ActionTypeSearch):
		return ActionTypeSearch, nil
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
	case string(SelfActionTypeReflection):
		return SelfActionTypeReflection, nil
	case string(SelfActionTypeSummarize):
		return SelfActionTypeSummarize, nil
	default:
		return "", errors.NewDomainError(errors.ValidationError, "invalid action type")
	}
}
