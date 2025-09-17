package handler

import (
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/handler/document"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/handler/problem"
	gen "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
)

type AdminHandlers struct {
	*document.CreateDocumentHandler
	*document.DeleteDocumentHandler
	*document.GetDocumentHandler
	*document.ListDocumentHandler
	*problem.CreateProblemHandler
	*problem.DeleteProblemHandler
	*problem.GetProblemHandler
	*problem.ListProblemHandler
}

func NewAdminHandlers(
	createDocumentHandler *document.CreateDocumentHandler,
	deleteDocumentHandler *document.DeleteDocumentHandler,
	getDocumentHandler *document.GetDocumentHandler,
	listDocumentHandler *document.ListDocumentHandler,
	createProblemHandler *problem.CreateProblemHandler,
	deleteProblemHandler *problem.DeleteProblemHandler,
	getProblemHandler *problem.GetProblemHandler,
	listProblemHandler *problem.ListProblemHandler) gen.StrictServerInterface {
	return &AdminHandlers{
		createDocumentHandler,
		deleteDocumentHandler,
		getDocumentHandler,
		listDocumentHandler,
		createProblemHandler,
		deleteProblemHandler,
		getProblemHandler,
		listProblemHandler,
	}
}
