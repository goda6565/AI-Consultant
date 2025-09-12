package document

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/document/entity"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/document"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type GetDocumentHandler struct {
	getDocumentUseCase document.GetDocumentInputPort
}

func NewGetDocumentHandler(getDocumentUseCase document.GetDocumentInputPort) *GetDocumentHandler {
	return &GetDocumentHandler{getDocumentUseCase: getDocumentUseCase}
}

func (h *GetDocumentHandler) GetDocument(ctx context.Context, request gen.GetDocumentRequestObject) (gen.GetDocumentResponseObject, error) {
	getDocumentOutput, err := h.getDocumentUseCase.Execute(ctx, document.GetDocumentUseCaseInput{DocumentID: request.DocumentId.String()})
	if err != nil {
		return nil, err
	}
	return toDocumentJSONResponse(getDocumentOutput.Document), nil
}

func toDocumentJSONResponse(document *entity.Document) gen.GetDocumentResponseObject {
	return gen.GetDocument200JSONResponse{
		GetDocumentSuccessJSONResponse: gen.GetDocumentSuccessJSONResponse{
			BucketName:     document.GetStorageInfo().BucketName(),
			CreatedAt:      *document.GetCreatedAt(),
			DocumentType:   gen.DocumentType(document.GetDocumentType()),
			DocumentStatus: gen.DocumentStatus(document.GetStatus()),
			Id:             openapi_types.UUID(uuid.MustParse(document.GetID().Value())),
			ObjectName:     document.GetStorageInfo().ObjectName(),
			RetryCount:     document.GetRetryCount().Value(),
			Title:          document.GetTitle().Value(),
			UpdatedAt:      *document.GetUpdatedAt(),
		},
	}
}
