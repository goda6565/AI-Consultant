package hearing

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing/entity"
	gen "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/hearing"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type CreateHearingHandler struct {
	createHearingUseCase hearing.CreateHearingInputPort
}

func NewCreateHearingHandler(createHearingUseCase hearing.CreateHearingInputPort) *CreateHearingHandler {
	return &CreateHearingHandler{createHearingUseCase: createHearingUseCase}
}

func (h *CreateHearingHandler) CreateHearing(ctx context.Context, request gen.CreateHearingRequestObject) (gen.CreateHearingResponseObject, error) {
	createHearingOutput, err := h.createHearingUseCase.Create(ctx, hearing.CreateHearingUseCaseInput{
		ProblemID: request.ProblemId.String(),
	}, nil)
	if err != nil {
		return nil, err
	}
	return toCreateHearingJSONResponse(createHearingOutput.Hearing), nil
}

func toCreateHearingJSONResponse(hearing *entity.Hearing) gen.CreateHearingResponseObject {
	return gen.CreateHearing201JSONResponse{
		CreateHearingSuccessJSONResponse: gen.CreateHearingSuccessJSONResponse{
			HearingId: openapi_types.UUID(uuid.MustParse(hearing.GetID().Value())),
		},
	}
}
