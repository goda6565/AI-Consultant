package document

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/document"
)

type CreateDocumentHandler struct {
	createDocumentUseCase document.CreateDocumentInputPort
}

func NewCreateDocumentHandler(createDocumentUseCase document.CreateDocumentInputPort) *CreateDocumentHandler {
	return &CreateDocumentHandler{createDocumentUseCase: createDocumentUseCase}
}

func (h *CreateDocumentHandler) CreateDocument(ctx context.Context, request gen.CreateDocumentRequestObject) (gen.CreateDocumentResponseObject, error) {
	return nil, nil
}
