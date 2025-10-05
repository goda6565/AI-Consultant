package event

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/event/entity"
	gen "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/event"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type ListEventHandler struct {
	listEventInputPort event.ListEventInputPort
}

func NewListEventHandler(listEventInputPort event.ListEventInputPort) *ListEventHandler {
	return &ListEventHandler{listEventInputPort: listEventInputPort}
}

func (h *ListEventHandler) ListEvents(ctx context.Context, request gen.ListEventsRequestObject) (gen.ListEventsResponseObject, error) {
	events, err := h.listEventInputPort.Execute(ctx, event.ListEventUseCaseInput{
		ProblemID: request.ProblemId.String(),
	})
	if err != nil {
		return nil, err
	}
	return gen.ListEvents200JSONResponse{
		ListEventsSuccessJSONResponse: gen.ListEventsSuccessJSONResponse{
			Events: toEventsJSONResponse(events.Events),
		},
	}, nil
}

func toEventsJSONResponse(events []entity.Event) []gen.Event {
	eventsJSON := make([]gen.Event, len(events))
	for i, event := range events {
		eventsJSON[i] = toSingleEventJSON(&event)
	}
	return eventsJSON
}

func toSingleEventJSON(event *entity.Event) gen.Event {
	id := openapi_types.UUID(uuid.MustParse(event.ID.Value()))
	actionType := gen.ActionType(event.ActionType.Value())
	eventType := gen.EventType(event.EventType.Value())
	return gen.Event{
		Id:         id,
		ActionType: actionType,
		EventType:  eventType,
		Message:    event.Message.Value(),
	}
}
