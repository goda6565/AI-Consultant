//go:build wireinject
// +build wireinject

package di

import (
	"context"

	"github.com/google/wire"

	documentService "github.com/goda6565/ai-consultant/backend/internal/domain/document/service"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database"
	chunkRepository "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/chunk"
	documentRepository "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/document"
	storageClient "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/storage"
	baseServer "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo"
	adminRouter "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin"
	adminHandler "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/handler"
	documentHandler "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/handler/document"
	zap "github.com/goda6565/ai-consultant/backend/internal/infrastructure/zap"
	documentUseCase "github.com/goda6565/ai-consultant/backend/internal/usecase/document"
)

func InitAdminApplication(ctx context.Context) (*App, func(), error) {
	panic(wire.Build(
		environment.Set,
		zap.Set,
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
