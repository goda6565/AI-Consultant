package ports

import (
	"context"
	"io"

	"github.com/goda6565/ai-consultant/backend/internal/domain/document/value"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type StoragePort interface {
	Upload(ctx context.Context, info value.StorageInfo, reader io.Reader) error
	Download(ctx context.Context, info value.StorageInfo) (io.ReadCloser, error)
	Delete(ctx context.Context, info value.StorageInfo) error
}
