package main

import (
	"log"
	"online-questionnaire/internal/routers"
	"online-questionnaire/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	app := fiber.New()
	app.Use(logger.New())
	app.Use(middleware.MockAuthMiddleware()) // Set mock user ID

	routers.SetupRoutes(app)
	log.Fatal(app.Listen(":8080"))

}
