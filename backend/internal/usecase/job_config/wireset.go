package jobconfig

import "github.com/google/wire"

var Set = wire.NewSet(
	NewUpdateJobConfigUseCase,
	NewGetJobConfigUseCase,
)
