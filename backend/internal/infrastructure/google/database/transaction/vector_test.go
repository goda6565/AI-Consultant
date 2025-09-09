package transaction

import (
	"context"
	"errors"
	"os"
	"strconv"
	"testing"

	"github.com/goda6565/ai-consultant/backend/internal/domain/chunk/entity"
	chunkValue "github.com/goda6565/ai-consultant/backend/internal/domain/chunk/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
	vectorDatabase "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database"
	chunkRepository "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/chunk"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/uuid"
	portsTransaction "github.com/goda6565/ai-consultant/backend/internal/usecase/ports/transaction"
)

func setup(t *testing.T) (portsTransaction.VectorUnitOfWork, func()) {
	host := os.Getenv("VECTOR_DB_HOST")
	if host == "" {
		t.Skip("VECTOR_DB_HOST is not set")
	}
	port := os.Getenv("VECTOR_DB_PORT")
	if port == "" {
		t.Skip("VECTOR_DB_PORT is not set")
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		t.Fatalf("failed to convert port to int: %v", err)
	}
	username := os.Getenv("VECTOR_DB_USERNAME")
	if username == "" {
		t.Skip("VECTOR_DB_USERNAME is not set")
	}
	password := os.Getenv("VECTOR_DB_PASSWORD")
	if password == "" {
		t.Skip("VECTOR_DB_PASSWORD is not set")
	}
	database := os.Getenv("VECTOR_DB_NAME")
	if database == "" {
		t.Skip("VECTOR_DB_NAME is not set")
	}
	sslMode := os.Getenv("VECTOR_DB_SSL_MODE")
	if sslMode == "" {
		t.Skip("VECTOR_DB_SSL_MODE is not set")
	}
	env := environment.Environment{
		VectorDatabaseEnvironment: environment.VectorDatabaseEnvironment{
			Host:     host,
			Port:     portInt,
			Username: username,
			Password: password,
			Database: database,
			SSLMode:  sslMode,
		},
	}
	ctx := context.Background()
	pool, originalCleanup := vectorDatabase.ProvideVectorPool(ctx, &env)
	// create test table
	_, err = pool.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS vector")
	if err != nil {
		t.Fatalf("failed to create extension: %v", err)
	}
	_, err = pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS vectors (id UUID PRIMARY KEY, document_id UUID, content TEXT, parent_content TEXT, embedding vector(1536))")
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}
	repo := chunkRepository.NewChunkRepository(pool)
	vectorUnitOfWork := NewVectorUnitOfWork(ctx, pool, repo)
	cleanup := func() {
		defer originalCleanup()
		_, err = pool.Exec(ctx, "DROP TABLE IF EXISTS vectors")
		if err != nil {
			t.Fatalf("failed to drop test table: %v", err)
		}
		_, err = pool.Exec(ctx, "DROP EXTENSION IF EXISTS vector")
		if err != nil {
			t.Fatalf("failed to drop extension: %v", err)
		}
	}
	return vectorUnitOfWork, cleanup
}

func TestVectorUnitOfWork_WithTx_Commit(t *testing.T) {
	uow, cleanup := setup(t)
	defer cleanup()

	ctx := context.Background()

	// テスト用のChunkデータを作成
	chunkID, _ := sharedValue.NewID(uuid.NewUUID())
	documentID, _ := sharedValue.NewID(uuid.NewUUID())
	content, err := chunkValue.NewContent("テストコンテンツ")
	if err != nil {
		t.Fatalf("Failed to create content: %v", err)
	}
	parentContent, err := chunkValue.NewContent("親コンテンツ")
	if err != nil {
		t.Fatalf("Failed to create parent content: %v", err)
	}

	// 1536要素のテスト用embeddingを作成
	embeddingData := make([]float32, 1536)
	for i := range embeddingData {
		embeddingData[i] = float32(i % 100)
	}
	embedding, err := chunkValue.NewEmbedding(embeddingData)
	if err != nil {
		t.Fatalf("Failed to create embedding: %v", err)
	}

	chunk := entity.NewChunk(chunkID, documentID, content, parentContent, embedding)

	err = uow.WithTx(ctx, func(ctx context.Context) error {
		repo := uow.ChunkRepository(ctx)
		if repo == nil {
			return errors.New("ChunkRepository should not be nil")
		}

		// Chunkを作成してコミットされることを確認
		return repo.Create(ctx, chunk)
	})

	if err != nil {
		t.Errorf("WithTx should not return error: %v", err)
	}

	// トランザクション外でデータが実際にコミットされているか確認
	err = uow.WithTx(ctx, func(ctx context.Context) error {
		repo := uow.ChunkRepository(ctx)
		numDeleted, err := repo.Delete(ctx, documentID)
		if err != nil {
			return err
		}
		if numDeleted != 1 {
			t.Errorf("Expected 1 deleted record (confirming commit), got %d", numDeleted)
		}
		return nil
	})

	if err != nil {
		t.Errorf("Failed to verify commit: %v", err)
	}
}

func TestVectorUnitOfWork_WithTx_Rollback(t *testing.T) {
	uow, cleanup := setup(t)
	defer cleanup()

	ctx := context.Background()

	// テスト用のChunkデータを作成
	chunkID, _ := sharedValue.NewID(uuid.NewUUID())
	documentID, _ := sharedValue.NewID(uuid.NewUUID())
	content, err := chunkValue.NewContent("ロールバックテストコンテンツ")
	if err != nil {
		t.Fatalf("Failed to create content: %v", err)
	}
	parentContent, err := chunkValue.NewContent("ロールバック親コンテンツ")
	if err != nil {
		t.Fatalf("Failed to create parent content: %v", err)
	}

	// 1536要素のテスト用embeddingを作成
	embeddingData := make([]float32, 1536)
	for i := range embeddingData {
		embeddingData[i] = float32(i % 100)
	}
	embedding, err := chunkValue.NewEmbedding(embeddingData)
	if err != nil {
		t.Fatalf("Failed to create embedding: %v", err)
	}

	chunk := entity.NewChunk(chunkID, documentID, content, parentContent, embedding)
	expectedErr := errors.New("test error")

	err = uow.WithTx(ctx, func(ctx context.Context) error {
		repo := uow.ChunkRepository(ctx)

		// Chunkを作成
		err := repo.Create(ctx, chunk)
		if err != nil {
			return err
		}

		// エラーを返してロールバックをテスト
		return expectedErr
	})

	if err == nil {
		t.Error("WithTx should return error")
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}

	// ロールバックが正しく動作してデータが存在しないことを確認
	err = uow.WithTx(ctx, func(ctx context.Context) error {
		repo := uow.ChunkRepository(ctx)
		numDeleted, err := repo.Delete(ctx, documentID)
		if err != nil {
			return err
		}
		if numDeleted != 0 {
			t.Errorf("Expected 0 deleted records (confirming rollback), got %d", numDeleted)
		}
		return nil
	})

	if err != nil {
		t.Errorf("Failed to verify rollback: %v", err)
	}
}
