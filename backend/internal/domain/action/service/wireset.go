package service

import "github.com/google/wire"

var Set = wire.NewSet(
	NewPlanAction,
	NewSearchAction,
	NewAnalyzeAction,
	NewWriteAction,
	NewReviewAction,
)

var ActionFactorySet = wire.NewSet(
	NewActionFactory,
)
