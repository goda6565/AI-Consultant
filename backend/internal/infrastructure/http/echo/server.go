package echo

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/shared/handler"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/shared/middleware"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
	"github.com/labstack/echo/v4"
)

type Router interface {
	RegisterRoutes(e *echo.Echo) *echo.Echo
}

type Server interface {
	Start()
	Stop(ctx context.Context)
}

type BaseServer struct {
	logger      logger.Logger
	echo        *echo.Echo
	environment *environment.Environment
}

func NewBaseServer(environment *environment.Environment, logger logger.Logger, router Router) Server {
	echo := echo.New()
	echo.Use(middleware.Middleware(logger))
	echo.HTTPErrorHandler = handler.CustomErrorHandler
	echo = router.RegisterRoutes(echo)
	return &BaseServer{logger: logger, echo: echo, environment: environment}
}

func (s *BaseServer) Start() {
	addr := s.environment.ListenAddress
	s.logger.Info("Starting server", "address", addr)
	go func() {
		if err := s.echo.Start(addr); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit // block until signal is received

	_, cancel := context.WithTimeout(context.Background(), s.environment.ShutdownTimeout)
	defer cancel()
}

func (s *BaseServer) Stop(ctx context.Context) {
	err := s.echo.Shutdown(ctx)
	if err != nil {
		panic(err)
	}
}
