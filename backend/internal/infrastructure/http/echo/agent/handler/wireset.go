package handler

import (
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/agent/handler/hearing"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	hearing.Set,
	wire.Struct(new(AgentHandlers), "*"),
)

type AgentHandlers struct {
	Hearing hearing.HearingHandlers
}
