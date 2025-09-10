package chunk

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewCreateChunkHandler,
	wire.Struct(new(ChunkHandlers), "*"),
)

type ChunkHandlers struct {
	Create CreateHandler
}
