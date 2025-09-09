package document

import (
	"context"

	gen "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/document"
)

type DeleteDocumentHandler struct {
	deleteDocumentUseCase document.DeleteDocumentInputPort
}

func NewDeleteDocumentHandler(deleteDocumentUseCase document.DeleteDocumentInputPort) *DeleteDocumentHandler {
	return &DeleteDocumentHandler{deleteDocumentUseCase: deleteDocumentUseCase}
}

func (h *DeleteDocumentHandler) DeleteDocument(ctx context.Context, request gen.DeleteDocumentRequestObject) (gen.DeleteDocumentResponseObject, error) {
	err := h.deleteDocumentUseCase.Execute(ctx, document.DeleteDocumentUseCaseInput{DocumentID: request.DocumentId.String()})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
