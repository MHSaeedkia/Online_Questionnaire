package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"log"
	config "online-questionnaire/configs"
	"online-questionnaire/internal/db"
	"online-questionnaire/internal/handlers"
)

// @title			online Questionnaire
// @version		1.0
// @description Questionnaire Management System API
func main() {
	cfg, err := config.LoadConfig("./configs/")
	if err != nil {
		log.Fatal(err)
	}

	//database connect
	DB, err := db.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(DB, "is connected successfully")

	// Initialize Fiber app
	app := fiber.New()

	// Serve the Swagger JSON file at /swagger/doc.json
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json") // Adjust the path if necessary
	})

	// Serve the Swagger UI at /swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Handlers
	app.Get("/version", handlers.GetVersion)

	// App port
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
