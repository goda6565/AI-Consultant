package jobconfig

import (
	"context"

	jobConfigEntity "github.com/goda6565/ai-consultant/backend/internal/domain/job_config/entity"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/job_config"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type GetJobConfigHandler struct {
	getJobConfigUseCase jobconfig.GetJobConfigInputPort
}

func NewGetJobConfigHandler(getJobConfigUseCase jobconfig.GetJobConfigInputPort) *GetJobConfigHandler {
	return &GetJobConfigHandler{getJobConfigUseCase: getJobConfigUseCase}
}

func (h *GetJobConfigHandler) GetJobConfig(ctx context.Context, request gen.GetJobConfigRequestObject) (gen.GetJobConfigResponseObject, error) {
	getJobConfigOutput, err := h.getJobConfigUseCase.Execute(ctx, jobconfig.GetJobConfigUseCaseInput{ProblemID: request.ProblemId.String()})
	if err != nil {
		return nil, err
	}
	return toGetJobConfigJSONResponse(getJobConfigOutput.JobConfig), nil
}

func toGetJobConfigJSONResponse(jobConfig *jobConfigEntity.JobConfig) gen.GetJobConfigResponseObject {
	id := jobConfig.GetID()
	problemID := jobConfig.GetProblemID()
	enableInternalSearch := jobConfig.GetEnableInternalSearch()
	return gen.GetJobConfig200JSONResponse{
		GetJobConfigSuccessJSONResponse: gen.GetJobConfigSuccessJSONResponse{
			Id:                   openapi_types.UUID(uuid.MustParse(id.Value())),
			ProblemId:            openapi_types.UUID(uuid.MustParse(problemID.Value())),
			EnableInternalSearch: enableInternalSearch,
		},
	}
}
