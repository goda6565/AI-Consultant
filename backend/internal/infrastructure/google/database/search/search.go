package search

import (
	"context"
	"fmt"

	searchClient "github.com/goda6565/ai-consultant/backend/internal/domain/search"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/internal/gen/app"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/internal/gen/vector"
	"github.com/pgvector/pgvector-go"
)

type SearchClient struct {
	vectorPool *database.VectorPool
	appPool    *database.AppPool
}

func NewSearchClient(vectorPool *database.VectorPool, appPool *database.AppPool) searchClient.DocumentSearchClient {
	return &SearchClient{vectorPool: vectorPool, appPool: appPool}
}

func (v *SearchClient) Search(ctx context.Context, input searchClient.DocumentSearchInput) (*searchClient.DocumentSearchOutput, error) {
	vectorQ := vector.New(v.vectorPool)
	appQ := app.New(v.appPool)
	if input.Embedding == nil {
		return nil, errors.NewInfrastructureError(errors.InternalError, "embedding is required")
	}
	pgVector := pgvector.NewVector(*input.Embedding)
	// <=> cosine similarity
	rows, err := vectorQ.SearchVector(ctx, vector.SearchVectorParams{
		Embedding: pgVector,
		Limit:     int32(input.MaxNumResults),
	})
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to search vector: %v", err))
	}

	results := []searchClient.DocumentSearchResult{}
	for _, row := range rows {
		document, err := appQ.GetDocument(ctx, row.DocumentID)
		if err != nil {
			return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to get document: %v", err))
		}
		url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", document.BucketName, document.ObjectName)
		result := searchClient.DocumentSearchResult{
			Title:   document.Title,
			Content: row.ParentContent,
			URL:     url,
		}
		results = append(results, result)
	}
	if input.MaxNumResults > 0 && input.MaxNumResults < len(results) {
		results = results[:input.MaxNumResults]
	}
	return &searchClient.DocumentSearchOutput{Results: results}, nil
}
