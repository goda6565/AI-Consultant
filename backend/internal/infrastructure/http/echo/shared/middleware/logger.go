package middleware

import (
	log "github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
	"github.com/labstack/echo/v4"
)

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
			ctx := log.WithLogger(req.Context(), logger)
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
