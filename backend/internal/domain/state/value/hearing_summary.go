package value

import (
	"github.com/goda6565/ai-consultant/backend/internal/domain/errors"
)

type HearingSummary struct {
	value string
}

func NewHearingSummary(value string) (*HearingSummary, error) {
	if value == "" {
		return nil, errors.NewDomainError(errors.ValidationError, "value is required")
	}
	return &HearingSummary{value: value}, nil
}

func (h *HearingSummary) Value() string {
	return h.value
}

func (h *HearingSummary) ToPrompt() string {
	return h.value
}

func (h *HearingSummary) Equals(other HearingSummary) bool {
	return h.value == other.value
}
