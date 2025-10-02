package transaction

import (
	"context"

	actionRepository "github.com/goda6565/ai-consultant/backend/internal/domain/action/repository"
	documentRepository "github.com/goda6565/ai-consultant/backend/internal/domain/document/repository"
	hearingRepository "github.com/goda6565/ai-consultant/backend/internal/domain/hearing/repository"
	hearingMessageRepository "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/repository"
	problemRepository "github.com/goda6565/ai-consultant/backend/internal/domain/problem/repository"
	problemFieldRepository "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/repository"
	reportRepository "github.com/goda6565/ai-consultant/backend/internal/domain/report/repository"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database"
	actionRepositoryImpl "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/action"
	documentRepositoryImpl "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/document"
	hearingRepositoryImpl "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/hearing"
	hearingMessageRepositoryImpl "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/hearing_message"
	problemRepositoryImpl "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/problem"
	problemFieldRepositoryImpl "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/problem_field"
	reportRepositoryImpl "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/report"
	transaction "github.com/goda6565/ai-consultant/backend/internal/usecase/ports/transaction"
	"github.com/jackc/pgx/v5"
)

type AdminUnitOfWork struct {
	pool                     *database.AppPool
	documentRepository       documentRepository.DocumentRepository
	problemRepository        problemRepository.ProblemRepository
	hearingRepository        hearingRepository.HearingRepository
	hearingMessageRepository hearingMessageRepository.HearingMessageRepository
	problemFieldRepository   problemFieldRepository.ProblemFieldRepository
	actionRepository         actionRepository.ActionRepository
	reportRepository         reportRepository.ReportRepository
}

func NewAdminUnitOfWork(
	ctx context.Context,
	pool *database.AppPool,
	documentRepository documentRepository.DocumentRepository,
	problemRepository problemRepository.ProblemRepository,
	hearingRepository hearingRepository.HearingRepository,
	hearingMessageRepository hearingMessageRepository.HearingMessageRepository,
	problemFieldRepository problemFieldRepository.ProblemFieldRepository,
	actionRepository actionRepository.ActionRepository,
	reportRepository reportRepository.ReportRepository,
) transaction.AdminUnitOfWork {
	return &AdminUnitOfWork{
		pool:                     pool,
		documentRepository:       documentRepository,
		problemRepository:        problemRepository,
		hearingRepository:        hearingRepository,
		hearingMessageRepository: hearingMessageRepository,
		problemFieldRepository:   problemFieldRepository,
		actionRepository:         actionRepository,
		reportRepository:         reportRepository,
	}
}

func (u *AdminUnitOfWork) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := u.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, transaction.AdminTxKey, tx)

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

func (u *AdminUnitOfWork) DocumentRepository(ctx context.Context) documentRepository.DocumentRepository {
	tx, ok := ctx.Value(transaction.AdminTxKey).(pgx.Tx)
	if !ok {
		panic("tx is not a pgx.Tx")
	}
	impl := u.documentRepository.(*documentRepositoryImpl.DocumentRepository)
	if impl == nil {
		panic("documentRepository is not a documentRepositoryImpl.DocumentRepository")
	}
	return impl.WithTx(tx)
}

func (u *AdminUnitOfWork) ProblemRepository(ctx context.Context) problemRepository.ProblemRepository {
	tx, ok := ctx.Value(transaction.AdminTxKey).(pgx.Tx)
	if !ok {
		panic("tx is not a pgx.Tx")
	}
	impl := u.problemRepository.(*problemRepositoryImpl.ProblemRepository)
	if impl == nil {
		panic("problemRepository is not a problemRepositoryImpl.ProblemRepository")
	}
	return impl.WithTx(tx)
}

func (u *AdminUnitOfWork) HearingRepository(ctx context.Context) hearingRepository.HearingRepository {
	tx, ok := ctx.Value(transaction.AdminTxKey).(pgx.Tx)
	if !ok {
		panic("tx is not a pgx.Tx")
	}
	impl := u.hearingRepository.(*hearingRepositoryImpl.HearingRepository)
	if impl == nil {
		panic("hearingRepository is not a hearingRepositoryImpl.HearingRepository")
	}
	return impl.WithTx(tx)
}

func (u *AdminUnitOfWork) HearingMessageRepository(ctx context.Context) hearingMessageRepository.HearingMessageRepository {
	tx, ok := ctx.Value(transaction.AdminTxKey).(pgx.Tx)
	if !ok {
		panic("tx is not a pgx.Tx")
	}
	impl := u.hearingMessageRepository.(*hearingMessageRepositoryImpl.HearingMessageRepository)
	if impl == nil {
		panic("hearingMessageRepository is not a hearingMessageRepositoryImpl.HearingMessageRepository")
	}
	return impl.WithTx(tx)
}

func (u *AdminUnitOfWork) ProblemFieldRepository(ctx context.Context) problemFieldRepository.ProblemFieldRepository {
	tx, ok := ctx.Value(transaction.AdminTxKey).(pgx.Tx)
	if !ok {
		panic("tx is not a pgx.Tx")
	}
	impl := u.problemFieldRepository.(*problemFieldRepositoryImpl.ProblemFieldRepository)
	if impl == nil {
		panic("problemFieldRepository is not a problemFieldRepositoryImpl.ProblemFieldRepository")
	}
	return impl.WithTx(tx)
}

func (u *AdminUnitOfWork) ActionRepository(ctx context.Context) actionRepository.ActionRepository {
	tx, ok := ctx.Value(transaction.AdminTxKey).(pgx.Tx)
	if !ok {
		panic("tx is not a pgx.Tx")
	}
	impl := u.actionRepository.(*actionRepositoryImpl.ActionRepository)
	if impl == nil {
		panic("actionRepository is not a actionRepositoryImpl.ActionRepository")
	}
	return impl.WithTx(tx)
}

func (u *AdminUnitOfWork) ReportRepository(ctx context.Context) reportRepository.ReportRepository {
	tx, ok := ctx.Value(transaction.AdminTxKey).(pgx.Tx)
	if !ok {
		panic("tx is not a pgx.Tx")
	}
	impl := u.reportRepository.(*reportRepositoryImpl.ReportRepository)
	if impl == nil {
		panic("reportRepository is not a reportRepositoryImpl.ReportRepository")
	}
	return impl.WithTx(tx)
}
