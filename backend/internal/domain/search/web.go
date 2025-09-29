package search

import (
	"context"
)

type WebSearchInput struct {
	Query         string
	MaxNumResults int
}

type WebSearchResult struct {
	Title   string
	Snippet string
	URL     string
}

type WebSearchOutput struct {
	Results []WebSearchResult
}

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type WebSearchClient interface {
	Search(ctx context.Context, input WebSearchInput) (*WebSearchOutput, error)
}
