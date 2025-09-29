//go:build wireinject
// +build wireinject

package di

import (
	"context"

	"github.com/google/wire"

	actionService "github.com/goda6565/ai-consultant/backend/internal/domain/action/service"
	tools "github.com/goda6565/ai-consultant/backend/internal/domain/action/tools"
	agentService "github.com/goda6565/ai-consultant/backend/internal/domain/agent/service"
	chunkService "github.com/goda6565/ai-consultant/backend/internal/domain/chunk/service"
	documentService "github.com/goda6565/ai-consultant/backend/internal/domain/document/service"
	hearingService "github.com/goda6565/ai-consultant/backend/internal/domain/hearing/service"
	hearingMessageService "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/service"
	problemService "github.com/goda6565/ai-consultant/backend/internal/domain/problem/service"
	problemFieldService "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/service"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
	cloudtasksClient "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/cloudtasks"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database"
	actionRepository "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/action"
	chunkRepository "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/chunk"
	documentRepository "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/document"
	hearingRepository "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/hearing"
	hearingMessageRepository "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/hearing_message"
	problemRepository "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/problem"
	problemFieldRepository "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/repository/problem_field"
	documentSearchClient "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/search"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/transaction"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/firebase"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/gemini"
	googleSearchClient "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/google_search"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/ocr"
	storageClient "github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/storage"
	baseServer "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo"
	adminRouter "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin"
	adminHandler "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/handler"
	documentHandler "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/handler/document"
	hearingHandlerAdmin "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/handler/hearing"
	hearingMessageHandler "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/handler/hearing_message"
	problemHandler "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/handler/problem"
	agentRouter "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/agent"
	agentHandler "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/agent/handler"
	hearingHandlerAgent "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/agent/handler/hearing"
	vectorRouter "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/vector"
	vectorHandler "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/vector/handler"
	baseJob "github.com/goda6565/ai-consultant/backend/internal/infrastructure/job"
	proposalJob "github.com/goda6565/ai-consultant/backend/internal/infrastructure/job/proposal"
	redis "github.com/goda6565/ai-consultant/backend/internal/infrastructure/upstash/redis"
	eventRepository "github.com/goda6565/ai-consultant/backend/internal/infrastructure/upstash/redis/repository/event"
	zap "github.com/goda6565/ai-consultant/backend/internal/infrastructure/zap"
	chunkUseCase "github.com/goda6565/ai-consultant/backend/internal/usecase/chunk"
	documentUseCase "github.com/goda6565/ai-consultant/backend/internal/usecase/document"
	eventUseCase "github.com/goda6565/ai-consultant/backend/internal/usecase/event"
	hearingUseCase "github.com/goda6565/ai-consultant/backend/internal/usecase/hearing"
	hearingMessageUseCase "github.com/goda6565/ai-consultant/backend/internal/usecase/hearing_message"
	problemUseCase "github.com/goda6565/ai-consultant/backend/internal/usecase/problem"
	proposalUseCase "github.com/goda6565/ai-consultant/backend/internal/usecase/proposal"
)

func InitAdminApplication(ctx context.Context) (*App, func(), error) {
	panic(wire.Build(
		environment.Set,
		zap.Set,
		firebase.Set,
		redis.Set,
		gemini.Set,
		database.Set,
		chunkRepository.Set,
		documentRepository.Set,
		hearingRepository.Set,
		hearingMessageRepository.Set,
		problemRepository.Set,
		problemFieldRepository.Set,
		eventRepository.Set,
		transaction.Set,
		storageClient.Set,
		cloudtasksClient.Set,
		documentService.Set,
		problemService.Set,
		problemFieldService.Set,
		documentUseCase.Set,
		problemUseCase.Set,
		hearingUseCase.Set,
		hearingMessageUseCase.Set,
		eventUseCase.Set,
		documentHandler.Set,
		problemHandler.Set,
		hearingHandlerAdmin.Set,
		hearingMessageHandler.Set,
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
		firebase.Set,
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
		hearingHandlerAgent.Set,
		agentHandler.Set,
		agentRouter.Set,
		baseServer.Set,
		wire.Struct(new(App), "*"),
	))
}

func InitProposalJob(ctx context.Context) (*Job, func(), error) {
	panic(wire.Build(
		environment.Set,
		zap.Set,
		gemini.Set,
		database.Set,
		problemRepository.Set,
		problemFieldRepository.Set,
		hearingRepository.Set,
		hearingMessageRepository.Set,
		actionRepository.Set,
		actionService.Set,
		actionService.ActionFactorySet,
		agentService.Set,
		googleSearchClient.Set,
		documentSearchClient.Set,
		tools.Set,
		proposalUseCase.Set,
		proposalJob.Set,
		baseJob.Set,
		wire.Struct(new(Job), "*"),
	))
}
