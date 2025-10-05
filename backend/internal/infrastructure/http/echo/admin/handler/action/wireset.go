package action

import "github.com/google/wire"

var Set = wire.NewSet(
	NewListActionHandler,
)
