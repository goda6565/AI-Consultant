package value

import "github.com/goda6565/ai-consultant/backend/internal/domain/errors"

type DocumentType string

const (
	DocumentExtensionPDF      DocumentType = "pdf"
	DocumentExtensionMarkdown DocumentType = "markdown"
	DocumentExtensionCSV      DocumentType = "csv"
)

func (d DocumentType) Equals(other DocumentType) bool {
	return d == other
}

func (d DocumentType) Value() string {
	return string(d)
}

func (d DocumentType) GetExtension() string {
	switch d {
	case DocumentExtensionPDF:
		return "pdf"
	case DocumentExtensionMarkdown:
		return "md"
	case DocumentExtensionCSV:
		return "csv"
	default:
		return ""
	}
}

func NewDocumentType(value string) (DocumentType, error) {
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
