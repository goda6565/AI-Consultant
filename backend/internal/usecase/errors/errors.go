package errors

type UseCaseErrorType string

const (
	DuplicateError UseCaseErrorType = "duplicate_error"
	NotFoundError  UseCaseErrorType = "not_found_error"
	InternalError  UseCaseErrorType = "internal_error"
)

type UseCaseError struct {
	ErrorType UseCaseErrorType
	Message   string
}

func (e *UseCaseError) Error() string {
	return e.Message
}

func NewUseCaseError(errorType UseCaseErrorType, message string) error {
	return &UseCaseError{
		ErrorType: errorType,
		Message:   message,
	}
}
