package main

import (
	"log"
	"online-questionnaire/internal/routers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	app := fiber.New()
	app.Use(logger.New())

	routers.SetupRoutes(app)
	log.Fatal(app.Listen(":8080"))

}
