package uuid

import (
	"github.com/google/uuid"
)

func Validate(value string) error {
	if _, err := uuid.Parse(value); err != nil {
		return err
	}
	return nil
}
