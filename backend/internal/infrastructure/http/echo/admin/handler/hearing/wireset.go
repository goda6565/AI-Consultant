package hearing

import "github.com/google/wire"

var Set = wire.NewSet(
	NewGetHearingHandler,
	NewCreateHearingHandler,
)
