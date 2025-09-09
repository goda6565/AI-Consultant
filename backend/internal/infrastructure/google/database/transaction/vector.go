package transaction

import (
	"context"

	chunkRepository "github.com/goda6565/ai-consultant/backend/internal/domain/chunk/repository"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database"
	chunkRepositoryImpl "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/chunk"
	transaction "github.com/goda6565/ai-consultant/backend/internal/usecase/ports/transaction"
	"github.com/jackc/pgx/v5"
)

type VectorUnitOfWork struct {
	pool            *database.VectorPool
	chunkRepository chunkRepository.ChunkRepository
}

func NewVectorUnitOfWork(ctx context.Context, pool *database.VectorPool, chunkRepository chunkRepository.ChunkRepository) transaction.VectorUnitOfWork {
	return &VectorUnitOfWork{pool: pool, chunkRepository: chunkRepository}
}

func (u *VectorUnitOfWork) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := u.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, transaction.VectorTxKey, tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(ctx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (u *VectorUnitOfWork) ChunkRepository(ctx context.Context) chunkRepository.ChunkRepository {
	tx, ok := ctx.Value(transaction.VectorTxKey).(pgx.Tx)
	if !ok {
		panic("tx is not a pgx.Tx")
	}
	impl := u.chunkRepository.(*chunkRepositoryImpl.ChunkRepository)
	if impl == nil {
		panic("chunkRepository is not a chunkRepositoryImpl.ChunkRepository")
	}
	return impl.WithTx(tx)
}
