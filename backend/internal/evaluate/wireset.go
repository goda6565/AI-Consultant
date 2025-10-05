package evaluate

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewBaseEvaluator,
)
