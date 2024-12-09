package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"log"
	config "online-questionnaire/configs"
	"online-questionnaire/internal/db"
	"online-questionnaire/internal/handlers"
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

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis address
		Password: "",               // No password by default
		DB:       0,                // Default DB
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Redis is connected successfully")

	// Initialize the UserRepository
	userRepository := repositories.NewUserRepository(DB)

	// Initialize the UserService
	userService := services.NewUserService(userRepository, cfg)

	// Initialize the OAuthService
	oauthService := services.NewOAuthService(cfg, userRepository)

	// Initialize the OAuthHandler
	oauthHandler := handlers.NewOAuthHandler(oauthService)

	// Initialize Fiber app
	app := fiber.New()

	// Setup routes using the external routes file
	routes.SetupRoutes(app, userService, oauthHandler)

	// Serve Swagger UI and documentation
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json") // Adjust the path if necessary
	})

	// Swagger UI for documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Start the app
	if err := app.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)); err != nil {
		log.Fatal(err)
	}
}
