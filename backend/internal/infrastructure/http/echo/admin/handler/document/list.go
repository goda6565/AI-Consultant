package document

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/document/entity"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/document"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type ListDocumentHandler struct {
	listDocumentUseCase document.ListDocumentInputPort
}

func NewListDocumentHandler(listDocumentUseCase document.ListDocumentInputPort) *ListDocumentHandler {
	return &ListDocumentHandler{listDocumentUseCase: listDocumentUseCase}
}

func (h *ListDocumentHandler) ListDocuments(ctx context.Context, request gen.ListDocumentsRequestObject) (gen.ListDocumentsResponseObject, error) {
	listDocumentOutput, err := h.listDocumentUseCase.Execute(ctx)
	if err != nil {
		return nil, err
	}
	return toDocumentsJSONResponse(listDocumentOutput.Documents), nil
}

func toDocumentsJSONResponse(documents []entity.Document) gen.ListDocumentsResponseObject {
	documentsJSON := make([]gen.Document, len(documents))
	for i, document := range documents {
		documentsJSON[i] = toSingleDocumentJSON(&document)
	}
	return gen.ListDocuments200JSONResponse{
		MultipleDocumentsJSONResponse: gen.MultipleDocumentsJSONResponse{
			Documents: documentsJSON,
		},
	}
}

func toSingleDocumentJSON(document *entity.Document) gen.Document {
	return gen.Document{
		BucketName:        document.GetStoragePath().BucketName(),
		CreatedAt:         *document.GetCreatedAt(),
		DocumentExtension: gen.DocumentExtension(document.GetDocumentExtension()),
		DocumentStatus:    gen.DocumentStatus(document.GetStatus()),
		Id:                openapi_types.UUID(uuid.MustParse(document.GetID().Value())),
		ObjectName:        document.GetStoragePath().ObjectName(),
		SyncStep:          gen.SyncStep(document.GetSyncStep()),
		Title:             document.GetTitle().Value(),
		UpdatedAt:         *document.GetUpdatedAt(),
	}
}
