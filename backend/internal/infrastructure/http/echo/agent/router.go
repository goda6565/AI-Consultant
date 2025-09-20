package agent

import (
	echoRouter "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/agent/handler"
	"github.com/labstack/echo/v4"
)

type AgentRouter struct {
	handlers *handler.AgentHandlers
}

func NewAgentRouter(handlers *handler.AgentHandlers) echoRouter.Router {
	return &AgentRouter{handlers: handlers}
}

func (r *AgentRouter) RegisterRoutes(e *echo.Echo) *echo.Echo {
	e.GET("/ws/:problemId", echo.HandlerFunc(r.handlers.Hearing.Execute))
	return e
}
