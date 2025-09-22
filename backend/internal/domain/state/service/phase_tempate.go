package service

import (
	"context"

	entity "github.com/goda6565/ai-consultant/backend/internal/domain/state/entity"
)

type PhaseTemplate interface {
	Execute(ctx context.Context, state *entity.State) (*entity.State, error)
}
