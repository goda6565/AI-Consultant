package middleware

import (
	"context"

	log "github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
	"github.com/labstack/echo/v4"
)

type loggerKeyType struct{}

var LoggerKey = loggerKeyType{}

func WithLogger(ctx context.Context, logger log.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

func GetLogger(ctx context.Context) log.Logger {
	return ctx.Value(LoggerKey).(log.Logger)
}

func Middleware(logger log.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			requestID := req.Header.Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = res.Header().Get(echo.HeaderXRequestID)
			}

			logger.Info("request",
				"method", req.Method,
				"uri", req.RequestURI,
				"request_id", requestID,
			)

			// set logger to context
			ctx := WithLogger(req.Context(), logger)
			c.SetRequest(req.WithContext(ctx))

			err := next(c)

			logger.Info("response",
				"status", res.Status,
				"size", res.Size,
				"request_id", requestID,
			)
			return err
		}
	}
}
