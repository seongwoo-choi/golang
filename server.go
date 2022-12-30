package main

import (
	"log"
	controller "seongwoo/go/fiber/controller/mainController"
	"seongwoo/go/fiber/myutils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/robfig/cron"
)

func main() {
	app := fiber.New()
	c := cron.New()

	app.Use(cors.New())

	myutils.EightCron(c)
	myutils.SeventeentCron(c)

	controller.MainController(app)

	log.Fatal(app.Listen(":3000"))
}
