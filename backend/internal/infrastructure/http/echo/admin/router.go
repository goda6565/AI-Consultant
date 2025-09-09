package admin

import (
	"encoding/json"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
	echoRouter "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/swaggo/swag"
)

type AdminRouter struct {
	handlers    gen.StrictServerInterface
	environment *environment.Environment
}

func NewAdminRouter(handlers gen.StrictServerInterface, environment *environment.Environment) echoRouter.Router {
	return &AdminRouter{handlers: handlers, environment: environment}
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

func (r *AdminRouter) RegisterRoutes(e *echo.Echo) *echo.Echo {
	_, err := setUpSwagger(e, r.environment)
	if err != nil {
		panic(err)
	}
	gen.RegisterHandlers(e, gen.NewStrictHandler(r.handlers, nil))
	return e
}
