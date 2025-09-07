package environment

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	ProvideEnvironment,
)
