package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	actionEntity "github.com/goda6565/ai-consultant/backend/internal/domain/action/entity"
	actionRepository "github.com/goda6565/ai-consultant/backend/internal/domain/action/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type MemoryActionRepository struct {
	mu      sync.RWMutex
	actions map[string]*actionEntity.Action
}

func NewMemoryActionRepository() actionRepository.ActionRepository {
	return &MemoryActionRepository{
		actions: make(map[string]*actionEntity.Action),
	}
}

func (r *MemoryActionRepository) FindByProblemID(ctx context.Context, problemID sharedValue.ID) ([]actionEntity.Action, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []actionEntity.Action
	for _, action := range r.actions {
		if action.GetProblemID().Equals(problemID) {
			result = append(result, *action)
		}
	}

	// 作成日時でソート（古い順）
	sort.Slice(result, func(i, j int) bool {
		return result[i].GetCreatedAt().Before(*result[j].GetCreatedAt())
	})

	return result, nil
}

func (r *MemoryActionRepository) Create(ctx context.Context, action *actionEntity.Action) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	action.SetCreatedAt(&now)

	r.actions[action.GetID().Value()] = action
	return nil
}

func (r *MemoryActionRepository) DeleteByProblemID(ctx context.Context, problemID sharedValue.ID) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var deleted int64
	for id, action := range r.actions {
		if action.GetProblemID().Equals(problemID) {
			delete(r.actions, id)
			deleted++
		}
	}
	return deleted, nil
}
