package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/otyang/go-pkg/utils"
)

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqHeaderID := c.Get("X-Request-Id", "")
		if reqHeaderID == "" {
			newRequestId := utils.RandomID(18)
			c.Set("X-Request-Id", newRequestId)
		}
		return c.Next()
	}
}

func TimeTakenToProcessEndpoint() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			start    = time.Now()
			_        = c.Next()
			duration = time.Since(start)
		)

		c.Set("X-Latency", duration.String())

		return nil
	}
}
