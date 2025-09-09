package chunk

import "github.com/google/wire"

var Set = wire.NewSet(
	NewCreateChunkUseCase,
)
