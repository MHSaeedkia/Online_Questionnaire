package handlers

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"online-questionnaire/internal/services"
)

type OAuthHandler struct {
	Service *services.OAuthService
}

func NewOAuthHandler(service *services.OAuthService) *OAuthHandler {
	return &OAuthHandler{
		Service: service,
	}
}

// GoogleLogin handler to authenticate a user using Google OAuth
func (h *OAuthHandler) GoogleLogin(c *fiber.Ctx) error {
	// Get Google token from the client (sent in the body)
	var request struct {
		GoogleToken string `json:"google_token"`
	}

	if err := c.BodyParser(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Verify the Google token and get user info
	userInfo, err := h.Service.ValidateGoogleToken(request.GoogleToken)
	if err != nil {
		log.Printf("Google token validation error: %v", err)
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid Google token")
	}

	// Generate JWT token for the user
	token, err := h.Service.GenerateJWTToken(userInfo.Email, "oauth-client")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error generating JWT token")
	}

	// Return the JWT token in the response
	return c.JSON(fiber.Map{
		"message": "User authenticated successfully via Google",
		"token":   token,
	})
}
