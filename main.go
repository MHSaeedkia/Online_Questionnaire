package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"log"
	"online-questionnaire/internal/routers"
)

func main() {

	app := fiber.New()
	app.Use(logger.New())

	//app.Use(middleware.MockAuthMiddleware()) // Set mock user ID

	routers.SetupRoutes(app)

	//Serve Swagger UI and documentation
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json") // Adjust the path if necessary
	})

	// Swagger UI for documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	log.Fatal(app.Listen(":8080"))

}
