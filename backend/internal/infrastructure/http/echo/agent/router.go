package agent

import (
	"encoding/json"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
	echoRouter "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/agent/internal"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/agent/middleware"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/auth"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
	"github.com/labstack/echo/v4"
	oapiMiddleware "github.com/oapi-codegen/echo-middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/swaggo/swag"
)

type AgentRouter struct {
	authenticator auth.Authenticator
	handlers      gen.StrictServerInterface
	environment   *environment.Environment
	logger        logger.Logger
}

func NewAgentRouter(authenticator auth.Authenticator, handlers gen.StrictServerInterface, environment *environment.Environment, logger logger.Logger) echoRouter.Router {
	return &AgentRouter{authenticator: authenticator, handlers: handlers, environment: environment, logger: logger}
}

func setUpSwagger(e *echo.Echo, environment *environment.Environment) (*openapi3.T, error) {
	swagger, err := gen.GetSwagger()
	if err != nil {
		return nil, err
	}

	if environment.Env == "development" {
		swaggerJson, _ := json.Marshal(swagger)
		var SwaggerInfo = &swag.Spec{
			InfoInstanceName: "swagger",
			SwaggerTemplate:  string(swaggerJson),
		}
		swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}
	return swagger, nil
}

func (r *AgentRouter) RegisterRoutes(e *echo.Echo) *echo.Echo {
	e.Use(middleware.CORSMiddleware())
	swagger, err := setUpSwagger(e, r.environment)
	if err != nil {
		panic(err)
	}
	subGroup := e.Group("")

	subGroup.Use(oapiMiddleware.OapiRequestValidatorWithOptions(swagger, &oapiMiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: middleware.AuthMiddlewareFunc(r.authenticator),
		},
	}))
	gen.RegisterHandlers(subGroup, gen.NewStrictHandler(r.handlers, nil))
	return e
}
