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
)

type adminTxKeyType struct{}

var AdminTxKey = adminTxKeyType{}

type AdminUnitOfWork interface {
	DocumentRepository(ctx context.Context) documentRepository.DocumentRepository
	ProblemRepository(ctx context.Context) problemRepository.ProblemRepository
	ProblemFieldRepository(ctx context.Context) problemFieldRepository.ProblemFieldRepository
	HearingRepository(ctx context.Context) hearingRepository.HearingRepository
	HearingMessageRepository(ctx context.Context) hearingMessageRepository.HearingMessageRepository
	ActionRepository(ctx context.Context) actionRepository.ActionRepository
	ReportRepository(ctx context.Context) reportRepository.ReportRepository
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
