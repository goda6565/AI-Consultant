package googlesearch

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewGoogleSearchClient,
)
