package agent

import "github.com/google/wire"

var Set = wire.NewSet(
	NewAgentRouter,
)
