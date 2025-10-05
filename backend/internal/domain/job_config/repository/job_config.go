package repository

import (
	"context"

	jobConfigEntity "github.com/goda6565/ai-consultant/backend/internal/domain/job_config/entity"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type JobConfigRepository interface {
	FindByProblemID(ctx context.Context, problemID sharedValue.ID) (*jobConfigEntity.JobConfig, error)
	Create(ctx context.Context, jobConfig *jobConfigEntity.JobConfig) error
	Update(ctx context.Context, jobConfig *jobConfigEntity.JobConfig) error
	DeleteByProblemID(ctx context.Context, problemID sharedValue.ID) (numDeleted int64, err error)
}
