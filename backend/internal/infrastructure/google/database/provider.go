package database

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxvec "github.com/pgvector/pgvector-go/pgx"
)

// VectorPool はベクターデータベース用のプール
type VectorPool struct {
	*pgxpool.Pool
}

// AppPool はアプリケーションデータベース用のプール
type AppPool struct {
	*pgxpool.Pool
}

func ProvideVectorPool(ctx context.Context, e *environment.Environment) (*VectorPool, func()) {
	config, err := pgxpool.ParseConfig(e.VectorDatabaseURL())
	if err != nil {
		panic(err)
	}

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		// Ensure pgvector extension exists before registering types
		if _, err := conn.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS vector"); err != nil {
			return err
		}
		return pgxvec.RegisterTypes(ctx, conn)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		panic(err)
	}

	return &VectorPool{Pool: pool}, func() {
		pool.Close()
	}
}

func ProvideAppPool(ctx context.Context, e *environment.Environment) (*AppPool, func()) {
	config, err := pgxpool.ParseConfig(e.AppDatabaseURL())
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		panic(err)
	}

	return &AppPool{Pool: pool}, func() {
		pool.Close()
	}
}
