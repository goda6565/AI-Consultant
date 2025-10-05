package hearing

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing/entity"
	gen "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/hearing"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type GetHearingHandler struct {
	getHearingUseCase hearing.GetHearingInputPort
}

func NewGetHearingHandler(getHearingUseCase hearing.GetHearingInputPort) *GetHearingHandler {
	return &GetHearingHandler{getHearingUseCase: getHearingUseCase}
}

func (h *GetHearingHandler) GetHearing(ctx context.Context, request gen.GetHearingRequestObject) (gen.GetHearingResponseObject, error) {
	getHearingOutput, err := h.getHearingUseCase.Execute(ctx, hearing.GetHearingUseCaseInput{ProblemID: request.ProblemId.String()})
	if err != nil {
		return nil, err
	}
	return toHearingJSONResponse(getHearingOutput.Hearing), nil
}

func toHearingJSONResponse(hearing *entity.Hearing) gen.GetHearingResponseObject {
	return gen.GetHearing200JSONResponse{
		GetHearingSuccessJSONResponse: gen.GetHearingSuccessJSONResponse{
			Id:        openapi_types.UUID(uuid.MustParse(hearing.GetID().Value())),
			ProblemId: openapi_types.UUID(uuid.MustParse(hearing.GetProblemID().Value())),
			CreatedAt: *hearing.GetCreatedAt(),
		},
	}
}
