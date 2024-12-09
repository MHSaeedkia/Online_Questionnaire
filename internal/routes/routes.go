package routes

import (
	"github.com/gofiber/fiber/v2"
	"online-questionnaire/internal/handlers"
	"online-questionnaire/internal/middlewares"
	"online-questionnaire/internal/services"
)

// SetupRoutes registers all routes
func SetupRoutes(app *fiber.App, userService *services.UserService, oauthHandler *handlers.OAuthHandler) {
	userHandler := handlers.NewUserHandler(userService)

	api := app.Group("/api")

	// User routes
	api.Post("/user/signup", middlewares.FixDateOfBirth, userHandler.Signup)
	api.Post("/user/login", userHandler.Login)

	// Google OAuth login
	api.Post("/user/oath", oauthHandler.GoogleLogin)
}
