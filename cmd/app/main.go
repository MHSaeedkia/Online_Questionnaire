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

var redisClient *redis.Client
var ctx = context.Background()

// Redis Initialization (inside the main app)
func initRedis() {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Assuming Redis is running locally
	})

	// Check Redis connection
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	} else {
		log.Println("Connected to Redis successfully")
	}
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("./configs/")
	if err != nil {
		log.Fatal(err)
	}

	initRedis()

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
