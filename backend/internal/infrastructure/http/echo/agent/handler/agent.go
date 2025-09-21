package handler

import (
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/agent/handler/hearing"
	gen "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/agent/internal"
)

type AgentHandlers struct {
	*hearing.ExecuteHearingHandler
}

func NewAgentHandlers(
	executeHearingHandler *hearing.ExecuteHearingHandler) gen.StrictServerInterface {
	return &AgentHandlers{
		executeHearingHandler,
	}
}
