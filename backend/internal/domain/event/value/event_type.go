package value

import "github.com/goda6565/ai-consultant/backend/internal/domain/errors"

type EventType string

const (
	EventTypeAction EventType = "action"
	EventTypeInput  EventType = "input"
	EventTypeOutput EventType = "output"
)

func (e EventType) Equals(other EventType) bool {
	return e == other
}

func (e EventType) Value() string {
	return string(e)
}

func NewEventType(value string) (EventType, error) {
	switch value {
	case "action":
		return EventTypeAction, nil
	case "input":
		return EventTypeInput, nil
	case "output":
		return EventTypeOutput, nil
	default:
		return "", errors.NewDomainError(errors.ValidationError, "invalid event type")
	}
}
