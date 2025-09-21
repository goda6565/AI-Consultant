package hearingmessage

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/entity"
	gen "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/hearing_message"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type ListHearingMessageHandler struct {
	listHearingMessageUseCase hearing_message.ListHearingMessageInputPort
}

func NewListHearingMessageHandler(listHearingMessageUseCase hearing_message.ListHearingMessageInputPort) *ListHearingMessageHandler {
	return &ListHearingMessageHandler{listHearingMessageUseCase: listHearingMessageUseCase}
}

func (h *ListHearingMessageHandler) ListHearingMessages(ctx context.Context, request gen.ListHearingMessagesRequestObject) (gen.ListHearingMessagesResponseObject, error) {
	listHearingMessageOutput, err := h.listHearingMessageUseCase.Execute(ctx, hearing_message.ListHearingMessageUseCaseInput{HearingID: request.HearingId.String()})
	if err != nil {
		return nil, err
	}
	return toListHearingMessageJSONResponse(listHearingMessageOutput.HearingMessages), nil
}

func toListHearingMessageJSONResponse(hearingMessages []entity.HearingMessage) gen.ListHearingMessagesResponseObject {
	hearingMessagesJSON := make([]gen.HearingMessage, len(hearingMessages))
	for i, hearingMessage := range hearingMessages {
		hearingMessagesJSON[i] = toSingleHearingMessageJSON(&hearingMessage)
	}
	return gen.ListHearingMessages200JSONResponse{
		ListHearingMessagesSuccessJSONResponse: gen.ListHearingMessagesSuccessJSONResponse{
			HearingMessages: hearingMessagesJSON,
		},
	}
}

func toSingleHearingMessageJSON(hearingMessage *entity.HearingMessage) gen.HearingMessage {
	message := hearingMessage.GetMessage()
	return gen.HearingMessage{
		Id:             openapi_types.UUID(uuid.MustParse(hearingMessage.GetID().Value())),
		HearingId:      openapi_types.UUID(uuid.MustParse(hearingMessage.GetHearingID().Value())),
		ProblemFieldId: openapi_types.UUID(uuid.MustParse(hearingMessage.GetProblemFieldID().Value())),
		Role:           gen.HearingMessageRole(hearingMessage.GetRole().Value()),
		Message:        message.Value(),
		CreatedAt:      *hearingMessage.GetCreatedAt(),
	}
}
