package value

import (
	"fmt"
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

	ActionTypePlanDescription    = "ユーザー課題に対して、次に進むための計画を立てる。必要に応じて不足情報を明示する。"
	ActionTypeSearchDescription  = "解決に必要な追加情報を外部検索やドキュメント検索で収集する。"
	ActionTypeAnalyzeDescription = "収集した情報や計画を整理・比較・解釈し、次に進むためのインサイトや不足点を明確化する。"
	ActionTypeWriteDescription   = "提案書や回答文を具体的に生成・記述する。"
	ActionTypeReviewDescription  = "生成された内容を評価し、改善点や不足点を明示する。"
	ActionTypeDoneDescription    = "提案書が完成し、追加の調査や修正が不要な状態を確定させる。"

	ActionTypePlanDecisionGuide    = "論点が曖昧、前提が揃っていない、または次に調べるべき方向性を整理したいときに選択する。"
	ActionTypeSearchDecisionGuide  = "仮説検証や根拠づけに必要な情報が不足しているときに選択する。"
	ActionTypeAnalyzeDecisionGuide = "収集した情報をそのまま使うのではなく、整理・比較・解釈して次の計画や文章化につなげたいときに選択する。"
	ActionTypeWriteDecisionGuide   = "計画や情報が揃い、具体的な文章化に進める段階で選択する。"
	ActionTypeReviewDecisionGuide  = "生成内容を検証し、改善点を抽出したいときに選択する。"
	ActionTypeDoneDecisionGuide    = "提案内容が完成し、これ以上のアクションが不要と判断されたときに選択する。"

	// self action type
	SelfActionTypeOrchestrator ActionType = "orchestrator"
	SelfActionTypeReflection   ActionType = "reflection"
	SelfActionTypeSummarize    ActionType = "summarize"
)

func AvailableActionTypes() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("- %s: %s\n  選択基準: %s\n",
		ActionTypePlan.Value(), ActionTypePlanDescription, ActionTypePlanDecisionGuide))
	b.WriteString(fmt.Sprintf("- %s: %s\n  選択基準: %s\n",
		ActionTypeSearch.Value(), ActionTypeSearchDescription, ActionTypeSearchDecisionGuide))
	b.WriteString(fmt.Sprintf("- %s: %s\n  選択基準: %s\n",
		ActionTypeAnalyze.Value(), ActionTypeAnalyzeDescription, ActionTypeAnalyzeDecisionGuide))
	b.WriteString(fmt.Sprintf("- %s: %s\n  選択基準: %s\n",
		ActionTypeWrite.Value(), ActionTypeWriteDescription, ActionTypeWriteDecisionGuide))
	b.WriteString(fmt.Sprintf("- %s: %s\n  選択基準: %s\n",
		ActionTypeReview.Value(), ActionTypeReviewDescription, ActionTypeReviewDecisionGuide))
	b.WriteString(fmt.Sprintf("- %s: %s\n  選択基準: %s\n",
		ActionTypeDone.Value(), ActionTypeDoneDescription, ActionTypeDoneDecisionGuide))
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
