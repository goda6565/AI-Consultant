package handler

import (
	"errors"
	"net/http"

	domainErrors "github.com/goda6565/ai-consultant/backend/internal/domain/errors"
	infraErrors "github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
	usecaseErrors "github.com/goda6565/ai-consultant/backend/internal/usecase/errors"
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	InternalServerErrorMessage = "Internal Server Error"
	BadRequestErrorMessage     = "Bad Request"
	ConflictErrorMessage       = "Conflict"
	NotFoundErrorMessage       = "Not Found"
)

func CustomErrorHandler(err error, c echo.Context) {
	logger := logger.GetLogger(c.Request().Context())
	var code int
	var message string
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		code = httpErr.Code
		if msg, ok := httpErr.Message.(string); ok {
			message = msg
		}
	} else {
		var infraErr *infraErrors.InfrastructureError
		var domainErr *domainErrors.DomainError
		var usecaseErr *usecaseErrors.UseCaseError
		if errors.As(err, &infraErr) {
			code, message = infraErrToResponse(infraErr)
		} else if errors.As(err, &domainErr) {
			code, message = domainErrToResponse(domainErr)
		} else if errors.As(err, &usecaseErr) {
			code, message = usecaseErrToResponse(usecaseErr)
		} else {
			code, message = http.StatusInternalServerError, "Internal Server Error"
		}
	}

	logger.Error("request failed",
		"code", code,
		"message", message,
		"error", err,
		"method", c.Request().Method,
		"uri", c.Request().RequestURI,
	)
	if err := c.JSON(code, errToResponse(code, message)); err != nil {
		logger.Error("failed to write response", "error", err)
	}
}

func infraErrToResponse(err *infraErrors.InfrastructureError) (int, string) {
	switch err.ErrorType {
	case infraErrors.ExternalServiceError:
		return http.StatusInternalServerError, InternalServerErrorMessage
	default:
		return http.StatusInternalServerError, InternalServerErrorMessage
	}
}

func domainErrToResponse(err *domainErrors.DomainError) (int, string) {
	switch err.ErrorType {
	case domainErrors.ValidationError:
		return http.StatusBadRequest, err.Message
	default:
		return http.StatusInternalServerError, InternalServerErrorMessage
	}
}

func usecaseErrToResponse(err *usecaseErrors.UseCaseError) (int, string) {
	switch err.ErrorType {
	case usecaseErrors.DuplicateError:
		return http.StatusConflict, ConflictErrorMessage
	case usecaseErrors.NotFoundError:
		return http.StatusNotFound, NotFoundErrorMessage
	case usecaseErrors.InternalError:
		return http.StatusInternalServerError, InternalServerErrorMessage
	default:
		return http.StatusInternalServerError, InternalServerErrorMessage
	}
}

func errToResponse(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}
