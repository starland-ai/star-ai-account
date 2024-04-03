package middlewares

import (
	"net/http"
	"starland-account/configs"

	"github.com/gofiber/fiber/v2"
)

func Auth() func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		token := ctx.Get("X-Token")
		if configs.GetConfig().Token != token {
			return ctx.SendStatus(http.StatusUnauthorized)
		} else {
			return ctx.Next()
		}
	}
}
