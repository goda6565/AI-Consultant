package service

import "github.com/google/wire"

var Set = wire.NewSet(
	NewPlanAction,
	NewExternalSearchAction,
	NewInternalSearchAction,
	NewAnalyzeAction,
	NewWriteAction,
	NewReviewAction,
)

var ActionFactorySet = wire.NewSet(
	NewActionFactory,
)
