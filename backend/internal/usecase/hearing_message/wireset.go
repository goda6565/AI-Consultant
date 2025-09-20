package hearing_message

import "github.com/google/wire"

var Set = wire.NewSet(
	NewListHearingMessageUseCase,
)
