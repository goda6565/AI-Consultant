//go:build wireinject
// +build wireinject

package di

import (
	"context"

	"github.com/google/wire"

	chunkService "github.com/goda6565/ai-consultant/backend/internal/domain/chunk/service"
	documentService "github.com/goda6565/ai-consultant/backend/internal/domain/document/service"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database"
	chunkRepository "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/chunk"
	documentRepository "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/document"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/transaction"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/firebase"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/gemini"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/ocr"
	pubsubClient "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/pubsub"
	pubsubPublisher "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/pubsub/publish"
	storageClient "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/storage"
	baseServer "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo"
	adminRouter "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin"
	adminHandler "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/handler"
	documentHandler "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/handler/document"
	vectorRouter "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/vector"
	vectorHandler "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/vector/handler"
	zap "github.com/goda6565/ai-consultant/backend/internal/infrastructure/zap"
	chunkUseCase "github.com/goda6565/ai-consultant/backend/internal/usecase/chunk"
	documentUseCase "github.com/goda6565/ai-consultant/backend/internal/usecase/document"
)

func InitAdminApplication(ctx context.Context) (*App, func(), error) {
	panic(wire.Build(
		environment.Set,
		zap.Set,
		firebase.Set,
		pubsubClient.Set,
		pubsubPublisher.Set,
		database.Set,
		chunkRepository.Set,
		documentRepository.Set,
		storageClient.Set,
		documentService.Set,
		documentUseCase.Set,
		documentHandler.Set,
		adminHandler.Set,
		adminRouter.Set,
		baseServer.Set,
		wire.Struct(new(App), "*"),
	))
}

func InitVectorApplication(ctx context.Context) (*App, func(), error) {
	panic(wire.Build(
		environment.Set,
		zap.Set,
		gemini.Set,
		ocr.Set,
		database.Set,
		transaction.Set,
		chunkRepository.Set,
		documentRepository.Set,
		storageClient.Set,
		chunkService.Set,
		chunkUseCase.Set,
		vectorHandler.Set,
		vectorRouter.Set,
		baseServer.Set,
		wire.Struct(new(App), "*"),
	))
}
