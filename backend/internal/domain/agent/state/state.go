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
	problem         problemEntity.Problem
	content         value.Content
	problemFields   []problemFieldEntity.ProblemField
	hearingMessages []hearingMessageEntity.HearingMessage
	history         value.History
}

func NewState(problem problemEntity.Problem, content value.Content, problemFields []problemFieldEntity.ProblemField, hearingMessages []hearingMessageEntity.HearingMessage, history value.History) *State {
	return &State{problem: problem, content: content, problemFields: problemFields, hearingMessages: hearingMessages, history: history}
}

func (s *State) GetProblem() problemEntity.Problem {
	return s.problem
}

func (s *State) GetContent() value.Content {
	return s.content
}

func (s *State) GetHistory() value.History {
	return s.history
}

func (s *State) SetContent(content value.Content) {
	s.content = content
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

func (s *State) ToPrompt() string {
	var b strings.Builder
	b.WriteString("=== 課題情報 ===\n")
	b.WriteString(fmt.Sprintf("タイトル: %s\n", s.problem.GetTitle().Value()))
	b.WriteString(fmt.Sprintf("説明: %s\n\n", s.problem.GetDescription().Value()))

	b.WriteString("=== ユーザーとのヒアリング履歴 ===\n")
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

	b.WriteString("\n=== 現在の内容 ===\n")
	b.WriteString(s.content.Value())
	b.WriteString("\n")

	b.WriteString("=== 現在の履歴 ===\n")
	b.WriteString(s.history.GetValue())

	// 情報の完全性を分析するためのガイダンスを追加
	b.WriteString("\n=== 情報分析ガイダンス ===\n")
	b.WriteString("上記の情報を分析し、以下の観点で不足している情報を特定してください:\n")
	b.WriteString("- 技術的な詳細や仕様\n")
	b.WriteString("- 業界のベストプラクティスや標準\n")
	b.WriteString("- 類似事例や成功パターン\n")
	b.WriteString("- リスクや課題の詳細\n")
	b.WriteString("- 実装方法や手順\n")
	b.WriteString("- 最新の動向やトレンド\n")

	return b.String()
}
