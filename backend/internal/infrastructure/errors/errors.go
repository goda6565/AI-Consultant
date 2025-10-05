package errors

type InfrastructureErrorType string

const (
	ExternalServiceError InfrastructureErrorType = "external_service_error"
	InternalError        InfrastructureErrorType = "internal_error"
	UnauthorizedError    InfrastructureErrorType = "unauthorized_error"
	ForbiddenError       InfrastructureErrorType = "forbidden_error"
	BadRequestError      InfrastructureErrorType = "bad_request_error"
	BadResponseError     InfrastructureErrorType = "bad_response_error"
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
