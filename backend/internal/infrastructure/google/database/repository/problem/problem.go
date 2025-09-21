package problem

import (
	"context"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/repository"
	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/internal/gen/app"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/helper"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProblemRepository struct {
	tx   pgx.Tx
	pool *database.AppPool
}

func NewProblemRepository(pool *database.AppPool) repository.ProblemRepository {
	return &ProblemRepository{tx: nil, pool: pool}
}

func (r *ProblemRepository) WithTx(tx pgx.Tx) *ProblemRepository {
	return &ProblemRepository{tx: tx, pool: r.pool}
}

func (r *ProblemRepository) FindAll(ctx context.Context) ([]entity.Problem, error) {
	var q *app.Queries
	if r.tx != nil {
		q = app.New(r.pool).WithTx(r.tx)
	} else {
		q = app.New(r.pool)
	}
	problems, err := q.GetAllProblems(ctx)
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to get all problems: %v", err))
	}
	entities := make([]entity.Problem, len(problems))
	for i, problem := range problems {
		entity, err := toEntity(problem)
		if err != nil {
			return nil, fmt.Errorf("failed to convert problem to entity: %v", err)
		}
		entities[i] = *entity
	}
	return entities, nil
}

func (r *ProblemRepository) Create(ctx context.Context, problem *entity.Problem) error {
	var q *app.Queries
	if r.tx != nil {
		q = app.New(r.pool).WithTx(r.tx)
	} else {
		q = app.New(r.pool)
	}
	var id pgtype.UUID
	if err := id.Scan(problem.GetID().Value()); err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan id: %v", err))
	}
	err := q.CreateProblem(ctx, app.CreateProblemParams{
		ID:          id,
		Title:       problem.GetTitle().Value(),
		Description: problem.GetDescription().Value(),
		Status:      problem.GetStatus().Value(),
	})
	if err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to create problem: %v", err))
	}
	return nil
}

func (r *ProblemRepository) UpdateStatus(ctx context.Context, id sharedValue.ID, status value.Status) error {
	var q *app.Queries
	if r.tx != nil {
		q = app.New(r.pool).WithTx(r.tx)
	} else {
		q = app.New(r.pool)
	}
	var problemID pgtype.UUID
	if err := problemID.Scan(id.Value()); err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan id: %v", err))
	}
	_, err := q.UpdateProblemStatus(ctx, app.UpdateProblemStatusParams{
		ID:     problemID,
		Status: status.Value(),
	})
	if err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to update problem status: %v", err))
	}
	return nil
}

func (r *ProblemRepository) Delete(ctx context.Context, id sharedValue.ID) (numDeleted int64, err error) {
	var q *app.Queries
	if r.tx != nil {
		q = app.New(r.pool).WithTx(r.tx)
	} else {
		q = app.New(r.pool)
	}
	var problemID pgtype.UUID
	if err := problemID.Scan(id.Value()); err != nil {
		return 0, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan id: %v", err))
	}
	numDeleted, err = q.DeleteProblem(ctx, problemID)
	if err != nil {
		return 0, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to delete problem: %v", err))
	}
	return numDeleted, nil
}

func (r *ProblemRepository) FindById(ctx context.Context, id sharedValue.ID) (*entity.Problem, error) {
	var q *app.Queries
	if r.tx != nil {
		q = app.New(r.pool).WithTx(r.tx)
	} else {
		q = app.New(r.pool)
	}
	var problemID pgtype.UUID
	if err := problemID.Scan(id.Value()); err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan id: %v", err))
	}
	problem, err := q.GetProblemById(ctx, problemID)
	if helper.IsNoRowsError(err) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to get problem by id: %v", err))
	}
	return toEntity(problem)
}

func toEntity(problem app.Problem) (*entity.Problem, error) {
	id, err := sharedValue.NewID(problem.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create id: %w", err)
	}
	title, err := value.NewTitle(problem.Title)
	if err != nil {
		return nil, fmt.Errorf("failed to create title: %w", err)
	}
	description, err := value.NewDescription(problem.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create description: %w", err)
	}
	status, err := value.NewStatus(problem.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to create status: %w", err)
	}
	createdAt := problem.CreatedAt.Time

	return entity.NewProblem(id, *title, *description, status, &createdAt), nil
}
