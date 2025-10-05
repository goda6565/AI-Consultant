package jobconfig

import (
	"context"

	jobConfigEntity "github.com/goda6565/ai-consultant/backend/internal/domain/job_config/entity"
	gen "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/job_config"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type UpdateJobConfigHandler struct {
	handler jobconfig.UpdateJobConfigInputPort
}

func NewUpdateJobConfigHandler(handler jobconfig.UpdateJobConfigInputPort) *UpdateJobConfigHandler {
	return &UpdateJobConfigHandler{handler: handler}
}

func (h *UpdateJobConfigHandler) UpdateJobConfig(ctx context.Context, request gen.UpdateJobConfigRequestObject) (gen.UpdateJobConfigResponseObject, error) {
	problemID := request.ProblemId.String()
	enableInternalSearch := request.Body.EnableInternalSearch
	output, err := h.handler.Execute(ctx, jobconfig.UpdateJobConfigUseCaseInput{
		ProblemID:            problemID,
		EnableInternalSearch: enableInternalSearch,
	})
	if err != nil {
		return nil, err
	}
	return toUpdateJobConfigJSONResponse(output.JobConfig), nil
}

func toUpdateJobConfigJSONResponse(jobConfig *jobConfigEntity.JobConfig) gen.UpdateJobConfigResponseObject {
	id := jobConfig.GetID()
	problemID := jobConfig.GetProblemID()
	enableInternalSearch := jobConfig.GetEnableInternalSearch()
	return gen.UpdateJobConfig200JSONResponse{
		UpdateJobConfigSuccessJSONResponse: gen.UpdateJobConfigSuccessJSONResponse{
			Id:                   openapi_types.UUID(uuid.MustParse(id.Value())),
			ProblemId:            openapi_types.UUID(uuid.MustParse(problemID.Value())),
			EnableInternalSearch: enableInternalSearch,
		},
	}
}
