package main

import (
	"log"
	"online-questionnaire/configs"
	logging "online-questionnaire/internal/logger"
	"online-questionnaire/internal/routers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

func main() {

	cfg, err := configs.LoadConfig("../../configs/")
	if err != nil {
		log.Fatal(err)
	}
	logErr := logging.NewLogger(cfg, "questionnaire")
	if logErr != nil {
		log.Fatal(logErr)
	}

	logging.GetLogger().Info("Application started", "questionnaire", logging.Logctx{})

	app := fiber.New()
	app.Use(logger.New())

	//app.Use(middleware.MockAuthMiddleware()) // Set mock user ID

	routers.SetupRoutes(cfg, app)

	//Serve Swagger UI and documentation
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json") // Adjust the path if necessary
	})

	// Swagger UI for documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	log.Fatal(app.Listen(":8080"))

}
