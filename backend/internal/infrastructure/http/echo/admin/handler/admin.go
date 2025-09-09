package handler

import (
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/handler/document"
	gen "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
)

type AdminHandlers struct {
	*document.CreateDocumentHandler
	*document.DeleteDocumentHandler
	*document.GetDocumentHandler
	*document.ListDocumentHandler
}

func NewAdminHandlers(createDocumentHandler *document.CreateDocumentHandler, deleteDocumentHandler *document.DeleteDocumentHandler, getDocumentHandler *document.GetDocumentHandler, listDocumentHandler *document.ListDocumentHandler) gen.StrictServerInterface {
	return &AdminHandlers{
		createDocumentHandler,
		deleteDocumentHandler,
		getDocumentHandler,
		listDocumentHandler,
	}
}
