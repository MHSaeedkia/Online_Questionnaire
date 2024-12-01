package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"online-questionnaire/internal/routers"
)

func main() {

	app := fiber.New()
	app.Use(logger.New())

	routers.SetupRoutes(app)

	log.Fatal(app.Listen(":8080"))

}
