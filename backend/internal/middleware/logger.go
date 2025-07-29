package middleware

import (
	"time"

	"github.com/defi-dashboard/backend/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log request details
		logger.Info("HTTP request",
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"latency", latency.Milliseconds(),
			"ip", c.IP(),
			"user_agent", c.Get("User-Agent"),
			"request_id", c.Locals("requestid"),
			"address", c.Locals("address"), // Will be nil for unauthenticated requests
		)

		return err
	}
}