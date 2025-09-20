//go:build wireinject
// +build wireinject

package di

import (
	"context"

	"github.com/google/wire"

	chunkService "github.com/goda6565/ai-consultant/backend/internal/domain/chunk/service"
	documentService "github.com/goda6565/ai-consultant/backend/internal/domain/document/service"
	hearingService "github.com/goda6565/ai-consultant/backend/internal/domain/hearing/service"
	hearingMessageService "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/service"
	problemService "github.com/goda6565/ai-consultant/backend/internal/domain/problem/service"
	problemFieldService "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/service"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
	cloudtasksClient "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/cloudtasks"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database"
	chunkRepository "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/chunk"
	documentRepository "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/document"
	hearingRepository "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/hearing"
	hearingMessageRepository "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/hearing_message"
	problemRepository "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/problem"
	problemFieldRepository "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/problem_field"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/transaction"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/firebase"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/gemini"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/ocr"
	storageClient "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/storage"
	baseServer "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo"
	adminRouter "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin"
	adminHandler "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/handler"
	documentHandler "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/handler/document"
	problemHandler "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/handler/problem"
	agentRouter "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/agent"
	agentHandler "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/agent/handler"
	vectorRouter "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/vector"
	vectorHandler "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/vector/handler"
	zap "github.com/goda6565/ai-consultant/backend/internal/infrastructure/zap"
	chunkUseCase "github.com/goda6565/ai-consultant/backend/internal/usecase/chunk"
	documentUseCase "github.com/goda6565/ai-consultant/backend/internal/usecase/document"
	hearingUseCase "github.com/goda6565/ai-consultant/backend/internal/usecase/hearing"
	problemUseCase "github.com/goda6565/ai-consultant/backend/internal/usecase/problem"
)

func InitAdminApplication(ctx context.Context) (*App, func(), error) {
	panic(wire.Build(
		environment.Set,
		zap.Set,
		firebase.Set,
		gemini.Set,
		database.Set,
		chunkRepository.Set,
		documentRepository.Set,
		hearingRepository.Set,
		hearingMessageRepository.Set,
		problemRepository.Set,
		problemFieldRepository.Set,
		transaction.Set,
		storageClient.Set,
		cloudtasksClient.Set,
		documentService.Set,
		problemService.Set,
		problemFieldService.Set,
		documentUseCase.Set,
		problemUseCase.Set,
		documentHandler.Set,
		problemHandler.Set,
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

func InitAgentApplication(ctx context.Context) (*App, func(), error) {
	panic(wire.Build(
		environment.Set,
		zap.Set,
		gemini.Set,
		database.Set,
		transaction.Set,
		documentRepository.Set,
		hearingRepository.Set,
		hearingMessageRepository.Set,
		problemRepository.Set,
		problemFieldRepository.Set,
		problemFieldService.Set,
		hearingService.Set,
		hearingMessageService.Set,
		hearingUseCase.Set,
		problemUseCase.Set,
		agentHandler.Set,
		agentRouter.Set,
		baseServer.Set,
		wire.Struct(new(App), "*"),
	))
}
