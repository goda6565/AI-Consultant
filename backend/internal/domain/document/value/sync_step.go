package value

import "github.com/goda6565/ai-consultant/backend/internal/domain/errors"

type SyncStep string

const (
	SyncStepPending SyncStep = "pending"
	SyncStepVector  SyncStep = "vector"
	SyncStepDone    SyncStep = "done"
)

func (s SyncStep) Equals(other SyncStep) bool {
	return s == other
}

func (s SyncStep) Value() string {
	return string(s)
}

func NewSyncStep(value string) (SyncStep, error) {
	switch value {
	case "pending":
		return SyncStepPending, nil
	case "vector":
		return SyncStepVector, nil
	case "done":
		return SyncStepDone, nil
	default:
		return "", errors.NewDomainError(errors.ValidationError, "invalid sync step")
	}
}
