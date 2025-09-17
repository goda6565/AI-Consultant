package entity

import (
	"time"

	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type HearingMessage struct {
	id        sharedValue.ID
	hearingID sharedValue.ID
	role      value.Role
	message   value.Message
	createdAt *time.Time
}

func NewHearingMessage(id sharedValue.ID, hearingID sharedValue.ID, role value.Role, message value.Message, createdAt *time.Time) *HearingMessage {
	return &HearingMessage{id: id, hearingID: hearingID, role: role, message: message, createdAt: createdAt}
}

func (h *HearingMessage) GetID() sharedValue.ID {
	return h.id
}

func (h *HearingMessage) GetHearingID() sharedValue.ID {
	return h.hearingID
}

func (h *HearingMessage) GetRole() value.Role {
	return h.role
}

func (h *HearingMessage) GetMessage() value.Message {
	return h.message
}

func (h *HearingMessage) GetCreatedAt() *time.Time {
	return h.createdAt
}
