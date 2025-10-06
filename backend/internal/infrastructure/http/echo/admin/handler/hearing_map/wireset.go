package hearingmap

import "github.com/google/wire"

var Set = wire.NewSet(
	NewGetHearingMapHandler,
)
