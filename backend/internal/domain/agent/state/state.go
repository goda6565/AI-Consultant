package state

import (
	"fmt"
	"strings"

	actionValue "github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	"github.com/goda6565/ai-consultant/backend/internal/domain/agent/value"
	hearingMessageEntity "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/entity"
	problemEntity "github.com/goda6565/ai-consultant/backend/internal/domain/problem/entity"
	problemFieldEntity "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/entity"
)

type State struct {
	problem            problemEntity.Problem
	goal               value.Goal
	content            value.Content
	problemFields      []problemFieldEntity.ProblemField
	hearingMessages    []hearingMessageEntity.HearingMessage
	history            value.History
	currentAction      actionValue.ActionType
	actionHistory      []actionValue.ActionType
	currentActionCount int
	actionLoopCount    int
}

func NewState(problem problemEntity.Problem, content value.Content, problemFields []problemFieldEntity.ProblemField, hearingMessages []hearingMessageEntity.HearingMessage, history value.History, actionHistory []actionValue.ActionType) *State {
	return &State{problem: problem, content: content, problemFields: problemFields, hearingMessages: hearingMessages, history: history, currentAction: actionValue.ActionTypePlan, actionHistory: actionHistory, currentActionCount: 0, actionLoopCount: 0}
}

func (s *State) GetProblem() problemEntity.Problem {
	return s.problem
}

func (s *State) GetGoal() value.Goal {
	return s.goal
}

func (s *State) GetContent() value.Content {
	return s.content
}

func (s *State) GetHistory() value.History {
	return s.history
}

func (s *State) GetCurrentAction() actionValue.ActionType {
	return s.currentAction
}

func (s *State) ToNextAction(canProceed bool) {
	if canProceed {
		s.currentAction = s.currentAction.Proceed()
		s.currentActionCount = 0
	} else {
		s.currentActionCount++
	}
	s.actionHistory = append(s.actionHistory, s.currentAction)
}

func (s *State) Done() {
	s.currentAction = actionValue.ActionTypeDone
}

func (s *State) GetActionHistory() []actionValue.ActionType {
	return s.actionHistory
}

func (s *State) SetGoal(goal value.Goal) {
	s.goal = goal
}

func (s *State) SetContent(content value.Content) {
	s.content = content
}

func (s *State) IncrementActionLoopCount() {
	s.actionLoopCount++
}

func (s *State) GetActionLoopCount() int {
	return s.actionLoopCount
}

func (s *State) IsInitialAction() bool {
	return len(s.actionHistory) == 0
}

func (s *State) GetCurrentActionCount() int {
	return s.currentActionCount
}

func (s *State) AddHistory(actionType actionValue.ActionType, content string) {
	currentHistory := s.history.GetValue()
	var b strings.Builder
	b.WriteString(currentHistory)
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("=== %s ===\n", actionType.Value()))
	b.WriteString(content)
	b.WriteString("\n")
	newHistory := value.NewHistory(b.String())
	s.history = *newHistory
}

func (s *State) SetHistory(history value.History) {
	s.history = history
}

func (s *State) ToActionHistory() string {
	if len(s.actionHistory) == 0 {
		return ""
	}
	var b strings.Builder
	for i, actionType := range s.actionHistory {
		b.WriteString(fmt.Sprintf("Step%d: %s", i+1, actionType.Value()))
	}
	return b.String()
}

func (s *State) ToGoalPrompt() string {
	var b strings.Builder
	b.WriteString("=== 課題情報 ===\n")
	b.WriteString(fmt.Sprintf("タイトル: %s\n", s.problem.GetTitle().Value()))
	b.WriteString(fmt.Sprintf("説明: %s\n\n", s.problem.GetDescription().Value()))

	b.WriteString("=== ユーザーとのヒアリング履歴（参考情報のみ。ヒアリングは今後することができません。） ===\n")
	fieldMap := make(map[string]string)
	for _, f := range s.problemFields {
		fid := f.GetID()
		fld := f.GetField()
		fieldMap[(&fid).Value()] = (&fld).Value()
	}
	for _, m := range s.hearingMessages {
		msg := m.GetMessage()
		fid := m.GetProblemFieldID()
		fld := fieldMap[(&fid).Value()]
		b.WriteString(fmt.Sprintf("%s: %s [対象項目:%s]\n", m.GetRole().Value(), (&msg).Value(), fld))
	}
	return b.String()
}

func (s *State) ToPrompt() string {
	var b strings.Builder
	b.WriteString("=== 最終ゴール ===\n")
	b.WriteString(s.goal.Value())
	b.WriteString("\n")

	b.WriteString("=== 課題情報 ===\n")
	b.WriteString(fmt.Sprintf("タイトル: %s\n", s.problem.GetTitle().Value()))
	b.WriteString(fmt.Sprintf("説明: %s\n\n", s.problem.GetDescription().Value()))

	b.WriteString("=== ユーザーとのヒアリング履歴（参考情報のみ。ヒアリングは今後することができません。） ===\n")
	fieldMap := make(map[string]string)
	for _, f := range s.problemFields {
		fid := f.GetID()
		fld := f.GetField()
		fieldMap[(&fid).Value()] = (&fld).Value()
	}
	for _, m := range s.hearingMessages {
		msg := m.GetMessage()
		fid := m.GetProblemFieldID()
		fld := fieldMap[(&fid).Value()]
		b.WriteString(fmt.Sprintf("%s: %s [対象項目:%s]\n", m.GetRole().Value(), (&msg).Value(), fld))
	}

	b.WriteString("\n=== 現在のレポートの本文 ===\n")
	b.WriteString(s.content.Value())
	b.WriteString("\n")

	b.WriteString("=== 現在の履歴 ===\n")
	b.WriteString(s.history.GetValue())

	b.WriteString("\n=== 現在のアクション ===\n")
	b.WriteString(s.currentAction.Value())

	b.WriteString("\n=== 現在のアクション履歴 ===\n")
	b.WriteString(s.ToActionHistory())

	b.WriteString("\n=== アクションルート ===\n")
	b.WriteString("**このアクション以外はできないので、今後の計画にこれら以外のActionは考慮しないでください**")
	b.WriteString(actionValue.ActionRoute())

	return b.String()
}
