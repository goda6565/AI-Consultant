package admin

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewAdminRouter,
)
