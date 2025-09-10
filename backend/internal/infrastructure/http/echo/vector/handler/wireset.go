package handler

import (
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/vector/handler/chunk"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	chunk.Set,
	wire.Struct(new(VectorHandlers), "*"),
)

type VectorHandlers struct {
	Chunk chunk.ChunkHandlers
}
