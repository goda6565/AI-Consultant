package value

import (
	"fmt"
	"strings"

	"github.com/goda6565/ai-consultant/backend/internal/domain/errors"
)

type ActionType string

const (
	ActionTypePlan   ActionType = "plan"
	ActionTypeSearch ActionType = "search"
	ActionTypeStruct ActionType = "struct"
	ActionTypeWrite  ActionType = "write"
	ActionTypeReview ActionType = "review"
	ActionTypeDone   ActionType = "done"

	ActionTypePlanDescription   = "ユーザーの課題に対して、解決のための計画を立てる。"
	ActionTypeSearchDescription = "解決に必要な追加情報を外部検索やドキュメント検索で集める。"
	ActionTypeStructDescription = "提案書やレポートの構成（章立て・目次・論点整理）を設計する。"
	ActionTypeWriteDescription  = "提案書や回答文を具体的に生成・記述する。"
	ActionTypeReviewDescription = "生成された内容をレビューし、改善点を検討する。"
	ActionTypeDoneDescription   = "提案書が完成し、追加の調査や修正が不要な状態。最終成果物を確定させる。"

	ActionTypePlanDecisionGuide   = "論点が曖昧、前提が揃っていない、次に何を調べるべきか方向性が必要なときに選択する。"
	ActionTypeSearchDecisionGuide = "仮説検証や根拠づけに必要な情報が不足しているときに選択する。"
	ActionTypeStructDecisionGuide = "書く準備ができたが、提案書の流れや章立てを整理してから執筆したいときに選択する。"
	ActionTypeWriteDecisionGuide  = "計画や情報が揃い、具体的な提案や文章化に進める段階で選択する。"
	ActionTypeReviewDecisionGuide = "生成された内容に論理の飛躍・不整合・不足があるかを確認し、改善点を抽出したいときに選択する。"
	ActionTypeDoneDecisionGuide   = "提案内容が完成し、これ以上のアクションが不要だと判断されたときに選択する。"
)

func AvailableActionTypes() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("- %s: %s\n  選択基準: %s\n",
		ActionTypePlan.Value(), ActionTypePlanDescription, ActionTypePlanDecisionGuide))
	b.WriteString(fmt.Sprintf("- %s: %s\n  選択基準: %s\n",
		ActionTypeSearch.Value(), ActionTypeSearchDescription, ActionTypeSearchDecisionGuide))
	b.WriteString(fmt.Sprintf("- %s: %s\n  選択基準: %s\n",
		ActionTypeSearch.Value(), ActionTypeStructDescription, ActionTypeStructDecisionGuide))
	b.WriteString(fmt.Sprintf("- %s: %s\n  選択基準: %s\n",
		ActionTypeWrite.Value(), ActionTypeWriteDescription, ActionTypeWriteDecisionGuide))
	b.WriteString(fmt.Sprintf("- %s: %s\n  選択基準: %s\n",
		ActionTypeReview.Value(), ActionTypeReviewDescription, ActionTypeReviewDecisionGuide))
	b.WriteString(fmt.Sprintf("- %s: %s\n  選択基準: %s\n",
		ActionTypeDone.Value(), ActionTypeDoneDescription, ActionTypeDoneDecisionGuide))
	return b.String()
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
	case string(ActionTypeStruct):
		return ActionTypeStruct, nil
	case string(ActionTypeWrite):
		return ActionTypeWrite, nil
	case string(ActionTypeReview):
		return ActionTypeReview, nil
	case string(ActionTypeDone):
		return ActionTypeDone, nil
	default:
		return "", errors.NewDomainError(errors.ValidationError, "invalid action type")
	}
}
