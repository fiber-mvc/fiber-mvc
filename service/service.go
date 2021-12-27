package service

import "github.com/gofiber/fiber/v2"

type service interface {
	Setup(app *fiber.App)
}

type Service struct {
	service
}

func (c *Service) Setup(app *fiber.App) {}
