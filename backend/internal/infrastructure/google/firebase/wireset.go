package firebase

import "github.com/google/wire"

var Set = wire.NewSet(
	NewFirebaseClient,
)
