package search

import (
	"context"

	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type SearchInput struct {
	Query      string
	Embedding  []float32
	NumResults int
}

type SearchResult struct {
	ID         sharedValue.ID
	DocumentID sharedValue.ID
	Content    string
	Similarity float64
}

type SearchOutput struct {
	Results []SearchResult
}

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type SearchPort interface {
	Search(ctx context.Context, input SearchInput) (*SearchOutput, error)
}
