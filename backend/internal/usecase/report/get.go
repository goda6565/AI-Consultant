package report

import (
	"context"
	"fmt"

	reportEntity "github.com/goda6565/ai-consultant/backend/internal/domain/report/entity"
	reportRepository "github.com/goda6565/ai-consultant/backend/internal/domain/report/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/errors"
)

type GetReportInputPort interface {
	Execute(ctx context.Context, input GetReportUseCaseInput) (*GetReportUseCaseOutput, error)
}

type GetReportUseCaseInput struct {
	ProblemID string
}

type GetReportUseCaseOutput struct {
	Report *reportEntity.Report
}

type GetReportInteractor struct {
	reportRepository reportRepository.ReportRepository
}

func NewGetReportUseCase(reportRepository reportRepository.ReportRepository) GetReportInputPort {
	return &GetReportInteractor{reportRepository: reportRepository}
}

func (i *GetReportInteractor) Execute(ctx context.Context, input GetReportUseCaseInput) (*GetReportUseCaseOutput, error) {
	problemID, err := sharedValue.NewID(input.ProblemID)
	if err != nil {
		return nil, fmt.Errorf("failed to create problem id: %w", err)
	}
	report, err := i.reportRepository.FindByProblemID(ctx, problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to find report: %w", err)
	}
	if report == nil {
		return nil, errors.NewUseCaseError(errors.NotFoundError, "report not found")
	}
	return &GetReportUseCaseOutput{Report: report}, nil
}
