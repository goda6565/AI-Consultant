package value

import "github.com/goda6565/ai-consultant/backend/internal/domain/errors"

type DocumentStatus string

const (
	DocumentStatusPending    DocumentStatus = "pending"
	DocumentStatusProcessing DocumentStatus = "processing"
	DocumentStatusDone       DocumentStatus = "done"
	DocumentStatusFailed     DocumentStatus = "failed"
)

func (d DocumentStatus) Equals(other DocumentStatus) bool {
	return d == other
}

func (d DocumentStatus) Value() string {
	return string(d)
}

func NewDocumentStatus(value string) (DocumentStatus, error) {
	switch value {
	case "pending":
		return DocumentStatusPending, nil
	case "processing":
		return DocumentStatusProcessing, nil
	case "done":
		return DocumentStatusDone, nil
	case "failed":
		return DocumentStatusFailed, nil
	default:
		return "", errors.NewDomainError(errors.ValidationError, "invalid document status")
	}
}
