package middleware

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/otyang/go-starter/pkg/logger"
	"github.com/otyang/go-starter/pkg/response"
)

var (
	e1        = "You are seeing this, cause we didnt design "
	e2        = "our system for this endpoint or method. "
	errSystem = e1 + e2 + "Kindly verify you are doing the right thing."
)

// logIt just an utility helper func to log http errors.
// It helps avoid boring repition of codes.
func logIt(c *fiber.Ctx, logger logger.Interface, err error, statusCode, duration int) {
	logger.Error(
		"log-request",
		"method", c.Method(),
		"statusCode", statusCode, // status code always differnt
		"duration", duration,
		"requestID", c.Get("X-Request-Id", ""),
		"Link", string(c.Request().URI().FullURI()),
		"userAgent", c.GetReqHeaders()["User-Agent"],
		"error", err.Error(),
		//"body", c.Response().String(),
	)
}

func ErrorHandler(log logger.Interface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			start    = time.Now()
			err      = c.Next()
			duration = time.Since(start).Milliseconds()
		)

		switch e := err.(type) {
		case nil:
			if c.Response().StatusCode() > http.StatusInternalServerError {
				rsp := response.InternalServerError(e.Error(), "")
				logIt(c, log, err, c.Response().StatusCode(), int(duration))
				return c.Status(http.StatusInternalServerError).JSON(rsp)
			}

		case *response.Response:
			logIt(c, log, err, e.StatusCode, int(duration))
			return c.Status(e.StatusCode).JSON(e)

		case *fiber.Error:
			logIt(c, log, err, e.Code, int(duration))
			return c.Status(e.Code).SendString(e.Error() + "\n\n" + errSystem)

		default:
			rsp := response.InternalServerError(e.Error(), "")
			logIt(c, log, err, c.Response().StatusCode(), int(duration))
			return c.Status(http.StatusInternalServerError).JSON(rsp)
		}

		return nil
	}
}
