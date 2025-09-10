package publish

import "github.com/google/wire"

var Set = wire.NewSet(
	NewPublisher,
)
