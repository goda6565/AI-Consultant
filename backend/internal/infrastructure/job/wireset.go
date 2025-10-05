package job

import "github.com/google/wire"

var Set = wire.NewSet(
	NewBaseJob,
)
