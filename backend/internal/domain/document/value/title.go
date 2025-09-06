package value

import (
	"github.com/goda6565/ai-consultant/backend/internal/domain/errors"
	"unicode/utf8"
)

const (
	maxTitleLength = 50
)

type Title string

func (t Title) Equals(other Title) bool {
	return t == other
}

func (t Title) Value() string {
	return string(t)
}

func NewTitle(value string) (Title, error) {
	if utf8.RuneCountInString(value) > maxTitleLength {
		return "", errors.NewDomainError(errors.ValidationError, "title must be less than 50 characters")
	}
	return Title(value), nil
}
