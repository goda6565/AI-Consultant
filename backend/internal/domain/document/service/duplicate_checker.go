package service

import (
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

func (dc *DuplicateChecker) CheckDuplicateByTitle(title value.Title) (bool, error) {
	existingDoc, err := dc.documentRepository.FindByTitle(title)
	if err != nil {
		return false, fmt.Errorf("failed to find document by title: %w", err)
	}

	return existingDoc != nil, nil
}

func (dc *DuplicateChecker) CheckDuplicateByPath(path *value.StoragePath) (bool, error) {
	if path == nil {
		return false, nil
	}

	existingDoc, err := dc.documentRepository.FindByPath(*path)
	if err != nil {
		return false, fmt.Errorf("failed to find document by path: %w", err)
	}

	return existingDoc != nil, nil
}
