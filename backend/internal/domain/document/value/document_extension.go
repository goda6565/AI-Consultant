package value

import "github.com/goda6565/ai-consultant/backend/internal/domain/errors"

type DocumentExtension string

const (
	DocumentExtensionPDF      DocumentExtension = "pdf"
	DocumentExtensionMarkdown DocumentExtension = "markdown"
	DocumentExtensionCSV      DocumentExtension = "csv"
)

func (d DocumentExtension) Equals(other DocumentExtension) bool {
	return d == other
}

func (d DocumentExtension) Value() string {
	return string(d)
}

func NewDocumentExtension(value string) (DocumentExtension, error) {
	switch value {
	case "pdf":
		return DocumentExtensionPDF, nil
	case "markdown":
		return DocumentExtensionMarkdown, nil
	case "csv":
		return DocumentExtensionCSV, nil
	default:
		return "", errors.NewDomainError(errors.ValidationError, "invalid document extension")
	}
}
