package hearing

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewExecuteHearingHandler,
	wire.Struct(new(HearingHandlers), "*"),
)

type HearingHandlers struct {
	Execute ExecuteHearingHandler
}
