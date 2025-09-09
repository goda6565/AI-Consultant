package di

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo"
)

type App struct {
	Server echo.Server
}

func (a *App) StartApp() {
	a.Server.Start()
}

func (a *App) StopApp(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case error:
				// TODO: Log error
			default:
				// TODO: Log panic
			}
		}
	}()
	a.Server.Stop(ctx)
}
