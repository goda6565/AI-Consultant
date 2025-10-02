package value

import "github.com/goda6565/ai-consultant/backend/internal/domain/errors"

type Message struct {
	value string
}

func (m *Message) Value() string {
	return m.value
}

func (m *Message) Equals(other Message) bool {
	return m.value == other.value
}

func NewMessage(value string) (*Message, error) {
	if value == "" {
		return nil, errors.NewDomainError(errors.ValidationError, "value is required")
	}
	return &Message{value: value}, nil
}
