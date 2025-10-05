package middleware

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/auth"
	echoMiddleware "github.com/oapi-codegen/echo-middleware"
)

const AuthenticatedValueKey = "authenticated_value"

func AuthMiddlewareFunc(authenticator auth.Authenticator) func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		token := input.RequestValidationInput.Request.Header.Get("Authorization")
		if token == "" {
			return errors.NewInfrastructureError(errors.UnauthorizedError, "unauthorized")
		}
		authenticatedValue, err := authenticator.Validate(ctx, token)
		if err != nil {
			return errors.NewInfrastructureError(errors.ForbiddenError, "forbidden")
		}
		eCtx := echoMiddleware.GetEchoContext(ctx)
		eCtx.Set(AuthenticatedValueKey, authenticatedValue)
		return nil
	}
}
