package value

import (
	"github.com/goda6565/ai-consultant/backend/internal/domain/errors"
)

type DocumentType string

const (
	DocumentTypeStructured   DocumentType = "structured"
	DocumentTypeUnstructured DocumentType = "unstructured"
)

func (d DocumentType) Equals(other DocumentType) bool {
	return d == other
}

func (d DocumentType) Value() string {
	return string(d)
}

func NewDocumentType(value string) (DocumentType, error) {
	switch value {
	case "structured":
		return DocumentTypeStructured, nil
	case "unstructured":
		return DocumentTypeUnstructured, nil
	default:
		return "", errors.NewDomainError(errors.ValidationError, "invalid document type")
	}
}
