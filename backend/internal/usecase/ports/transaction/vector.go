package transaction

import (
	"context"

	chunkRepository "github.com/goda6565/ai-consultant/backend/internal/domain/chunk/repository"
)

type vectorTxKeyType struct{}

var VectorTxKey = vectorTxKeyType{}

type VectorUnitOfWork interface {
	ChunkRepository(ctx context.Context) chunkRepository.ChunkRepository
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
