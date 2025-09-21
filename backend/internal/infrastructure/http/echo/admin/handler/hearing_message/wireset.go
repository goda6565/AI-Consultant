package hearingmessage

import "github.com/google/wire"

var Set = wire.NewSet(
	NewListHearingMessageHandler,
)
