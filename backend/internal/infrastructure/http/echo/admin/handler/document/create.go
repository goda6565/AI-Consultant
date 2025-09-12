package document

import (
	"bytes"
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/document/entity"
	gen "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/document"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type CreateDocumentHandler struct {
	createDocumentUseCase document.CreateDocumentInputPort
}

func NewCreateDocumentHandler(createDocumentUseCase document.CreateDocumentInputPort) *CreateDocumentHandler {
	return &CreateDocumentHandler{createDocumentUseCase: createDocumentUseCase}
}

func (h *CreateDocumentHandler) CreateDocument(ctx context.Context, request gen.CreateDocumentRequestObject) (gen.CreateDocumentResponseObject, error) {
	createDocumentOutput, err := h.createDocumentUseCase.Execute(ctx, document.CreateDocumentUseCaseInput{
		Title:        request.Body.Title,
		DocumentType: string(request.Body.DocumentType),
		File:         bytes.NewReader(request.Body.Data),
	})
	if err != nil {
		return nil, err
	}
	return toCreateDocumentJSONResponse(createDocumentOutput.Document), nil
}

func toCreateDocumentJSONResponse(document *entity.Document) gen.CreateDocumentResponseObject {
	return gen.CreateDocument201JSONResponse{
		CreateDocumentSuccessJSONResponse: gen.CreateDocumentSuccessJSONResponse{
			Id: openapi_types.UUID(uuid.MustParse(document.GetID().Value())),
		},
	}
}
