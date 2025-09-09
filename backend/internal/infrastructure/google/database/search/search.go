package search

import (
	"context"
	"fmt"

	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/internal/gen/vector"
	searchPort "github.com/goda6565/ai-consultant/backend/internal/usecase/ports/search"
	"github.com/pgvector/pgvector-go"
)

type SearchClient struct {
	pool *database.VectorPool
}

func NewSearchClient(pool *database.VectorPool) searchPort.SearchPort {
	return &SearchClient{pool: pool}
}

func (v *SearchClient) Search(ctx context.Context, input searchPort.SearchInput) (*searchPort.SearchOutput, error) {
	q := vector.New(v.pool)
	pgVector := pgvector.NewVector(input.Embedding)
	// <=> cosine similarity
	rows, err := q.SearchVector(ctx, vector.SearchVectorParams{
		Embedding: pgVector,
		Limit:     int32(input.NumResults),
	})
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to search vector: %v", err))
	}

	results := []searchPort.SearchResult{}
	for _, row := range rows {
		id, err := sharedValue.NewID(row.ID.String())
		if err != nil {
			return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to convert id to sharedValue.ID: %v", err))
		}
		documentID, err := sharedValue.NewID(row.DocumentID.String())
		if err != nil {
			return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to convert documentId to sharedValue.ID: %v", err))
		}
		result := searchPort.SearchResult{
			ID:         id,
			DocumentID: documentID,
			Content:    row.Content,
			Similarity: row.Similarity,
		}
		results = append(results, result)
	}
	return &searchPort.SearchOutput{Results: results}, nil
}
