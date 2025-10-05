package value

import (
	"unicode/utf8"

	"github.com/goda6565/ai-consultant/backend/internal/domain/errors"
)

const (
	maxTitleLength = 50
)

type Title struct {
	value string
}

func (t Title) Equals(other Title) bool {
	return t.value == other.value
}

func (t Title) Value() string {
	return t.value
}

func NewTitle(value string) (*Title, error) {
	if utf8.RuneCountInString(value) > maxTitleLength {
		return nil, errors.NewDomainError(errors.ValidationError, "title must be less than 50 characters")
	}
	return &Title{value: value}, nil
}
