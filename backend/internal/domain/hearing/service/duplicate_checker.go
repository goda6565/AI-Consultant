package service

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type DuplicateCheckerService struct {
	hearingRepository repository.HearingRepository
}

func NewDuplicateCheckerService(hearingRepository repository.HearingRepository) *DuplicateCheckerService {
	return &DuplicateCheckerService{hearingRepository: hearingRepository}
}

func (s *DuplicateCheckerService) Execute(ctx context.Context, problemID sharedValue.ID) (bool, error) {
	existingHearing, err := s.hearingRepository.FindByProblemId(ctx, problemID)
	if err != nil {
		return false, err
	}

	return existingHearing != nil, nil
}
