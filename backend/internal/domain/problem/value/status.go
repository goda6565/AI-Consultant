package value

import "github.com/goda6565/ai-consultant/backend/internal/domain/errors"

type Status string

const (
	StatusPending    Status = "pending"
	StatusHearing    Status = "hearing"
	StatusProcessing Status = "processing"
	StatusDone       Status = "done"
	StatusFailed     Status = "failed"
)

func (s Status) Equals(other Status) bool {
	return s == other
}

func (s Status) Value() string {
	return string(s)
}

func NewStatus(value string) (Status, error) {
	switch value {
	case "pending":
		return StatusPending, nil
	case "hearing":
		return StatusHearing, nil
	case "processing":
		return StatusProcessing, nil
	case "done":
		return StatusDone, nil
	default:
		return "", errors.NewDomainError(errors.ValidationError, "invalid status")
	}
}
