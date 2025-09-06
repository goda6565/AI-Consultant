package errors

type DomainErrorType string

const (
	ValidationError DomainErrorType = "validation_error"
)

type DomainError struct {
	ErrorType DomainErrorType
	Message   string
}

func (e *DomainError) Error() string {
	return e.Message
}

func NewDomainError(errorType DomainErrorType, message string) error {
	return &DomainError{
		ErrorType: errorType,
		Message:   message,
	}
}
