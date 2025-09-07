package search

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewSearchClient,
)
