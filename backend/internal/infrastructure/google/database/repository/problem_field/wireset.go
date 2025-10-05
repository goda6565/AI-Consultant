package problemfield

import "github.com/google/wire"

var Set = wire.NewSet(
	NewProblemFieldRepository,
)
