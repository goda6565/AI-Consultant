package event

import (
	"encoding/json"

	"github.com/goda6565/ai-consultant/backend/internal/pkg/sse"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/event"
	"github.com/labstack/echo/v4"
)

type StreamEventHandler echo.HandlerFunc

type streamEventHandler struct {
	streamEventInputPort event.StreamEventInputPort
}

func NewStreamEventHandler(
	streamEventInputPort event.StreamEventInputPort,
) StreamEventHandler {
	h := &streamEventHandler{
		streamEventInputPort: streamEventInputPort,
	}
	return StreamEventHandler(h.Handle)
}

func (h *streamEventHandler) Handle(c echo.Context) error {
	ctx := c.Request().Context()

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	output, err := h.streamEventInputPort.Execute(ctx, event.StreamEventUseCaseInput{
		ProblemID: c.Param("problemId"),
	})
	if err != nil {
		return err
	}

	for ev := range output.Stream {
		select {
		case <-ctx.Done():
			return nil
		default:
			evJson := EventDataJson{
				ActionType: ev.ActionType.Value(),
				Message:    ev.Message.Value(),
			}
			evJsonBytes, err := json.Marshal(evJson)
			if err != nil {
				return err
			}
			sseEvent := sse.Event{
				ID:    []byte(ev.ID.Value()),
				Event: []byte(ev.EventType),
				Data:  evJsonBytes,
			}
			err = sseEvent.MarshalTo(w)
			if err != nil {
				return err
			}
			w.Flush()
		}
	}

	return nil
}

type EventDataJson struct {
	ActionType string `json:"actionType"`
	Message    string `json:"message"`
}
