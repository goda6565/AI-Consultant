package pubsub

import "github.com/google/wire"

var Set = wire.NewSet(
	ProvidePubsubClient,
)
