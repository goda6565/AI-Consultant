package document

import "github.com/google/wire"

var Set = wire.NewSet(
	NewListDocumentHandler,
	NewCreateDocumentHandler,
	NewDeleteDocumentHandler,
	NewGetDocumentHandler,
)
