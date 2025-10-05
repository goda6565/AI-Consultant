package handler

import (
	"errors"
	"net/http"

	domainErrors "github.com/goda6565/ai-consultant/backend/internal/domain/errors"
	infraErrors "github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	logger "github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
	usecaseErrors "github.com/goda6565/ai-consultant/backend/internal/usecase/errors"
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	BadRequestErrorMessage400     = "Bad Request"
	UnauthorizedErrorMessage401   = "Unauthorized"
	ForbiddenErrorMessage403      = "Forbidden"
	NotFoundErrorMessage404       = "Not Found"
	ConflictErrorMessage409       = "Conflict"
	InternalServerErrorMessage500 = "Internal Server Error"
)

func CustomErrorHandler(err error, c echo.Context) {
	logger := logger.GetLogger(c.Request().Context())
	var code int
	var message string
	var infraErr *infraErrors.InfrastructureError
	var domainErr *domainErrors.DomainError
	var usecaseErr *usecaseErrors.UseCaseError
	var httpErr *echo.HTTPError
	if errors.As(err, &infraErr) {
		code = infraErrToResponse(infraErr)
	} else if errors.As(err, &domainErr) {
		code = domainErrToResponse(domainErr)
	} else if errors.As(err, &usecaseErr) {
		code = usecaseErrToResponse(usecaseErr)
	} else if errors.As(err, &httpErr) {
		code = httpErr.Code
	} else {
		code = http.StatusInternalServerError
	}
	message = codeToMessage(code)

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

func infraErrToResponse(err *infraErrors.InfrastructureError) int {
	switch err.ErrorType {
	case infraErrors.ExternalServiceError:
		return http.StatusInternalServerError
	case infraErrors.UnauthorizedError:
		return http.StatusUnauthorized
	case infraErrors.ForbiddenError:
		return http.StatusForbidden
	case infraErrors.BadRequestError:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func domainErrToResponse(err *domainErrors.DomainError) int {
	switch err.ErrorType {
	case domainErrors.ValidationError:
		return http.StatusBadRequest
	case domainErrors.InvalidFunctionName:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func usecaseErrToResponse(err *usecaseErrors.UseCaseError) int {
	switch err.ErrorType {
	case usecaseErrors.DuplicateError:
		return http.StatusConflict
	case usecaseErrors.NotFoundError:
		return http.StatusNotFound
	case usecaseErrors.InternalError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func codeToMessage(code int) string {
	switch code {
	case http.StatusBadRequest:
		return BadRequestErrorMessage400
	case http.StatusConflict:
		return ConflictErrorMessage409
	case http.StatusNotFound:
		return NotFoundErrorMessage404
	case http.StatusInternalServerError:
		return InternalServerErrorMessage500
	case http.StatusUnauthorized:
		return UnauthorizedErrorMessage401
	case http.StatusForbidden:
		return ForbiddenErrorMessage403
	default:
		return InternalServerErrorMessage500
	}
}

func errToResponse(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}
