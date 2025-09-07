package ocr

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewDocumentAIClient,
)
