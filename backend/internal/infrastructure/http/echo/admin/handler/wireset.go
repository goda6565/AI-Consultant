package handler

import (
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/handler/event"
	gen "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewAdminHandlers,
	event.Set,
	wire.Struct(new(AdminHandlers), "*"),
)

type AdminHandlers struct {
	Rest   gen.StrictServerInterface
	Stream event.StreamEventHandler
}
