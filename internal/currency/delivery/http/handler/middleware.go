package handler

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/timeout"
)

func TimeoutMiddleware(t time.Duration) func(fiber.Ctx) error {
	return timeout.New(func(c fiber.Ctx) (err error) { return c.Next() }, t)
}
