package controller

import (
	weatherController "seongwoo/go/fiber/controller/weatherController"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func MainController(app *fiber.App) {
	weather := app.Group("/weather", logger.New())

	weatherController.WeatherController(weather)
}
