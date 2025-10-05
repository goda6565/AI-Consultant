package report

import (
	"context"
	"fmt"
	"time"

	reportEntity "github.com/goda6565/ai-consultant/backend/internal/domain/report/entity"
	reportRepository "github.com/goda6565/ai-consultant/backend/internal/domain/report/repository"
	reportValue "github.com/goda6565/ai-consultant/backend/internal/domain/report/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/internal/gen/app"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ReportRepository struct {
	tx   pgx.Tx
	pool *database.AppPool
}

func NewReportRepository(pool *database.AppPool) reportRepository.ReportRepository {
	return &ReportRepository{tx: nil, pool: pool}
}

func (r *ReportRepository) WithTx(tx pgx.Tx) *ReportRepository {
	return &ReportRepository{tx: tx, pool: r.pool}
}

func (r *ReportRepository) FindByProblemID(ctx context.Context, problemID sharedValue.ID) (*reportEntity.Report, error) {
	var q *app.Queries
	if r.tx != nil {
		q = app.New(r.pool).WithTx(r.tx)
	} else {
		q = app.New(r.pool)
	}

	var pID pgtype.UUID
	if err := pID.Scan(problemID.Value()); err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan problem id: %v", err))
	}

	report, err := q.GetReportByProblemID(ctx, pID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to get report by problem id: %v", err))
	}

	entity, err := toEntity(report)
	if err != nil {
		return nil, fmt.Errorf("failed to convert report to entity: %v", err)
	}

	return entity, nil
}

func (r *ReportRepository) Create(ctx context.Context, report *reportEntity.Report) error {
	var q *app.Queries
	if r.tx != nil {
		q = app.New(r.pool).WithTx(r.tx)
	} else {
		q = app.New(r.pool)
	}

	var id pgtype.UUID
	if err := id.Scan(report.GetID().Value()); err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan id: %v", err))
	}

	var problemID pgtype.UUID
	if err := problemID.Scan(report.GetProblemID().Value()); err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan problem id: %v", err))
	}

	content := report.GetContent()
	err := q.CreateReport(ctx, app.CreateReportParams{
		ID:        id,
		ProblemID: problemID,
		Content:   content.Value(),
	})
	if err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to create report: %v", err))
	}

	return nil
}

func (r *ReportRepository) DeleteByProblemID(ctx context.Context, problemID sharedValue.ID) (numDeleted int64, err error) {
	var q *app.Queries
	if r.tx != nil {
		q = app.New(r.pool).WithTx(r.tx)
	} else {
		q = app.New(r.pool)
	}

	var pID pgtype.UUID
	if err := pID.Scan(problemID.Value()); err != nil {
		return 0, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan problem id: %v", err))
	}

	numDeleted, err = q.DeleteReportsByProblemID(ctx, pID)
	if err != nil {
		return 0, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to delete reports by problem id: %v", err))
	}

	return numDeleted, nil
}

func toEntity(report app.Report) (*reportEntity.Report, error) {
	id, err := sharedValue.NewID(report.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create id: %v", err)
	}

	problemID, err := sharedValue.NewID(report.ProblemID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create problem id: %v", err)
	}

	content := reportValue.NewContent(report.Content)

	var createdAt *time.Time
	if report.CreatedAt.Valid {
		createdAt = &report.CreatedAt.Time
	}

	return reportEntity.NewReport(id, problemID, *content, createdAt), nil
}
