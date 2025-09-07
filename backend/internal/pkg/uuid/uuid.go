package uuid

import (
	"github.com/google/uuid"
)

func NewUUID() string {
	return uuid.New().String()
}

func Validate(value string) error {
	if _, err := uuid.Parse(value); err != nil {
		return err
	}
	return nil
}
