package action

import (
	"context"
	"fmt"
	"time"

	actionEntity "github.com/goda6565/ai-consultant/backend/internal/domain/action/entity"
	actionRepository "github.com/goda6565/ai-consultant/backend/internal/domain/action/repository"
	actionValue "github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/internal/gen/app"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ActionRepository struct {
	tx   pgx.Tx
	pool *database.AppPool
}

func NewActionRepository(pool *database.AppPool) actionRepository.ActionRepository {
	return &ActionRepository{tx: nil, pool: pool}
}

func (r *ActionRepository) WithTx(tx pgx.Tx) *ActionRepository {
	return &ActionRepository{tx: tx, pool: r.pool}
}

func (r *ActionRepository) FindByProblemID(ctx context.Context, problemID sharedValue.ID) ([]actionEntity.Action, error) {
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

	actions, err := q.GetActionsByProblemID(ctx, pID)
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to get actions by problem id: %v", err))
	}

	entities := make([]actionEntity.Action, len(actions))
	for i, action := range actions {
		entity, err := toEntity(action)
		if err != nil {
			return nil, fmt.Errorf("failed to convert action to entity: %v", err)
		}
		entities[i] = *entity
	}

	return entities, nil
}

func (r *ActionRepository) Create(ctx context.Context, action *actionEntity.Action) error {
	var q *app.Queries
	if r.tx != nil {
		q = app.New(r.pool).WithTx(r.tx)
	} else {
		q = app.New(r.pool)
	}

	var id pgtype.UUID
	if err := id.Scan(action.GetID().Value()); err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan id: %v", err))
	}

	var problemID pgtype.UUID
	if err := problemID.Scan(action.GetProblemID().Value()); err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan problem id: %v", err))
	}

	input := action.GetInput()
	output := action.GetOutput()
	err := q.CreateAction(ctx, app.CreateActionParams{
		ID:         id,
		ProblemID:  problemID,
		ActionType: action.GetActionType().Value(),
		Input:      input.Value(),
		Output:     output.Value(),
	})
	if err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to create action: %v", err))
	}

	return nil
}

func (r *ActionRepository) DeleteByProblemID(ctx context.Context, problemID sharedValue.ID) (numDeleted int64, err error) {
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

	numDeleted, err = q.DeleteActionsByProblemID(ctx, pID)
	if err != nil {
		return 0, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to delete actions by problem id: %v", err))
	}

	return numDeleted, nil
}

func toEntity(action app.Action) (*actionEntity.Action, error) {
	id, err := sharedValue.NewID(action.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create id: %w", err)
	}

	problemID, err := sharedValue.NewID(action.ProblemID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create problem id: %w", err)
	}

	actionType, err := actionValue.NewActionType(action.ActionType)
	if err != nil {
		return nil, fmt.Errorf("failed to create action type: %w", err)
	}

	input, err := actionValue.NewActionInput(action.Input)
	if err != nil {
		return nil, fmt.Errorf("failed to create action input: %w", err)
	}

	output, err := actionValue.NewActionOutput(action.Output)
	if err != nil {
		return nil, fmt.Errorf("failed to create action output: %w", err)
	}

	var createdAt *time.Time
	if action.CreatedAt.Valid {
		createdAt = &action.CreatedAt.Time
	}

	return actionEntity.NewAction(id, problemID, actionType, *input, *output, createdAt), nil
}
