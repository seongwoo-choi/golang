package weathercontroller

import (
	weatherService "seongwoo/go/fiber/service/weatherService"

	"github.com/gofiber/fiber/v2"
)

// /weather/*
func WeatherController(router fiber.Router) {
	router.Get("/", weatherService.GetWeather)
}
