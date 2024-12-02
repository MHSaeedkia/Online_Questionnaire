package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"log"
	config "online-questionnaire/configs"
	"online-questionnaire/internal/db"
	"online-questionnaire/internal/handlers"
	"online-questionnaire/internal/repositories"
	"online-questionnaire/internal/services"
)

func SetupRoutes(app *fiber.App, userService *services.UserService) {
	userHandler := handlers.NewUserHandler(userService)

	api := app.Group("/api")
	api.Post("/user/signup", userHandler.Signup)
	//api.Post("/user/login", userHandler.Login)
}

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

	// Initialize the UserRepository
	userRepository := repositories.NewUserRepository(DB) // Assuming you have a constructor for the UserRepository

	// Initialize Fiber app
	app := fiber.New()

	// Initialize services and routes
	userService := services.NewUserService(userRepository, cfg)
	SetupRoutes(app, userService)

	// Serve the Swagger JSON file at /swagger/doc.json
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json") // Adjust the path if necessary
	})

	// Serve the Swagger UI at /swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// App port
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
