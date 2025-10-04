package llmasjudge

import "github.com/google/wire"

var Set = wire.NewSet(
	NewJudge,
)
