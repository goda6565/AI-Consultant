package chunk

import (
	"context"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/chunk/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/chunk/repository"
	"github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/internal/gen/vector"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pgvector/pgvector-go"
)

type ChunkRepository struct {
	tx   pgx.Tx
	pool *database.VectorPool
}

func NewChunkRepository(pool *database.VectorPool) repository.ChunkRepository {
	return &ChunkRepository{tx: nil, pool: pool}
}

func (v *ChunkRepository) WithTx(tx pgx.Tx) *ChunkRepository {
	return &ChunkRepository{tx: tx, pool: v.pool}
}

func (v *ChunkRepository) Create(ctx context.Context, chunk *entity.Chunk) error {
	var q *vector.Queries
	if v.tx != nil {
		q = vector.New(v.pool).WithTx(v.tx)
	} else {
		q = vector.New(v.pool)
	}

	pgVector := pgvector.NewVector(chunk.GetEmbedding().Value())
	// pgtype.UUID is used to scan the UUID values from the input
	var id, documentID pgtype.UUID
	if err := id.Scan(chunk.GetID().Value()); err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan ID: %v", err))
	}
	if err := documentID.Scan(chunk.GetDocumentID().Value()); err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan DocumentID: %v", err))
	}

	err := q.CreateVector(ctx, vector.CreateVectorParams{
		ID:            id,
		DocumentID:    documentID,
		Content:       chunk.GetContent().Value(),
		ParentContent: chunk.GetParentContent().Value(),
		Embedding:     pgVector,
	})
	if err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to create vector: %v", err))
	}
	return nil
}

func (v *ChunkRepository) Delete(ctx context.Context, documentID value.ID) (numDeleted int64, err error) {
	var q *vector.Queries
	if v.tx != nil {
		q = vector.New(v.pool).WithTx(v.tx)
	} else {
		q = vector.New(v.pool)
	}

	var id pgtype.UUID
	if err := id.Scan(documentID.Value()); err != nil {
		return 0, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan documentID: %v", err))
	}

	numDeleted, err = q.DeleteVector(ctx, id)
	if err != nil {
		return 0, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to delete vector: %v", err))
	}
	return numDeleted, nil
}
