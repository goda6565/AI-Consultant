package service

import (
	"context"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/document/repository"
	"github.com/goda6565/ai-consultant/backend/internal/domain/document/value"
)

type DuplicateChecker struct {
	documentRepository repository.DocumentRepository
}

func NewDuplicateChecker(documentRepository repository.DocumentRepository) *DuplicateChecker {
	return &DuplicateChecker{
		documentRepository: documentRepository,
	}
}

func (dc *DuplicateChecker) CheckDuplicateByTitle(ctx context.Context, title value.Title) (bool, error) {
	existingDoc, err := dc.documentRepository.FindByTitle(ctx, title)
	if err != nil {
		return false, fmt.Errorf("failed to find document by title: %w", err)
	}

	return existingDoc != nil, nil
}
