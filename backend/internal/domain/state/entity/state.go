package entity

import (
	"fmt"
	"strings"

	problemEntity "github.com/goda6565/ai-consultant/backend/internal/domain/problem/entity"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/domain/state/value"
)

type State struct {
	id                 sharedValue.ID
	problem            *problemEntity.Problem
	phase              value.Phase
	statusMessage      value.StatusMessage
	hearingSummary     *value.HearingSummary
	problemDefinitions []ProblemDefinition
}

func NewState(id sharedValue.ID, problem *problemEntity.Problem, phase value.Phase, statusMessage value.StatusMessage, hearingSummaryState *value.HearingSummary) *State {
	return &State{id: id, problem: problem, phase: phase, statusMessage: statusMessage, hearingSummary: hearingSummaryState}
}

func (s *State) GetID() sharedValue.ID {
	return s.id
}

func (s *State) GetProblem() *problemEntity.Problem {
	return s.problem
}

func (s *State) GetPhase() value.Phase {
	return s.phase
}

func (s *State) GetStatusMessage() value.StatusMessage {
	return s.statusMessage
}

func (s *State) ToProblemDefinitionPrompt() string {
	var b strings.Builder
	b.WriteString("【課題情報】\n")
	b.WriteString(fmt.Sprintf("タイトル: %s\n", s.problem.GetTitle().Value()))
	b.WriteString(fmt.Sprintf("説明: %s\n\n", s.problem.GetDescription().Value()))

	b.WriteString("【ヒアリング結果の整理結果】\n")
	b.WriteString(s.hearingSummary.ToPrompt())
	return b.String()
}

func (s *State) SetHearingSummary(hearingSummary *value.HearingSummary) {
	s.hearingSummary = hearingSummary
}

func (s *State) SetProblemDefinitionState(problemDefinitionState []ProblemDefinition) {
	s.problemDefinitions = problemDefinitionState
}

func (s *State) ProceedPhase() {
	switch s.phase {
	case value.InitialPhase:
		s.phase = value.HearingSummaryPhase
	case value.HearingSummaryPhase:
		s.phase = value.ProblemDefinitionPhase
	}
}
