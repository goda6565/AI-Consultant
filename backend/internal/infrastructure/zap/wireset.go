package zap

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	ProvideZapLogger,
)
