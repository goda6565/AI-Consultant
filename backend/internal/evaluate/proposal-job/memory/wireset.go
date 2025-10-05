package memory

import "github.com/google/wire"

var Set = wire.NewSet(
	NewMemoryActionRepository,
	NewMemoryReportRepository,
)
