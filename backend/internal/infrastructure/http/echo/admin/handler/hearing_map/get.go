package hearingmap

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing_map/entity"
	gen "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
	hearingmap "github.com/goda6565/ai-consultant/backend/internal/usecase/hearing_map"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type GetHearingMapHandler struct {
	getHearingMapUseCase hearingmap.GetHearingMapInputPort
}

func NewGetHearingMapHandler(getHearingMapUseCase hearingmap.GetHearingMapInputPort) *GetHearingMapHandler {
	return &GetHearingMapHandler{getHearingMapUseCase: getHearingMapUseCase}
}

func (h *GetHearingMapHandler) GetHearingMap(ctx context.Context, request gen.GetHearingMapRequestObject) (gen.GetHearingMapResponseObject, error) {
	output, err := h.getHearingMapUseCase.Execute(ctx, hearingmap.GetHearingMapUseCaseInput{HearingID: request.HearingId.String()})
	if err != nil {
		return nil, err
	}
	return toHearingMapJSONResponse(output.HearingMap), nil
}

func toHearingMapJSONResponse(m *entity.HearingMap) gen.GetHearingMapResponseObject {
	content := m.GetContent()
	return gen.GetHearingMap200JSONResponse{
		GetHearingMapSuccessJSONResponse: gen.GetHearingMapSuccessJSONResponse{
			Id:        openapi_types.UUID(uuid.MustParse(m.GetID().Value())),
			HearingId: openapi_types.UUID(uuid.MustParse(m.GetHearingID().Value())),
			ProblemId: openapi_types.UUID(uuid.MustParse(m.GetProblemID().Value())),
			Content:   content.Value(),
		},
	}
}
