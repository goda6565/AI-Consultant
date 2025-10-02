package value

import (
	"github.com/goda6565/ai-consultant/backend/internal/pkg/uuid"

	"github.com/goda6565/ai-consultant/backend/internal/domain/errors"
)

type ID string

func (id ID) Equals(other ID) bool {
	return id == other
}

func (id ID) Value() string {
	return string(id)
}

func NewID(value string) (ID, error) {
	if err := uuid.Validate(value); err != nil {
		return "", errors.NewDomainError(errors.ValidationError, "Id must be a valid UUID")
	}
	return ID(value), nil
}
