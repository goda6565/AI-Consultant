package search

import (
	"context"
)

type DocumentSearchInput struct {
	Query         string
	Embedding     *[]float32
	MaxNumResults int
}

type DocumentSearchResult struct {
	Title   string
	Content string
	URL     string
}

type DocumentSearchOutput struct {
	Results []DocumentSearchResult
}

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type DocumentSearchClient interface {
	Search(ctx context.Context, input DocumentSearchInput) (*DocumentSearchOutput, error)
}
