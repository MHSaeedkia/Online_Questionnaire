package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"log"
	config "online-questionnaire/configs"
	"online-questionnaire/internal/db"
	"online-questionnaire/internal/repositories"
	"online-questionnaire/internal/routes"
	"online-questionnaire/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("./configs/")
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the database
	DB, err := db.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(DB, "is connected successfully")

	// Initialize the UserRepository
	userRepository := repositories.NewUserRepository(DB)

	// Initialize the UserService
	userService := services.NewUserService(userRepository, cfg)

	// Initialize Fiber app
	app := fiber.New()

	// Setup routes using the external routes file
	routes.SetupRoutes(app, userService)

	// Serve Swagger UI and documentation
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json") // Adjust the path if necessary
	})

	// Swagger UI for documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Start the app
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
