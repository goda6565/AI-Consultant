package hearing

import (
	"context"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/internal/gen/app"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/helper"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type HearingRepository struct {
	tx   pgx.Tx
	pool *database.AppPool
}

func NewHearingRepository(pool *database.AppPool) repository.HearingRepository {
	return &HearingRepository{tx: nil, pool: pool}
}

func (r *HearingRepository) WithTx(tx pgx.Tx) *HearingRepository {
	return &HearingRepository{tx: tx, pool: r.pool}
}

func (r *HearingRepository) FindById(ctx context.Context, id sharedValue.ID) (*entity.Hearing, error) {
	var q *app.Queries
	if r.tx != nil {
		q = app.New(r.pool).WithTx(r.tx)
	} else {
		q = app.New(r.pool)
	}
	var hearingID pgtype.UUID
	if err := hearingID.Scan(id.Value()); err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan id: %v", err))
	}
	hearing, err := q.GetHearingById(ctx, hearingID)
	if helper.IsNoRowsError(err) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to get hearing by id: %v", err))
	}
	return toEntity(hearing)
}

func (r *HearingRepository) FindByProblemId(ctx context.Context, problemID sharedValue.ID) (*entity.Hearing, error) {
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
	hearing, err := q.GetHearingByProblemId(ctx, pID)
	if helper.IsNoRowsError(err) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to get hearing by problem id: %v", err))
	}
	return toEntity(hearing)
}

func (r *HearingRepository) FindAllByProblemId(ctx context.Context, problemID sharedValue.ID) ([]entity.Hearing, error) {
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
	hearings, err := q.GetAllHearingsByProblemId(ctx, pID)
	if helper.IsNoRowsError(err) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to get hearing by problem id: %v", err))
	}
	entities := make([]entity.Hearing, len(hearings))
	for i, hearing := range hearings {
		entity, err := toEntity(hearing)
		if err != nil {
			return nil, fmt.Errorf("failed to convert hearing to entity: %v", err)
		}
		entities[i] = *entity
	}
	return entities, nil
}

func (r *HearingRepository) Create(ctx context.Context, hearing *entity.Hearing) error {
	var q *app.Queries
	if r.tx != nil {
		q = app.New(r.pool).WithTx(r.tx)
	} else {
		q = app.New(r.pool)
	}
	var id pgtype.UUID
	if err := id.Scan(hearing.GetID().Value()); err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan id: %v", err))
	}
	var problemID pgtype.UUID
	if err := problemID.Scan(hearing.GetProblemID().Value()); err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan problem id: %v", err))
	}
	err := q.CreateHearing(ctx, app.CreateHearingParams{
		ID:        id,
		ProblemID: problemID,
	})
	if err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to create hearing: %v", err))
	}
	return nil
}

func (r *HearingRepository) DeleteByProblemID(ctx context.Context, problemID sharedValue.ID) (numDeleted int64, err error) {
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
	numDeleted, err = q.DeleteHearingByProblemID(ctx, pID)
	if err != nil {
		return 0, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to delete hearing by problem id: %v", err))
	}
	return numDeleted, nil
}

func toEntity(hearing app.Hearing) (*entity.Hearing, error) {
	id, err := sharedValue.NewID(hearing.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create id: %w", err)
	}
	problemID, err := sharedValue.NewID(hearing.ProblemID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create problem id: %w", err)
	}
	createdAt := hearing.CreatedAt.Time

	return entity.NewHearing(id, problemID, &createdAt), nil
}
