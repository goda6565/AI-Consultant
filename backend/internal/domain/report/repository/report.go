package repository

import (
	"context"

	reportEntity "github.com/goda6565/ai-consultant/backend/internal/domain/report/entity"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type ReportRepository interface {
	FindByProblemID(ctx context.Context, problemID sharedValue.ID) (*reportEntity.Report, error)
	Create(ctx context.Context, report *reportEntity.Report) error
	DeleteByProblemID(ctx context.Context, problemID sharedValue.ID) (numDeleted int64, err error)
}
