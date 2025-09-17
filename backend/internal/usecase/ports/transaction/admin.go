package transaction

import (
	"context"

	documentRepository "github.com/goda6565/ai-consultant/backend/internal/domain/document/repository"
	hearingRepository "github.com/goda6565/ai-consultant/backend/internal/domain/hearing/repository"
	hearingMessageRepository "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/repository"
	problemRepository "github.com/goda6565/ai-consultant/backend/internal/domain/problem/repository"
)

type adminTxKeyType struct{}

var AdminTxKey = adminTxKeyType{}

type AdminUnitOfWork interface {
	DocumentRepository(ctx context.Context) documentRepository.DocumentRepository
	ProblemRepository(ctx context.Context) problemRepository.ProblemRepository
	HearingRepository(ctx context.Context) hearingRepository.HearingRepository
	HearingMessageRepository(ctx context.Context) hearingMessageRepository.HearingMessageRepository
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
