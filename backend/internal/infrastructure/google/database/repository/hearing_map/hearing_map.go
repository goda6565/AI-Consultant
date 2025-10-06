package hearingmap

import (
	"context"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing_map/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing_map/repository"
	value "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_map/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/internal/gen/app"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type HearingMapRepository struct {
	tx   pgx.Tx
	pool *database.AppPool
}

func NewHearingMapRepository(pool *database.AppPool) repository.HearingMapRepository {
	return &HearingMapRepository{tx: nil, pool: pool}
}

func (r *HearingMapRepository) WithTx(tx pgx.Tx) *HearingMapRepository {
	return &HearingMapRepository{tx: tx, pool: r.pool}
}

func (r *HearingMapRepository) FindByHearingID(ctx context.Context, hearingID sharedValue.ID) (*entity.HearingMap, error) {
	var q *app.Queries
	if r.tx != nil {
		q = app.New(r.pool).WithTx(r.tx)
	} else {
		q = app.New(r.pool)
	}

	var hID pgtype.UUID
	if err := hID.Scan(hearingID.Value()); err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan hearing id: %v", err))
	}

	row, err := q.GetHearingMapByHearingID(ctx, hID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to get hearing map by hearing id: %v", err))
	}

	return toEntity(row)
}

func (r *HearingMapRepository) Create(ctx context.Context, hearingMap *entity.HearingMap) error {
	var q *app.Queries
	if r.tx != nil {
		q = app.New(r.pool).WithTx(r.tx)
	} else {
		q = app.New(r.pool)
	}

	var id pgtype.UUID
	if err := id.Scan(hearingMap.GetID().Value()); err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan id: %v", err))
	}
	var hID pgtype.UUID
	if err := hID.Scan(hearingMap.GetHearingID().Value()); err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan hearing id: %v", err))
	}
	var pID pgtype.UUID
	if err := pID.Scan(hearingMap.GetProblemID().Value()); err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan problem id: %v", err))
	}

	content := hearingMap.GetContent()
	if err := q.CreateHearingMap(ctx, app.CreateHearingMapParams{
		ID:        id,
		HearingID: hID,
		ProblemID: pID,
		Content:   content.Value(),
	}); err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to create hearing map: %v", err))
	}
	return nil
}

func (r *HearingMapRepository) DeleteByHearingID(ctx context.Context, hearingID sharedValue.ID) (numDeleted int64, err error) {
	var q *app.Queries
	if r.tx != nil {
		q = app.New(r.pool).WithTx(r.tx)
	} else {
		q = app.New(r.pool)
	}

	var hID pgtype.UUID
	if err := hID.Scan(hearingID.Value()); err != nil {
		return 0, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan hearing id: %v", err))
	}

	numDeleted, err = q.DeleteHearingMapByHearingID(ctx, hID)
	if err != nil {
		return 0, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to delete hearing map by hearing id: %v", err))
	}
	return numDeleted, nil
}

func toEntity(row app.HearingMap) (*entity.HearingMap, error) {
	id, err := sharedValue.NewID(row.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create id: %w", err)
	}

	hearingID, err := sharedValue.NewID(row.HearingID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create hearing id: %w", err)
	}

	problemID, err := sharedValue.NewID(row.ProblemID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create problem id: %w", err)
	}

	content := value.NewContent(row.Content)
	return entity.NewHearingMap(id, hearingID, problemID, *content), nil
}
