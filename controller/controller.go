package controller

import (
	"github.com/gofiber/fiber/v2"
)

type controller interface {
	Setup(app *fiber.App)
}

type Controller struct {
	controller
}

func (c *Controller) Setup(app *fiber.App) {
}
