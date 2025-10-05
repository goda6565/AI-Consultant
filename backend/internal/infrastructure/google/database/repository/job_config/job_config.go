package jobconfig

import (
	"context"
	"fmt"

	jobConfigEntity "github.com/goda6565/ai-consultant/backend/internal/domain/job_config/entity"
	jobConfigRepository "github.com/goda6565/ai-consultant/backend/internal/domain/job_config/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/internal/gen/app"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/helper"
	"github.com/jackc/pgx/v5"
)

type JobConfigRepository struct {
	tx   pgx.Tx
	pool *database.AppPool
}

func NewJobConfigRepository(pool *database.AppPool) jobConfigRepository.JobConfigRepository {
	return &JobConfigRepository{tx: nil, pool: pool}
}

func (r *JobConfigRepository) WithTx(tx pgx.Tx) *JobConfigRepository {
	return &JobConfigRepository{tx: tx, pool: r.pool}
}

func (r *JobConfigRepository) FindByProblemID(ctx context.Context, problemID sharedValue.ID) (*jobConfigEntity.JobConfig, error) {
	var q *app.Queries
	if r.tx != nil {
		q = app.New(r.pool).WithTx(r.tx)
	} else {
		q = app.New(r.pool)
	}

	jobConfig, err := q.GetJobConfigByProblemID(ctx, problemID.Value())
	if helper.IsNoRowsError(err) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to get job config by problem id: %v", err))
	}

	return toEntity(jobConfig)
}

func (r *JobConfigRepository) Create(ctx context.Context, jobConfig *jobConfigEntity.JobConfig) error {
	var q *app.Queries
	if r.tx != nil {
		q = app.New(r.pool).WithTx(r.tx)
	} else {
		q = app.New(r.pool)
	}

	err := q.CreateJobConfig(ctx, app.CreateJobConfigParams{
		ID:                   jobConfig.GetID().Value(),
		ProblemID:            jobConfig.GetProblemID().Value(),
		EnableInternalSearch: jobConfig.GetEnableInternalSearch(),
	})
	if err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to create job config: %v", err))
	}

	return nil
}

func (r *JobConfigRepository) Update(ctx context.Context, jobConfig *jobConfigEntity.JobConfig) error {
	var q *app.Queries
	if r.tx != nil {
		q = app.New(r.pool).WithTx(r.tx)
	} else {
		q = app.New(r.pool)
	}

	err := q.UpdateJobConfig(ctx, app.UpdateJobConfigParams{
		ID:                   jobConfig.GetID().Value(),
		EnableInternalSearch: jobConfig.GetEnableInternalSearch(),
	})
	if err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to update job config: %v", err))
	}

	return nil
}

func (r *JobConfigRepository) DeleteByProblemID(ctx context.Context, problemID sharedValue.ID) (numDeleted int64, err error) {
	var q *app.Queries
	if r.tx != nil {
		q = app.New(r.pool).WithTx(r.tx)
	} else {
		q = app.New(r.pool)
	}

	numDeleted, err = q.DeleteJobConfigByProblemID(ctx, problemID.Value())
	if err != nil {
		return 0, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to delete job config by problem id: %v", err))
	}

	return numDeleted, nil
}

func toEntity(jobConfig app.JobConfig) (*jobConfigEntity.JobConfig, error) {
	id, err := sharedValue.NewID(jobConfig.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create id: %w", err)
	}

	problemID, err := sharedValue.NewID(jobConfig.ProblemID)
	if err != nil {
		return nil, fmt.Errorf("failed to create problem id: %w", err)
	}

	return jobConfigEntity.NewJobConfig(id, problemID, jobConfig.EnableInternalSearch), nil
}
