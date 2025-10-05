package problem

import "github.com/google/wire"

var Set = wire.NewSet(
	NewListProblemHandler,
	NewCreateProblemHandler,
	NewDeleteProblemHandler,
	NewGetProblemHandler,
)
