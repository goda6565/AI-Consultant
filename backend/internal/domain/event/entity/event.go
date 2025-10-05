package entity

import (
	actionValue "github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	eventValue "github.com/goda6565/ai-consultant/backend/internal/domain/event/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type Event struct {
	ID         sharedValue.ID
	ProblemID  sharedValue.ID
	EventType  eventValue.EventType
	ActionType actionValue.ActionType
	Message    eventValue.Message
}

func NewEvent(id sharedValue.ID, problemID sharedValue.ID, eventType eventValue.EventType, actionType actionValue.ActionType, message eventValue.Message) *Event {
	return &Event{ID: id, ProblemID: problemID, EventType: eventType, ActionType: actionType, Message: message}
}

func (e *Event) GetID() sharedValue.ID {
	return e.ID
}

func (e *Event) GetProblemID() sharedValue.ID {
	return e.ProblemID
}

func (e *Event) GetEventType() eventValue.EventType {
	return e.EventType
}

func (e *Event) GetActionType() actionValue.ActionType {
	return e.ActionType
}

func (e *Event) GetMessage() eventValue.Message {
	return e.Message
}
