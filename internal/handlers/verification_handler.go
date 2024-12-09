package handlers

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"online-questionnaire/internal/services"
)

// VerificationHandler struct
type VerificationHandler struct {
	Service *services.VerificationService
}

// NewVerificationHandler creates a new instance of VerificationHandler
func NewVerificationHandler(service *services.VerificationService) *VerificationHandler {
	return &VerificationHandler{
		Service: service,
	}
}

// SendVerificationCode handler to generate and send the verification code to the user's email
func (h *VerificationHandler) SendVerificationCode(c *fiber.Ctx) error {
	// Get the user's email from the request
	var request struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request data")
	}

	// Send verification code to the user's email
	err := h.Service.SendVerificationCode(request.Email)
	if err != nil {
		log.Printf("Error sending verification code: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to send verification code")
	}

	return c.JSON(fiber.Map{
		"message": "Verification code sent successfully",
	})
}

// ValidateVerificationCode handler to validate the verification code entered by the user
func (h *VerificationHandler) ValidateVerificationCode(c *fiber.Ctx) error {
	// Get the email and code from the request
	var request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := c.BodyParser(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request data")
	}

	// Validate the verification code
	isValid, err := h.Service.ValidateCode(request.Email, request.Code)
	if err != nil {
		log.Printf("Error validating code: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to validate verification code")
	}

	if !isValid {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid verification code")
	}

	return c.JSON(fiber.Map{
		"message": "Verification code validated successfully",
	})
}
