package value

import "github.com/goda6565/ai-consultant/backend/internal/domain/errors"

type ActionOutput struct {
	value string
}

func NewActionOutput(value string) (*ActionOutput, error) {
	if value == "" {
		return nil, errors.NewDomainError(errors.ValidationError, "value is required")
	}
	return &ActionOutput{value: value}, nil
}

func (a *ActionOutput) Value() string {
	return a.value
}

func (a *ActionOutput) Equals(other ActionOutput) bool {
	return a.value == other.value
}
