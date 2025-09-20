package problemfield

import (
	"context"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/repository"
	"github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/internal/gen/app"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/helper"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProblemFieldRepository struct {
	tx   pgx.Tx
	pool *database.AppPool
}

func NewProblemFieldRepository(pool *database.AppPool) repository.ProblemFieldRepository {
	return &ProblemFieldRepository{tx: nil, pool: pool}
}

func (r *ProblemFieldRepository) WithTx(tx pgx.Tx) *ProblemFieldRepository {
	return &ProblemFieldRepository{tx: tx, pool: r.pool}
}

func (r *ProblemFieldRepository) FindByProblemID(ctx context.Context, problemID sharedValue.ID) ([]entity.ProblemField, error) {
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

	rows, err := q.FindByProblemID(ctx, pID)
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to find problem fields by problem id: %v", err))
	}

	entities := make([]entity.ProblemField, len(rows))
	for i, row := range rows {
		e, err := toEntity(row)
		if err != nil {
			return nil, fmt.Errorf("failed to convert problem field to entity: %v", err)
		}
		entities[i] = *e
	}
	return entities, nil
}

func (r *ProblemFieldRepository) Create(ctx context.Context, problemField *entity.ProblemField) error {
	var q *app.Queries
	if r.tx != nil {
		q = app.New(r.pool).WithTx(r.tx)
	} else {
		q = app.New(r.pool)
	}

	var id pgtype.UUID
	if err := id.Scan(problemField.GetID().Value()); err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan id: %v", err))
	}
	var pID pgtype.UUID
	if err := pID.Scan(problemField.GetProblemID().Value()); err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan problem id: %v", err))
	}

	field := problemField.GetField()
	answered := problemField.GetAnswered()
	err := q.CreateProblemField(ctx, app.CreateProblemFieldParams{
		ID:        id,
		ProblemID: pID,
		Field:     field.Value(),
		Answered:  answered.Value(),
	})
	if err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to create problem field: %v", err))
	}
	return nil
}

func (r *ProblemFieldRepository) UpdateAnswered(ctx context.Context, id sharedValue.ID, answered value.Answered) error {
	var q *app.Queries
	if r.tx != nil {
		q = app.New(r.pool).WithTx(r.tx)
	} else {
		q = app.New(r.pool)
	}

	var pfID pgtype.UUID
	if err := pfID.Scan(id.Value()); err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan id: %v", err))
	}
	_, err := q.UpdateAnswered(ctx, app.UpdateAnsweredParams{
		ID:       pfID,
		Answered: answered.Value(),
	})
	if err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to update problem field: %v", err))
	}
	return nil
}

func (r *ProblemFieldRepository) Delete(ctx context.Context, id sharedValue.ID) (numDeleted int64, err error) {
	var q *app.Queries
	if r.tx != nil {
		q = app.New(r.pool).WithTx(r.tx)
	} else {
		q = app.New(r.pool)
	}

	var pfID pgtype.UUID
	if err := pfID.Scan(id.Value()); err != nil {
		return 0, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan id: %v", err))
	}
	numDeleted, err = q.DeleteProblemField(ctx, pfID)
	if helper.IsNoRowsError(err) {
		return 0, nil
	}
	if err != nil {
		return 0, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to delete problem field: %v", err))
	}
	return numDeleted, nil
}

func toEntity(pf app.ProblemField) (*entity.ProblemField, error) {
	id, err := sharedValue.NewID(pf.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create id: %w", err)
	}
	problemID, err := sharedValue.NewID(pf.ProblemID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create problem id: %w", err)
	}
	field, err := value.NewField(pf.Field)
	if err != nil {
		return nil, fmt.Errorf("failed to create field: %w", err)
	}
	answered := value.NewAnswered(pf.Answered)
	createdAt := pf.CreatedAt.Time
	return entity.NewProblemField(id, problemID, *field, *answered, &createdAt), nil
}
