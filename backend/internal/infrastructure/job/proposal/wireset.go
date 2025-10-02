package proposal

import "github.com/google/wire"

var Set = wire.NewSet(
	NewExecuteProposal,
)
