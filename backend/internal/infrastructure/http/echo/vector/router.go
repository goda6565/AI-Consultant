package vector

import (
	echoRouter "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/vector/handler"
	"github.com/labstack/echo/v4"
)

type VectorRouter struct {
	handlers *handler.VectorHandlers
}

func NewVectorRouter(handlers *handler.VectorHandlers) echoRouter.Router {
	return &VectorRouter{handlers: handlers}
}

func (r *VectorRouter) RegisterRoutes(e *echo.Echo) *echo.Echo {
	e.POST("/webhook", echo.HandlerFunc(r.handlers.Chunk.Create))
	return e
}
