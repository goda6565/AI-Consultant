package value

import "github.com/goda6565/ai-consultant/backend/internal/domain/errors"

type Phase string

const (
	InitialPhase           Phase = "initial_phase"
	HearingSummaryPhase    Phase = "hearing_summary_phase"
	ProblemDefinitionPhase Phase = "problem_definition_phase"
)

func (p Phase) Equals(other Phase) bool {
	return p == other
}

func (p Phase) Value() string {
	return string(p)
}

func NewPhase(value string) (Phase, error) {
	switch value {
	case "initial_phase":
		return InitialPhase, nil
	case "hearing_summary_phase":
		return HearingSummaryPhase, nil
	case "problem_definition_phase":
		return ProblemDefinitionPhase, nil
	default:
		return "", errors.NewDomainError(errors.ValidationError, "invalid phase")
	}
}
