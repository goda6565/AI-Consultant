package report

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/report/entity"
	gen "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/report"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type GetReportHandler struct {
	getReportUseCase report.GetReportInputPort
}

func NewGetReportHandler(getReportUseCase report.GetReportInputPort) *GetReportHandler {
	return &GetReportHandler{getReportUseCase: getReportUseCase}
}

func (h *GetReportHandler) GetReport(ctx context.Context, request gen.GetReportRequestObject) (gen.GetReportResponseObject, error) {
	getReportOutput, err := h.getReportUseCase.Execute(ctx, report.GetReportUseCaseInput{ProblemID: request.ProblemId.String()})
	if err != nil {
		return nil, err
	}
	return toReportJSONResponse(getReportOutput.Report), nil
}

func toReportJSONResponse(report *entity.Report) gen.GetReportResponseObject {
	content := report.GetContent()
	return gen.GetReport200JSONResponse{
		GetReportSuccessJSONResponse: gen.GetReportSuccessJSONResponse{
			Id:        openapi_types.UUID(uuid.MustParse(report.GetID().Value())),
			ProblemId: openapi_types.UUID(uuid.MustParse(report.GetProblemID().Value())),
			Content:   content.Value(),
			CreatedAt: *report.GetCreatedAt(),
		},
	}
}
