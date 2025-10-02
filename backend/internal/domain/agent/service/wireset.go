package service

import "github.com/google/wire"

var Set = wire.NewSet(
	NewOrchestrator,
	NewSummarizeService,
	NewGoalService,
)
