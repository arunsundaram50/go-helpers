package fiber

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var F = fmt.Sprintf

type KEY_VALUES map[string]any

func JSONResponse(c *fiber.Ctx, statusCode int, keyValuePairs *KEY_VALUES) error {
	return c.Status(statusCode).JSON(keyValuePairs)
}
