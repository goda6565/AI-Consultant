package handler

import (
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/handler/document"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/handler/event"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/handler/hearing"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/handler/hearing_message"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/handler/problem"
	gen "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
)

type AdminRestHandlers struct {
	*document.CreateDocumentHandler
	*document.DeleteDocumentHandler
	*document.GetDocumentHandler
	*document.ListDocumentHandler
	*problem.CreateProblemHandler
	*problem.DeleteProblemHandler
	*problem.GetProblemHandler
	*problem.ListProblemHandler
	*hearing.CreateHearingHandler
	*hearing.GetHearingHandler
	*hearingmessage.ListHearingMessageHandler
	*event.ListEventHandler
}

func NewAdminHandlers(
	createDocumentHandler *document.CreateDocumentHandler,
	deleteDocumentHandler *document.DeleteDocumentHandler,
	getDocumentHandler *document.GetDocumentHandler,
	listDocumentHandler *document.ListDocumentHandler,
	createProblemHandler *problem.CreateProblemHandler,
	deleteProblemHandler *problem.DeleteProblemHandler,
	getProblemHandler *problem.GetProblemHandler,
	listProblemHandler *problem.ListProblemHandler,
	createHearingHandler *hearing.CreateHearingHandler,
	getHearingHandler *hearing.GetHearingHandler,
	listHearingMessageHandler *hearingmessage.ListHearingMessageHandler,
	listEventHandler *event.ListEventHandler,
) gen.StrictServerInterface {
	return &AdminRestHandlers{
		createDocumentHandler,
		deleteDocumentHandler,
		getDocumentHandler,
		listDocumentHandler,
		createProblemHandler,
		deleteProblemHandler,
		getProblemHandler,
		listProblemHandler,
		createHearingHandler,
		getHearingHandler,
		listHearingMessageHandler,
		listEventHandler,
	}
}
