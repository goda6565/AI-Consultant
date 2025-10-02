package report

import "github.com/google/wire"

var Set = wire.NewSet(
	NewReportRepository,
)
