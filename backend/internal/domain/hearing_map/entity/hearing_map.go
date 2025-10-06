package entity

import (
	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing_map/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type HearingMap struct {
	id        sharedValue.ID
	hearingID sharedValue.ID
	problemID sharedValue.ID
	content   value.Content
}

func NewHearingMap(id sharedValue.ID, hearingID sharedValue.ID, problemID sharedValue.ID, content value.Content) *HearingMap {
	return &HearingMap{id: id, hearingID: hearingID, problemID: problemID, content: content}
}

func (h *HearingMap) GetID() sharedValue.ID {
	return h.id
}

func (h *HearingMap) GetHearingID() sharedValue.ID {
	return h.hearingID
}

func (h *HearingMap) GetProblemID() sharedValue.ID {
	return h.problemID
}

func (h *HearingMap) GetContent() value.Content {
	return h.content
}
