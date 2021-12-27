package routing

import "github.com/gofiber/fiber/v2"

type Router interface {
	Setup(app *fiber.App)
}
