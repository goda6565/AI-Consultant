package errors

type InfrastructureErrorType string

const (
	ExternalServiceError InfrastructureErrorType = "external_service_error"
)

type InfrastructureError struct {
	ErrorType InfrastructureErrorType
	Message   string
}

func (e *InfrastructureError) Error() string {
	return e.Message
}

func NewInfrastructureError(errorType InfrastructureErrorType, message string) error {
	return &InfrastructureError{
		ErrorType: errorType,
		Message:   message,
	}
}
