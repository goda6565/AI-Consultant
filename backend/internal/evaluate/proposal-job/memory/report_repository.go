package memory

import (
	"context"
	"sync"

	reportEntity "github.com/goda6565/ai-consultant/backend/internal/domain/report/entity"
	reportRepository "github.com/goda6565/ai-consultant/backend/internal/domain/report/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type MemoryReportRepository struct {
	mu      sync.RWMutex
	reports map[string]*reportEntity.Report
}

func NewMemoryReportRepository() reportRepository.ReportRepository {
	return &MemoryReportRepository{
		reports: make(map[string]*reportEntity.Report),
	}
}

func (r *MemoryReportRepository) FindByProblemID(ctx context.Context, problemID sharedValue.ID) (*reportEntity.Report, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, report := range r.reports {
		if report.GetProblemID().Equals(problemID) {
			return report, nil
		}
	}
	return nil, nil
}

func (r *MemoryReportRepository) Create(ctx context.Context, report *reportEntity.Report) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.reports[report.GetID().Value()] = report
	return nil
}

func (r *MemoryReportRepository) DeleteByProblemID(ctx context.Context, problemID sharedValue.ID) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var deleted int64
	for id, report := range r.reports {
		if report.GetProblemID().Equals(problemID) {
			delete(r.reports, id)
			deleted++
		}
	}
	return deleted, nil
}
